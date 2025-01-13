package ui

import (
	"fmt"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/audio"
	"github.com/liqmix/slaptrax/internal/beats"
	"github.com/liqmix/slaptrax/internal/external"
	"github.com/liqmix/slaptrax/internal/l"
	"github.com/liqmix/slaptrax/internal/types"
	"github.com/tinne26/etxt"
)

type Leaderboard struct {
	Component

	scores    *UIGroup
	itemCache map[string]map[int][]*LeaderboardItem
	title     *Element

	loading   bool
	connected bool

	anim   *beats.PulseAnimation
	bmager *beats.Manager
}

type LeaderboardItem struct {
	Element
	userScore *external.Score

	elements []*Element
}

func NewLeaderboardItem(score *external.Score, center Point) *LeaderboardItem {
	i := &LeaderboardItem{
		userScore: score,
		elements:  make([]*Element, 0),
	}

	itemSize := Point{X: 0.35, Y: 0.1}
	elementSize := Point{X: itemSize.X / 2, Y: itemSize.Y * 0.75}
	leftSide := center.Translate(-itemSize.X/2, 0)
	rightSide := center.Translate(itemSize.X/2, 0)
	offset := 0.05
	textOpts := GetDefaultTextOptions()

	g := NewUIGroup()
	g.SetDisabled(true)
	g.SetHorizontal()
	g.SetSize(itemSize)
	g.SetCenter(center)

	// Username
	yOffset := 0.01

	e := NewElement()
	e.SetText(score.Username)
	e.SetSize(elementSize)
	e.SetTextBold(true)
	e.SetTextAlign(etxt.Left)
	e.SetTextScale(1.2)
	e.SetCenter(Point{X: leftSide.X + offset, Y: center.Y + yOffset})
	i.elements = append(i.elements, e)

	// Rank
	e = NewElement()
	e.SetText(fmt.Sprintf("%.2f", score.Rank))
	e.SetTextAlign(etxt.Right)
	e.SetSize(elementSize)
	e.SetCenter(Point{X: leftSide.X + offset*0.9, Y: center.Y + yOffset})
	i.elements = append(i.elements, e)

	// Title
	title := types.RankTitleFromRank(score.Rank)
	color := title.Color()
	e = NewElement()
	e.SetText(title.String())
	e.SetTextAlign(etxt.Left)
	e.SetTextScale(0.75)
	e.SetTextColor(color)
	e.SetSize(elementSize)
	e.SetCenter(Point{X: leftSide.X + offset, Y: center.Y + yOffset*3.2})
	i.elements = append(i.elements, e)

	e = NewElement()
	rating := types.GetSongRating(score.Score)
	e.SetText(rating.String())
	e.SetSize(elementSize)
	e.SetCenter(center)
	e.SetTextColor(rating.Color())
	i.elements = append(i.elements, e)

	e = NewElement()
	scoreText := strconv.Itoa(int(score.Score))
	scoreTextWidth := TextWidth(textOpts, scoreText)
	e.SetText(scoreText)
	e.SetSize(elementSize)
	e.SetCenter(Point{X: rightSide.X - offset - scoreTextWidth/2, Y: center.Y})
	i.elements = append(i.elements, e)

	return i
}

func (l *LeaderboardItem) SetElementScale(scale float64) {
	for _, e := range l.elements {
		e.SetRenderTextScale(scale)
	}
}

func (l *LeaderboardItem) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	for _, e := range l.elements {
		e.Draw(screen, opts)
	}
}

func NewLeaderboard() *Leaderboard {
	lb := &Leaderboard{
		Component: Component{},
		itemCache: make(map[string]map[int][]*LeaderboardItem),
		loading:   true,
		connected: external.HasConnection(),
		anim:      beats.NewPulseAnimation(1.25, 0.01),
	}
	bmager := beats.NewManager(125, 0)
	for i := 0; i < 4; i++ {
		bmager.SetTrigger(beats.BeatPosition{Numerator: i, Denominator: 4}, func() {
			lb.anim.Pulse()
		})
	}
	lb.bmager = bmager

	size := Point{X: 0.33, Y: 0.6}
	center := Point{X: 0.2, Y: 0.35}
	items := NewUIGroup()
	items.SetDisabled(true)
	items.SetVertical()
	items.SetSize(size)
	items.SetCenter(center)
	lb.scores = items

	text := NewElement()
	text.SetDisabled(true)
	text.SetText(l.String(l.LEADERBOARD))
	text.SetCenter(Point{X: center.X, Y: 0.25})
	text.SetTextBold(true)
	text.SetTextScale(2.5)
	lb.title = text
	return lb
}

func (l *Leaderboard) CreateItems(scores []external.Score, song string, difficulty int) {
	if cachedSong, ok := l.itemCache[song]; ok {
		if cached, ok := cachedSong[difficulty]; ok {
			l.SetItems(cached)
			return
		}
	}

	center := l.scores.GetCenter()
	offset := 0.05
	items := make([]*LeaderboardItem, len(scores))
	for i, score := range scores {
		items[i] = NewLeaderboardItem(&score, Point{X: center.X, Y: center.Y + offset*float64(i)})
	}
	if _, ok := l.itemCache[song]; !ok {
		l.itemCache[song] = make(map[int][]*LeaderboardItem)
	}
	if _, ok := l.itemCache[song][difficulty]; !ok {
		l.itemCache[song][difficulty] = make([]*LeaderboardItem, len(items))
	}
	l.itemCache[song][difficulty] = items
	l.SetItems(items)
}

func (l *Leaderboard) SetItems(items []*LeaderboardItem) {
	l.scores.items = make([]Componentable, 0)
	for _, item := range items {
		l.scores.Add(item)
	}
	l.loading = false
}

func (l *Leaderboard) FetchScores(song string, difficulty int, bpm float64) {
	l.loading = true
	l.bmager.SetBPM(bpm)

	if difficulties, ok := l.itemCache[song]; ok {
		if scores, ok := difficulties[difficulty]; ok {
			l.SetItems(scores)
			return
		}
	}
	go func() {
		scores, err := external.GetLeaderboard(song, difficulty)
		if err != nil {
			l.loading = false
			l.SetItems([]*LeaderboardItem{})
			return
		}
		l.CreateItems(scores, song, difficulty)
	}()
}

func (l *Leaderboard) Update() {
	if !l.connected {
		return
	}
	l.bmager.Update(audio.GetSongPreviewPositionMS())
	l.anim.Update()
	scale := l.anim.GetScale()
	if len(l.scores.items) > 0 {
		first := l.scores.items[0].(*LeaderboardItem)
		first.SetElementScale(scale)
	}
	l.scores.Update()
	l.title.Update()
}

func (l *Leaderboard) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	if !l.connected {
		return
	}
	l.title.Draw(screen, opts)
	if !l.loading {
		l.scores.Draw(screen, opts)
	}
}
