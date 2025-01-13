package state

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/assets"
	"github.com/liqmix/ebiten-holiday-2024/internal/external"
	"github.com/liqmix/ebiten-holiday-2024/internal/l"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
	"github.com/liqmix/ebiten-holiday-2024/internal/user"
)

type HowToPlayState struct {
	types.BaseGameState

	text       *ui.Element
	loginText  *ui.Element
	offsetText *ui.Element
	gameImage  *ui.Element
	panel      *ui.UIGroup

	textOpts   *ui.TextOptions
	wasNewUser bool
}

const gameImagePath = "game.png"
const slapNotePath = "slap_note.png"
const holdNotePath = "hold_note.png"
const multiNotePath = "multi_note.png"

func NewHowToPlayState() *HowToPlayState {
	state := HowToPlayState{}

	center := ui.Point{X: 0.5, Y: 0.2}
	text := ui.NewElement()
	text.SetSize(ui.Point{X: 0.8, Y: 0.8})
	text.SetCenter(center)
	text.SetText(l.String(l.DIALOG_HOW_TO_PLAY))
	text.SetTextScale(1.3)
	state.text = text

	offsetText := ui.NewElement()
	offsetText.SetSize(ui.Point{X: 0.8, Y: 0.8})
	offsetText.SetCenter(center.Translate(0, 0.15))
	offsetText.SetText(l.String(l.DIALOG_BE_SURE_TO_OFFSET))
	offsetText.SetTextScale(1.3)
	offsetText.SetTextColor(ui.CornerTrackColor())
	offsetText.SetTextBold(true)
	state.offsetText = offsetText

	if external.HasConnection() {
		loginText := ui.NewElement()
		loginText.SetSize(ui.Point{X: 0.35, Y: 0.5})
		loginText.SetTextScale(1.2)
		loginText.SetCenter(ui.Point{X: 0.5, Y: 0.80})
		loginText.SetText(l.String(l.DIALOG_BE_SURE_TO_LOGIN))
		loginText.SetTextColor(ui.CenterTrackColor())
		loginText.SetTextBold(true)
		state.loginText = loginText
	}

	gameImage := ui.NewElement()
	gameImage.SetSize(ui.Point{X: 0.4, Y: 0.4})
	gameImage.SetCenter(ui.Point{X: 0.5, Y: 0.57})
	gameImage.SetImage(assets.GetImage(gameImagePath))
	state.gameImage = gameImage

	panel := ui.NewUIGroup()
	panel.SetPaneled(true)
	panel.SetCenter(ui.Point{X: 0.5, Y: 0.5})
	panel.SetSize(ui.Point{X: 0.85, Y: 0.85})

	b := ui.NewElement()
	b.SetCenter(ui.Point{X: 0.5, Y: 0.87})
	b.SetSize(ui.Point{X: 0.2, Y: 0.1})
	b.SetText(l.String(l.OK))
	b.SetTrigger(func() {
		if user.Current().Settings.IsNewUser {
			state.SetNextState(types.GameStateKeyConfig, &FloatStateArgs{
				Cb: func() {
					user.Current().Settings.IsNewUser = false
					state.SetNextState(types.GameStateBack, nil)
				},
			})
		} else {
			state.SetNextState(types.GameStateBack, nil)
		}
	})
	panel.Add(b)
	state.panel = panel

	return &state
}

func (s *HowToPlayState) Update() error {
	s.BaseGameState.Update()
	s.panel.Update()
	s.text.Update()
	s.offsetText.Update()
	if s.loginText != nil {
		s.loginText.Update()
	}
	return nil
}

func (s *HowToPlayState) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	s.panel.Draw(screen, opts)
	s.gameImage.Draw(screen, opts)
	s.text.Draw(screen, opts)
	s.offsetText.Draw(screen, opts)
	if s.loginText != nil {
		s.loginText.Draw(screen, opts)
	}
}
