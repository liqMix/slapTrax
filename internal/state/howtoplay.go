package state

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/assets"
	"github.com/liqmix/ebiten-holiday-2024/internal/l"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
)

type HowToPlayState struct {
	types.BaseGameState

	gameImage *ui.Element
	text      *ui.Element

	buttons *ui.UIGroup
}

const gameImagePath = "game.png"

func NewHowToPlayState() *HowToPlayState {
	state := HowToPlayState{}

	center := ui.Point{
		X: 0.5,
		Y: 0.5,
	}

	gameImage := ui.NewElement()
	gameImage.SetSize(ui.Point{X: 0.5, Y: 0.5})
	gameImage.SetCenter(center)
	gameImage.SetImage(assets.GetImage(gameImagePath))

	text := ui.NewElement()
	text.SetSize(ui.Point{X: 0.5, Y: 0.5})
	text.SetCenter(center)
	text.SetText(l.String(l.DIALOG_HOW_TO_PLAY))

	buttons := ui.NewUIGroup()
	buttons.SetPaneled(true)
	buttons.SetCenter(ui.Point{X: 0.5, Y: 0.5})
	buttons.SetSize(ui.Point{X: 0.5, Y: 0.5})

	b := ui.NewElement()
	b.SetCenter(ui.Point{X: 0.5, Y: 0.7})
	b.SetSize(ui.Point{X: 0.2, Y: 0.1})
	b.SetText(l.String(l.BACK))
	b.SetTrigger(func() {
		state.SetNextState(types.GameStateBack, nil)
	})
	buttons.Add(b)

	state.gameImage = gameImage
	state.text = text
	state.buttons = buttons

	return &state
}

func (s *HowToPlayState) Update() error {
	s.buttons.Update()
	return nil
}

func (s *HowToPlayState) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	s.buttons.Draw(screen, opts)
	s.gameImage.Draw(screen, opts)
	s.text.Draw(screen, opts)
}
