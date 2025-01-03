package state

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/assets"
	"github.com/liqmix/ebiten-holiday-2024/internal/input"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
)

type Title struct {
	types.BaseGameState

	group *ui.UIGroup
}

func NewTitleState() *Title {
	state := Title{}
	group := ui.NewUIGroup()

	center := ui.Point{
		X: 0.75,
		Y: 0.5,
	}

	offset := float64(ui.TextHeight() * 2)

	// Play
	e := ui.NewElement()
	e.SetCenter(center)
	e.SetText(assets.String(types.L_STATE_PLAY))
	e.SetTrigger(func() {
		state.SetNextState(types.GameStateSongSelection, nil)
	})
	group.Add(e)
	center.Y += offset

	// Settings
	e = ui.NewElement()
	e.SetCenter(center)
	e.SetText(assets.String(types.L_STATE_SETTINGS))
	e.SetTrigger(func() {
		state.NextState = types.GameStateSettings
	})
	group.Add(e)
	center.Y += offset

	// Offset
	e = ui.NewElement()
	e.SetCenter(center)
	e.SetText(assets.String(types.L_STATE_OFFSET))
	e.SetTrigger(func() {
		state.NextState = types.GameStateOffset
	})
	group.Add(e)
	center.Y += offset

	// Exit
	e = ui.NewElement()
	e.SetCenter(center)
	e.SetText(assets.String(types.L_EXIT))
	e.SetTrigger(func() {
		panic("lol, lmao")
	})
	group.Add(e)

	state.group = group
	return &state
}

func (s *Title) Update() error {
	s.group.Update()

	if input.K.Is(ebiten.KeyEscape, input.JustPressed) {
		panic("lol, lmao")
	}
	return nil
}

func (s *Title) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	s.group.Draw(screen, opts)
}
