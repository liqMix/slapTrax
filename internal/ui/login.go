package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/l"
	"github.com/liqmix/slaptrax/internal/types"
)

type LoginModal struct {
	Component
	group *UIGroup
	panel *Component

	errorText  string
	onLogin    func(string, string)
	onContinue func()
}

func NewLoginModal() *LoginModal {
	p := &Component{}
	p.SetCenter(Point{X: 0.5, Y: 0.5})
	p.SetSize(Point{X: 0.6, Y: 0.8})
	p.SetPaneled(true)

	m := &LoginModal{}

	// Text Inputs
	inputSize := Point{X: 0.2, Y: 0.1}
	inputPos := Point{X: 0.5, Y: 0.4}
	textScale := 1.5

	username := NewTextInput(l.String(l.LOGIN_USERNAME))
	username.SetPaneled(true)
	username.SetSize(inputSize)
	username.SetTextScale(textScale)
	username.SetCenter(inputPos)

	password := NewTextInput(l.String(l.LOGIN_PASSWORD))
	password.SetIsPassword(true)
	password.SetPaneled(true)
	password.SetSize(inputSize)
	password.SetTextScale(textScale)
	password.SetCenter(inputPos.Translate(0, 0.12))

	// Login Button
	textScale = 3
	login := NewElement()
	loginSize := Point{X: 0.15, Y: 0.1}
	loginCenter := Point{X: 0.5, Y: 0.65}
	login.SetSize(loginSize)
	login.SetCenter(loginCenter)
	login.SetText(l.String(l.LOGIN_LOGIN))
	login.SetTextColor(types.TrackTypeCenter.Color())
	login.SetTextScale(textScale)
	login.SetTrigger(func() {
		if m.onLogin != nil {
			m.onLogin(username.GetText(), password.GetText())
		}
	})

	// Continue Button
	textScale = 1
	continueButton := NewElement()
	continueButton.SetSize(loginSize)
	continueButton.SetCenter(Point{X: loginCenter.X, Y: 0.8})
	continueButton.SetText(l.String(l.LOGIN_CONTINUE))
	continueButton.SetTextScale(textScale)
	continueButton.SetTrigger(func() {
		if m.onContinue != nil {
			m.onContinue()
		}
	})

	g := NewUIGroup()
	g.Add(username)
	g.Add(password)
	g.Add(login)
	g.Add(continueButton)

	m.panel = p
	m.group = g
	return m
}

func (m *LoginModal) SetOnLogin(f func(string, string)) {
	m.onLogin = f
}

func (m *LoginModal) SetOnContinue(f func()) {
	m.onContinue = f
}

func (m *LoginModal) SetError(err string) {
	m.errorText = err
}

func (m *LoginModal) Update() {
	if m.hidden {
		return
	}
	m.group.Update()
}

func (m *LoginModal) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	if m.hidden {
		return
	}

	m.Component.Draw(screen, opts)

	m.panel.Draw(screen, opts)
	// Draw components
	m.group.Draw(screen, opts)

	// Draw error if any
	if m.errorText != "" {
		errorColor := color.RGBA{255, 100, 100, 255}
		pos := &Point{X: m.center.X, Y: 0.15}
		textOpts := GetDefaultTextOptions()
		textOpts.Color = errorColor
		DrawTextAt(screen, m.errorText, pos, textOpts, opts)
	}
}
