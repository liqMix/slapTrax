package state

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/external"
	"github.com/liqmix/ebiten-holiday-2024/internal/l"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
	"github.com/liqmix/ebiten-holiday-2024/internal/user"
)

type LoginState struct {
	types.BaseGameState

	connectedToServer bool
	username          string
	password          string

	text        *ui.Element
	inputGroup  *ui.UIGroup
	buttonGroup *ui.UIGroup
}

func NewLoginState() *LoginState {
	ls := LoginState{}
	if !external.PingServer() {
		ls.connectedToServer = false
	} else {
		ls.connectedToServer = true
	}

	center := ui.Point{
		X: 0.5,
		Y: 0.4,
	}
	loginText := ui.NewElement()
	loginText.SetCenter(center)
	loginText.SetDisabled(true)
	loginText.SetText(l.String(l.LOGIN_TEXT_OFFLINE))

	center.Y += 0.1

	// Buttons
	buttons := ui.NewUIGroup()
	buttons.SetHorizontal()
	buttons.SetPaneled(true)

	// Inputs
	inputs := ui.NewUIGroup()
	buttonPosition := ui.Point{
		X: 0.44,
		Y: 0.57,
	}
	if ls.connectedToServer {
		fmt.Println("Connected to server")
		loginText.SetText(l.String(l.LOGIN_TEXT_ONLINE))
		kie := ui.NewKeyboardInputElement()
		kie.SetCenter(center)
		kie.SetLabel(l.String(l.LOGIN_USERNAME))
		inputs.Add(kie)
		center.Y += 0.1

		kie = ui.NewKeyboardInputElement()
		kie.SetCenter(center)
		kie.SetLabel(l.String(l.LOGIN_PASSWORD))
		inputs.Add(kie)
		center.Y += 0.1

		b := ui.NewElement()
		b.SetCenter(buttonPosition)
		b.SetText(l.String(l.LOGIN_LOGIN))
		b.SetTrigger(func() {
			ls.SetNextState(types.GameStateTitle, nil)
		})
		buttons.Add(b)
		buttonPosition.X += 0.15
	} else {
		b := ui.NewElement()
		b.SetCenter(buttonPosition)
		b.SetText(l.String(l.LOGIN_SAVE_LOCAL))
		b.SetTrigger(func() {
			user.GetCurrentUser().IsNewUser = false
			user.Save()
			user.GetCurrentUser().BypassLogin = true
			ls.SetNextState(types.GameStateTitle, nil)
		})
		buttons.Add(b)
		buttonPosition.X += 0.15
	}

	b := ui.NewElement()
	b.SetCenter(buttonPosition)
	b.SetText(l.String(l.LOGIN_GUEST))
	b.SetTrigger(func() {
		user.GetCurrentUser().IsNewUser = false
		user.GetCurrentUser().IsGuest = true
		user.GetCurrentUser().BypassLogin = true
		ls.SetNextState(types.GameStateTitle, nil)
	})
	buttons.Add(b)
	buttons.SetCenter(ui.Point{X: 0.5, Y: 0.5})
	buttons.SetSize(ui.Point{X: 0.35, Y: 0.35})

	inputs.SetCenter(ui.Point{X: 0.5, Y: 0.5})
	inputs.SetSize(ui.Point{X: 0.35, Y: 0.35})

	ls.text = loginText
	ls.inputGroup = inputs
	ls.buttonGroup = buttons

	// b := ui.NewButton()
	// b.SetCenter(center)
	// b.SetText("Create Account")
	// b.SetTrigger(func() {

	// g.SetCenter(ui.Point{X: 0.5, Y: 0.5})
	// g.SetSize(ui.Point{X: 0.25, Y: 0.25})
	return &ls
}

func (s *LoginState) Update() error {
	s.inputGroup.Update()
	s.buttonGroup.Update()
	return nil
}

func (s *LoginState) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	s.inputGroup.Draw(screen, opts)
	s.buttonGroup.Draw(screen, opts)
}
