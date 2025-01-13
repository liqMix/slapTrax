package state

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/audio"
	"github.com/liqmix/slaptrax/internal/beats"
	"github.com/liqmix/slaptrax/internal/display"
	"github.com/liqmix/slaptrax/internal/external"
	"github.com/liqmix/slaptrax/internal/l"
	"github.com/liqmix/slaptrax/internal/logger"
	"github.com/liqmix/slaptrax/internal/types"
	"github.com/liqmix/slaptrax/internal/ui"
	"github.com/tinne26/etxt"
)

type ResultStateArgs struct {
	Score *types.Score
}

type Result struct {
	types.BaseGameState

	previousScore *external.Score
	score         *types.Score
	rating        *ui.Element
	group         *ui.UIGroup
	text          *ebiten.Image

	anim   *beats.PulseAnimation
	bmager *beats.Manager
}

func NewResultState(args *ResultStateArgs) *Result {
	audio.FadeInBGM()

	score := args.Score
	r := &Result{
		score: score,
		anim:  beats.NewPulseAnimation(1.5, 0.01),
	}
	bmager := beats.NewManager(125, 0)
	for i := 0; i < 4; i++ {
		bmager.SetTrigger(beats.BeatPosition{Numerator: i, Denominator: 4}, func() {
			r.anim.Pulse()
		})
	}
	r.bmager = bmager
	if external.HasConnection() {
		go func() {
			// Get previous score
			previousScore, err := external.GetScore(score.Song.Hash, int(score.Difficulty))
			if err == nil {
				r.previousScore = previousScore
			}

			// Send score
			err = external.AddScore(&external.Score{
				SongHash:   score.Song.Hash,
				Difficulty: int(score.Difficulty),
				Score:      score.TotalScore,
				MaxCombo:   score.MaxCombo,
				Accuracy:   score.GetAccuracy(),
				PlayedAt:   time.Now(),
			})
			if err != nil {
				logger.Error(err.Error())
			}
		}()
	}

	g := ui.NewUIGroup()
	g.SetHorizontal()

	//// Buttons
	position := ui.Point{X: 0.45, Y: 0.85}
	e := ui.NewElement()
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

	img := display.NewRenderImage()
	center := ui.Point{X: 0.5, Y: 0.25}
	textOpts := ui.GetDefaultTextOptions()

	songRating := types.GetSongRating(score.TotalScore)
	e = ui.NewElement()
	e.SetSize(ui.Point{X: 0.1, Y: 0.1})
	e.SetCenter(center)
	e.SetTextScale(5)
	e.SetTextBold(true)
	e.SetText(songRating.String())
	e.SetTextColor(songRating.Color())
	e.SetDisabled(true)
	r.rating = e
	r.group = g

	textOpts.Scale = 2
	center.Y += 0.15
	ui.DrawTextAt(img, fmt.Sprintf("%d", score.TotalScore), &center, textOpts, nil)
	center.Y += 0.05

	textOpts.Scale = 1
	ui.DrawTextAt(img, fmt.Sprintf("MAX COMBO: %d", score.MaxCombo), &center, textOpts, nil)
	center.Y += 0.1

	detailsStart := center.Y

	left := ui.Point{X: 0.5, Y: center.Y}
	leftTextOpts := textOpts
	leftTextOpts.Align = etxt.Left

	right := ui.Point{X: 0.5, Y: center.Y}
	rightTextOpts := textOpts
	rightTextOpts.Align = etxt.Right

	yOffset := 0.05
	ui.DrawTextAt(img, "TOTAL", &left, leftTextOpts, nil)
	ui.DrawTextAt(img, fmt.Sprintf("%d", score.TotalNotes), &right, rightTextOpts, nil)
	right.Y += yOffset
	left.Y += yOffset

	leftTextOpts.Color = types.Perfect.Color().C()
	ui.DrawTextAt(img, types.Perfect.String(), &left, leftTextOpts, nil)
	ui.DrawTextAt(img, fmt.Sprintf("%d", score.Perfect), &right, rightTextOpts, nil)
	right.Y += yOffset
	left.Y += yOffset

	leftTextOpts.Color = types.Good.Color().C()
	ui.DrawTextAt(img, types.Good.String(), &left, leftTextOpts, nil)
	ui.DrawTextAt(img, fmt.Sprintf("%d", score.Good), &right, rightTextOpts, nil)
	right.Y += yOffset
	left.Y += yOffset

	leftTextOpts.Color = types.Bad.Color().C()
	ui.DrawTextAt(img, types.Bad.String(), &left, leftTextOpts, nil)
	ui.DrawTextAt(img, fmt.Sprintf("%d", score.Bad), &right, rightTextOpts, nil)
	right.Y += yOffset
	left.Y += yOffset

	leftTextOpts.Color = types.Miss.Color().C()
	ui.DrawTextAt(img, types.Miss.String(), &left, leftTextOpts, nil)
	ui.DrawTextAt(img, fmt.Sprintf("%d", score.Miss), &right, rightTextOpts, nil)

	leftTextOpts.Color = types.LightBlue.C()
	ui.DrawTextAt(img, fmt.Sprintf("EARLY\n%d", score.Early), &ui.Point{X: 0.33, Y: detailsStart}, leftTextOpts, nil)
	rightTextOpts.Color = types.Purple.C()
	ui.DrawTextAt(img, fmt.Sprintf("LATE\n%d", score.Late), &ui.Point{X: 0.66, Y: detailsStart}, rightTextOpts, nil)
	r.text = img
	return r
}

func (r *Result) Update() error {
	if audio.IsSongPlaying() {
		// hmm
		audio.StopSong()
	}

	r.BaseGameState.Update()
	r.bmager.Update(audio.GetBGMPositionMS())

	r.anim.Update()
	scale := r.anim.GetScale()
	r.rating.SetRenderTextScale(scale)

	r.group.Update()
	return nil
}

func (r *Result) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	screen.DrawImage(r.text, nil)
	r.group.Draw(screen, opts)
	r.rating.Draw(screen, opts)
	if r.previousScore != nil {
		// Draw previous score
		textOpts := ui.GetDefaultTextOptions()
		textOpts.Scale = 1.0
		textOpts.Color = types.Gray.C()
		ui.DrawTextAt(screen, fmt.Sprintf("PREVIOUS SCORE: %d", r.previousScore.Score), &ui.Point{X: 0.5, Y: 0.49}, textOpts, nil)
	}
}
