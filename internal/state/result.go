package state

import (
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/external"
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

	go func() {
		// Send score
		external.AddScore(&external.Score{
			Song:       score.Song.Hash,
			Difficulty: score.Difficulty.String(),
			Score:      score.TotalScore,
			MaxCombo:   score.MaxCombo,
			Accuracy:   score.GetAccuracy(),
			PlayedAt:   time.Now(),
		})
	}()

	r := &Result{score: score}
	g := ui.NewUIGroup()
	g.SetHorizontal()

	var e *ui.Element

	//// Score + Rating
	size := ui.Point{X: 0.4, Y: 0.2}
	position := ui.Point{X: 0.5, Y: 0.15}
	e = ui.NewElement()
	e.SetCenter(position)
	e.SetSize(size)
	e.SetText(types.GetSongRating(score.TotalScore).String())
	r.elements = append(r.elements, e)
	position.Y += 0.1

	size = size.Scale(0.5)
	e = ui.NewElement()
	e.SetCenter(position)
	e.SetSize(size)
	e.SetText(strconv.Itoa(score.TotalScore))
	r.elements = append(r.elements, e)
	position.Y += 0.2

	//// Accuracy

	//// Leaderboard

	//// Player Stats

	//// Buttons
	position = ui.Point{X: 0.45, Y: 0.85}
	e = ui.NewElement()
	e.SetCenter(position)
	e.SetText(l.String(l.CONTINUE))
	e.SetTrigger(func() {
		r.SetNextState(types.GameStateSongSelection, nil)
	})
	g.Add(e)

	position.X += 0.1
	e = ui.NewElement()
	e.SetCenter(position)
	e.SetText(l.String(l.STATE_PLAY_RESTART))
	e.SetTrigger(func() {
		r.SetNextState(types.GameStatePlay, &PlayArgs{Song: r.score.Song, Difficulty: r.score.Difficulty})
	})
	g.Add(e)
	r.buttonGroup = g
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
