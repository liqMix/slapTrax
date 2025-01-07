package state

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/assets"
	"github.com/liqmix/ebiten-holiday-2024/internal/input"
	"github.com/liqmix/ebiten-holiday-2024/internal/l"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
	"github.com/liqmix/ebiten-holiday-2024/internal/user"
)

type Title struct {
	types.BaseGameState

	locale *ui.Button
	group  *ui.UIGroup
}

func NewTitleState() *Title {
	state := Title{}
	group := ui.NewUIGroup()

	center := ui.Point{
		X: 0.75,
		Y: 0.5,
	}

	offset := float64(ui.TextHeight(nil) * 2)

	// Play
	play := ui.NewElement()
	play.SetCenter(center)
	play.SetText(l.String(l.STATE_PLAY))
	play.SetTrigger(func() {
		if !user.S.PromptedOffsetCheck {
			txt := ui.NewElement()
			txt.SetCenter(ui.Point{X: 0.5, Y: 0.45})
			txt.SetScale(1.2)
			txt.SetText(l.String(l.DIALOG_CHECKOFFSETS))
			txt.SetTextColor(types.TrackTypeCenter.Color())
			txt.SetDisabled(true)

			g := ui.NewUIGroup()
			g.SetPaneled(true)

			e := ui.NewElement()
			e.SetCenter(ui.Point{X: 0.5, Y: 0.6})
			e.SetText("OK")
			e.SetTrigger(func() {
				user.S.PromptedOffsetCheck = true
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
					if user.S.PromptedOffsetCheck {
						setNextState(types.GameStateBack, nil)
					} else {
						txt.Update()
						g.Update()
					}
				},
			})
			return
		}
		state.SetNextState(types.GameStateSongSelection, nil)
	})
	group.Add(play)
	center.Y += offset

	// Settings
	settings := ui.NewElement()
	settings.SetCenter(center)
	settings.SetText(l.String(l.STATE_SETTINGS))
	settings.SetTrigger(func() {
		state.NextState = types.GameStateSettings
	})
	group.Add(settings)
	center.Y += offset

	// // Offset
	// e = ui.NewElement()
	// e.SetCenter(center)
	// e.SetText(l.STATE_OFFSET))
	// e.SetTrigger(func() {
	// 	state.NextState = types.GameStateOffset
	// })
	// group.Add(e)
	// center.Y += offset

	// Exit
	exit := ui.NewElement()
	exit.SetCenter(center)
	exit.SetText(l.String(l.EXIT))
	exit.SetTrigger(func() {
		panic("lol, lmao")
	})
	group.Add(exit)

	locales := assets.Locales()
	current := assets.CurrentLocale()
	idx := 0
	for i, l := range locales {
		if l == current {
			idx = i
			break
		}
	}
	b := ui.NewButton()
	b.SetCenter(ui.Point{X: 0.9, Y: 0.1})
	b.SetSize(ui.Point{X: 0.05, Y: 0.05})
	b.SetTrigger(func() {
		idx = (idx + 1) % len(locales)
		assets.SetLocale(locales[idx])
		b.SetImage(assets.Flag())

		//lol
		play.SetText(l.String(l.STATE_PLAY))
		settings.SetText(l.String(l.STATE_SETTINGS))
		exit.SetText(l.String(l.EXIT))
	})
	flag := assets.Flag()
	if flag != nil {
		b.SetImage(flag)
	} else {
		b.SetText(l.String(l.LOCALE))
	}
	state.locale = b
	state.group = group
	return &state
}

func (s *Title) Update() error {
	s.group.Update()
	s.locale.Update()
	if input.K.Is(ebiten.KeyEscape, input.JustPressed) {
		panic("lol, lmao")
	}
	if user.ShowLogin() {
		s.SetNextState(types.GameStateLogin, nil)
	}
	return nil
}

func (s *Title) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	s.locale.Draw(screen, opts)
	s.group.Draw(screen, opts)
}
