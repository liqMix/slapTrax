package state

import (
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/l"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
)

type ResultStateArgs struct {
	Score *types.Score
}

type Result struct {
	types.BaseGameState

	score *types.Score

	elements    []*ui.Element
	buttonGroup *ui.UIGroup
}

func NewResultState(args *ResultStateArgs) *Result {
	score := args.Score
	r := &Result{score: score}
	g := ui.NewUIGroup()
	g.SetHorizontal()

	var e *ui.Element

	//// Score + Rating
	position := ui.Point{X: 0.5, Y: 0.15}
	e = ui.NewElement()
	e.SetCenter(position)
	e.SetScale(4.0)
	e.SetText(score.GetRating().String())
	r.elements = append(r.elements, e)

	position.Y += 0.1
	e = ui.NewElement()
	e.SetCenter(position)
	e.SetScale(2.0)
	e.SetText(strconv.Itoa(score.GetScore()))
	r.elements = append(r.elements, e)
	position.Y += 0.2

	//// Accuracy

	//// Leaderboard

	//// Player Stats

	//// Buttons
	position = ui.Point{X: 0.45, Y: 0.85}
	b := ui.NewButton()
	b.SetCenter(position)
	b.SetText(l.String(l.BACK))
	b.SetTrigger(func() {
		r.SetNextState(types.GameStateSongSelection, nil)
	})
	g.Add(b)

	position.X += 0.1
	b = ui.NewButton()
	b.SetCenter(position)
	b.SetText(l.String(l.STATE_PLAY_RESTART))
	b.SetTrigger(func() {
		r.SetNextState(types.GameStatePlay, &PlayArgs{Song: r.score.Song, Difficulty: r.score.Difficulty})
	})
	g.Add(b)
	return r
}

func (r *Result) Update() error {
	for _, e := range r.elements {
		e.Update()
	}
	r.buttonGroup.Update()
	return nil
}

func (r *Result) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	for _, e := range r.elements {
		e.Draw(screen, opts)
	}
	r.buttonGroup.Draw(screen, opts)
}
