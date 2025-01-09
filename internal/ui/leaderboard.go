package ui

import (
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/external"
	"github.com/liqmix/ebiten-holiday-2024/internal/l"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
)

type Leaderboard struct {
	Component

	items     *UIGroup
	itemCache map[string]map[string][]*LeaderboardItem
	title     *Element

	loading   bool
	connected bool
}

type LeaderboardItem struct {
	Element
	score *external.Score

	group *UIGroup
}

func NewLeaderboardItem(score *external.Score, center Point) *LeaderboardItem {
	i := &LeaderboardItem{
		Element: *NewElement(),
		score:   score,
	}

	itemSize := Point{X: 0.3, Y: 0.1}
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

	e := NewElement()
	e.SetText(score.Username)
	e.SetSize(elementSize)
	e.SetCenter(Point{X: leftSide.X + offset, Y: center.Y})
	g.Add(e)

	e = NewElement()
	rating := types.GetSongRating(score.Score)
	e.SetText(rating.String())
	e.SetSize(elementSize)
	e.SetCenter(center)
	e.SetTextColor(rating.Color())
	g.Add(e)

	e = NewElement()
	scoreText := strconv.Itoa(int(score.Score))
	scoreTextWidth := TextWidth(textOpts, scoreText)
	e.SetText(scoreText)
	e.SetSize(elementSize)
	e.SetCenter(Point{X: rightSide.X - offset - scoreTextWidth/2, Y: center.Y})
	g.Add(e)

	i.group = g
	return i
}

func (l *LeaderboardItem) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	l.group.Draw(screen, opts)
}

func NewLeaderboard() *Leaderboard {
	lb := &Leaderboard{
		Component: Component{},
		itemCache: make(map[string]map[string][]*LeaderboardItem),
		loading:   true,
		connected: external.HasConnection(),
	}

	size := Point{X: 0.33, Y: 0.6}
	center := Point{X: 0.2, Y: 0.5}
	items := NewUIGroup()
	items.SetDisabled(true)
	items.SetVertical()
	items.SetSize(size)
	items.SetCenter(center)
	lb.items = items

	text := NewElement()
	text.SetDisabled(true)
	text.SetText(l.String(l.LEADERBOARD))
	text.SetCenter(Point{X: center.X, Y: 0.35})
	text.SetTextBold(true)
	text.SetTextScale(3)
	lb.title = text
	return lb
}

func (l *Leaderboard) CreateItems(scores []external.Score, song, difficulty string) {
	if cachedSong, ok := l.itemCache[song]; ok {
		if cached, ok := cachedSong[difficulty]; ok {
			l.SetItems(cached)
			return
		}
	}

	center := l.items.GetCenter()
	offset := 0.025
	items := make([]*LeaderboardItem, len(scores))
	for i, score := range scores {
		items[i] = NewLeaderboardItem(&score, Point{X: center.X, Y: center.Y + offset*float64(i)})
	}
	if _, ok := l.itemCache[song]; !ok {
		l.itemCache[song] = make(map[string][]*LeaderboardItem)
	}
	if _, ok := l.itemCache[song][difficulty]; !ok {
		l.itemCache[song][difficulty] = make([]*LeaderboardItem, len(items))
	}
	l.itemCache[song][difficulty] = items
	l.SetItems(items)
}

func (l *Leaderboard) SetItems(items []*LeaderboardItem) {
	l.items.items = make([]Componentable, 0)
	for _, item := range items {
		l.items.Add(item)
	}
	l.loading = false
}

func (l *Leaderboard) FetchScores(song string, difficulty types.Difficulty) {
	l.loading = true

	diffText := strconv.Itoa(int(difficulty))
	if scores, ok := l.itemCache[song][diffText]; ok {
		l.SetItems(scores)
		return
	}

	go func() {
		scores, err := external.GetLeaderboard(song, diffText)
		if err != nil {
			l.loading = false
			return
		}
		l.CreateItems(scores, song, diffText)
	}()
}

func (l *Leaderboard) Update() {
	if !l.connected {
		return
	}
	l.items.Update()
	l.title.Update()
}

func (l *Leaderboard) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	if !l.connected {
		return
	}
	l.title.Draw(screen, opts)
	if !l.loading {
		l.items.Draw(screen, opts)
	}
}
