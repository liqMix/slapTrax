package state

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/l"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
)

type PauseArgs struct {
	song       *types.Song
	difficulty types.Difficulty
	cb         func()
}

type Pause struct {
	types.BaseGameState

	group *ui.UIGroup
}

func NewPauseState(args *PauseArgs) *Pause {
	p := &Pause{}

	group := ui.NewUIGroup()
	group.SetPaneled(true)
	center := ui.Point{
		X: 0.5,
		Y: 0.4,
	}

	offset := float64(ui.TextHeight(nil) * 2)

	// Resume
	e := ui.NewElement()
	e.SetCenter(center)
	e.SetText(l.String(l.BACK))
	e.SetTrigger(func() {
		p.SetNextState(types.GameStateBack, nil)
		args.cb()
	})
	group.Add(e)
	center.Y += offset

	// Settings
	e = ui.NewElement()
	e.SetCenter(center)
	e.SetText(l.String(l.STATE_SETTINGS))
	e.SetTrigger(func() {
		p.SetNextState(types.GameStateSettings, nil)
	})
	group.Add(e)
	center.Y += offset

	// Restart
	e = ui.NewElement()
	e.SetCenter(center)
	e.SetText(l.String(l.STATE_PLAY_RESTART))
	e.SetTrigger(func() {
		p.SetNextState(types.GameStatePlay, &PlayArgs{
			Song:       args.song,
			Difficulty: args.difficulty,
		})
	})
	group.Add(e)
	center.Y += offset

	// Quit
	quit := ui.NewElement()
	quit.SetCenter(center)
	quit.SetText(l.String(l.EXIT))
	quit.SetTrigger(func() {
		p.SetNextState(types.GameStateSongSelection, nil)
	})
	group.Add(quit)

	group.SetCenter(ui.Point{X: 0.5, Y: 0.5})
	group.SetSize(ui.Point{X: 0.25, Y: 0.5})
	p.group = group
	return p
}

func (s *Pause) Update() error {
	s.group.Update()
	return nil
}

func (s *Pause) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	s.group.Draw(screen, opts)
}
