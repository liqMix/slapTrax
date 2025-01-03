package state

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/audio"
	"github.com/liqmix/ebiten-holiday-2024/internal/input"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
)

type DifficultySelectionArgs struct {
	song *types.Song
}

type DifficultySelection struct {
	types.BaseGameState

	song  *types.Song
	text  *ui.Element
	group *ui.UIGroup
}

func NewDifficultySelectionState(args *DifficultySelectionArgs) *DifficultySelection {
	d := &DifficultySelection{song: args.song}
	diffs := d.song.GetDifficulties()
	if len(diffs) == 0 {
		d.SetNextState(types.GameStateBack, nil)
	}

	group := ui.NewUIGroup()
	group.SetPaneled(true)
	group.SetHorizontal()

	center := ui.Point{
		X: 0.5,
		Y: 0.4,
	}
	e := ui.NewElement()
	e.SetCenter(center)
	e.SetScale(2.0)
	e.SetTextBold(true)
	e.SetText("Select Difficulty")
	d.text = e

	// group.Add(e)
	center.Y += 0.1

	for _, diff := range diffs {
		e := ui.NewElement()
		e.SetCenter(center)
		e.SetText(diff.String())
		e.SetTrigger(func() {
			audio.StopAll()
			d.SetNextState(types.GameStatePlay, &PlayArgs{
				Song:       args.song,
				Difficulty: diff,
			})
		})
		group.Add(e)
		center.X += 0.025
	}
	d.group = group
	group.SetCenter(ui.Point{X: 0.5, Y: 0.5})

	return d
}

func (s *DifficultySelection) Update() error {
	s.group.Update()
	if input.K.Is(ebiten.KeyEscape, input.JustPressed) {
		s.SetNextState(types.GameStateBack, nil)
	}
	return nil
}

func (s *DifficultySelection) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	s.group.Draw(screen, opts)
	s.text.Draw(screen, opts)
}
