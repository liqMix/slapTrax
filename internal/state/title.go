package state

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/assets"
	"github.com/liqmix/ebiten-holiday-2024/internal/external"
	"github.com/liqmix/ebiten-holiday-2024/internal/input"
	"github.com/liqmix/ebiten-holiday-2024/internal/l"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
	"github.com/liqmix/ebiten-holiday-2024/internal/user"
)

type Title struct {
	types.BaseGameState

	group   *ui.UIGroup
	welcome *ui.Element
	name    *ui.Element
}

func NewTitleState() *Title {
	state := Title{}

	// center := ui.Point{
	// 	X: 0.85,
	// 	Y: 0.65,
	// }
	center := ui.Point{
		X: 0.5,
		Y: 0.5,
	}
	welcome := ui.NewElement()
	welcome.SetSize(ui.Point{X: 0.15, Y: 0.1})
	welcome.SetCenter(ui.Point{
		X: center.X,
		Y: 0.25,
	})
	welcome.SetText(l.String(l.WELCOME))
	welcome.SetTextScale(3)

	name := ui.NewElement()
	name.SetSize(ui.Point{X: 0.05, Y: 0.05})
	name.SetTextScale(2)
	name.SetCenter(ui.Point{
		X: center.X,
		Y: 0.35,
	})
	name.SetText(user.Current().Username)

	group := ui.NewUIGroup()
	group.SetPaneled(true)

	offset := float64(ui.TextHeight(nil) * 2)
	buttonSize := ui.Point{
		X: 0.1,
		Y: 0.1,
	}
	textScale := 1.5

	// Play
	play := ui.NewElement()
	play.SetCenter(center)
	play.SetSize(buttonSize)
	play.SetText(l.String(l.STATE_PLAY))
	play.SetTextScale(textScale)
	play.SetTrigger(func() {
		if !user.S().PromptedOffsetCheck {
			txt := ui.NewElement()
			txt.SetCenter(ui.Point{X: 0.5, Y: 0.45})
			txt.SetSize(ui.Point{X: 0.5, Y: 0.2})
			txt.SetText(l.String(l.DIALOG_NEW_PLAYER))
			txt.SetTextColor(types.TrackTypeCenter.Color())
			txt.SetDisabled(true)

			g := ui.NewUIGroup()
			g.SetPaneled(true)

			e := ui.NewElement()
			e.SetCenter(ui.Point{X: 0.5, Y: 0.6})
			e.SetText("OK")
			e.SetTrigger(func() {
				user.S().PromptedOffsetCheck = true
				user.Save()
			})
			g.Add(e)
			g.SetCenter(ui.Point{X: 0.5, Y: 0.5})
			g.SetSize(ui.Point{X: 0.5, Y: 0.3})
			state.SetNextState(types.GameStateModal, &ModalStateArgs{
				Draw: func(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
					g.Draw(screen, opts)
					txt.Draw(screen, opts)
				},
				Update: func(setNextState func(s types.GameState, args interface{})) {
					if user.S().PromptedOffsetCheck {
						setNextState(types.GameStateBack, nil)
					} else {
						txt.Update()
						g.Update()
					}
				},
			})
			return
		}
		state.SetNextState(types.GameStateSongSelection, &SongSelectionArgs{
			Song: nil,
		})
	})
	group.Add(play)
	center.Y += offset

	// Settings
	settings := ui.NewElement()
	settings.SetSize(buttonSize)
	settings.SetCenter(center)
	settings.SetText(l.String(l.STATE_SETTINGS))
	settings.SetTextScale(textScale)
	settings.SetTrigger(func() {
		state.NextState = types.GameStateSettings
	})
	group.Add(settings)
	center.Y += offset

	// Locale
	locales := assets.Locales()
	current := assets.CurrentLocale()
	idx := 0
	for i, l := range locales {
		if l == current {
			idx = i
			break
		}
	}

	// login
	loginState := external.GetLoginState()
	var login *ui.Element
	if loginState != external.StatePlayOffline && external.HasConnection() {
		login = ui.NewElement()
		login.SetCenter(center)
		login.SetSize(buttonSize)
		login.SetTextScale(textScale)

		if loginState == external.StateOnline {
			login.SetText(l.String(l.LOGIN_LOGOUT))
			login.SetTrigger(func() {
				external.Logout()
				state.SetNextState(types.GameStateTitle, nil)
			})
		} else if loginState == external.StateOffline {
			login.SetText(l.String(l.LOGIN_LOGIN))
			login.SetTrigger(func() {
				state.SetNextState(types.GameStateLogin, nil)
			})
		}
		group.Add(login)
		center.Y += offset
	}

	// How to play
	howToPlay := ui.NewElement()
	howToPlay.SetCenter(center)
	howToPlay.SetSize(buttonSize)
	howToPlay.SetTextScale(textScale)
	howToPlay.SetText(l.String(l.STATE_HOW_TO_PLAY))
	howToPlay.SetTrigger(func() {
		state.SetNextState(types.GameStateHowToPlay, nil)
	})
	group.Add(howToPlay)
	center.Y += offset

	// Exit
	exit := ui.NewElement()
	exit.SetCenter(center)
	exit.SetText(l.String(l.EXIT))
	exit.SetSize(buttonSize)
	exit.SetTextScale(textScale)
	exit.SetTrigger(func() {
		state.SetNextState(types.GameStateExit, nil)
	})
	group.Add(exit)
	center.Y += offset * 2

	localeSize := buttonSize.Scale(0.75)
	b := ui.NewElement()
	b.SetCenter(center)
	b.SetSize(localeSize)
	b.SetTextScale(2)
	b.SetTrigger(func() {
		idx = (idx + 1) % len(locales)
		assets.SetLocale(locales[idx])
		flag := assets.Flag()
		if flag != nil {
			b.SetImage(assets.Flag())
			b.SetText("")
		} else {
			b.SetImage(nil)
			b.SetText(l.String(l.LOCALE))
		}

		//lol
		play.SetText(l.String(l.STATE_PLAY))
		settings.SetText(l.String(l.STATE_SETTINGS))
		exit.SetText(l.String(l.EXIT))
		if login != nil {
			if loginState == external.StateOnline {
				login.SetText(l.String(l.LOGIN_LOGOUT))
			} else if loginState == external.StateOffline {
				login.SetText(l.String(l.LOGIN_LOGIN))
			}
		}
	})

	flag := assets.Flag()
	if flag != nil {
		b.SetImage(flag)
	} else {
		b.SetText(l.String(l.LOCALE))
	}
	group.Add(b)

	state.group = group
	state.name = name
	state.welcome = welcome
	return &state
}

func (s *Title) Update() error {
	s.name.Update()
	s.group.Update()

	if input.K.Is(ebiten.KeyEscape, input.JustPressed) {
		return ebiten.Termination
	}
	return nil
}

func (s *Title) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	if external.GetLoginState() == external.StateOnline {
		s.welcome.Draw(screen, opts)
		s.name.Draw(screen, opts)
	}
	s.group.Draw(screen, opts)
}
