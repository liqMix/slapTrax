package state

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/assets"
	"github.com/liqmix/slaptrax/internal/audio"
	"github.com/liqmix/slaptrax/internal/beats"
	"github.com/liqmix/slaptrax/internal/external"
	"github.com/liqmix/slaptrax/internal/input"
	"github.com/liqmix/slaptrax/internal/l"
	"github.com/liqmix/slaptrax/internal/types"
	"github.com/liqmix/slaptrax/internal/ui"
	"github.com/liqmix/slaptrax/internal/user"
	"github.com/tinne26/etxt"
)

type TitleLogo struct {
	s    *ui.Element
	T    *ui.Element
	x    *ui.Element
	rest *ui.Element

	cornerAnim *beats.PulseAnimation
	centerAnim *beats.PulseAnimation
}

func NewTitleLogo() *TitleLogo {
	center := ui.Point{
		X: 0.5,
		Y: 0.25,
	}

	sSize := ui.Point{
		X: 0.046,
		Y: 0,
	}
	tSize := ui.Point{
		X: 0.08,
		Y: 0.08,
	}
	xSize := ui.Point{
		X: 0.06,
		Y: 0.01,
	}
	restSize := ui.Point{
		X: 0.4,
		Y: 0.2,
	}
	s := assets.GetImage("logo_s.png")
	T := assets.GetImage("logo_t.png")
	x := assets.GetImage("logo_x.png")
	rest := assets.GetImage("logo_rest.png")

	sE := ui.NewElement()
	sE.SetSize(sSize)
	sE.SetCenter(ui.Point{
		X: center.X - 0.268,
		Y: center.Y,
	})
	sE.SetImage(s)

	TE := ui.NewElement()
	TE.SetSize(tSize)
	TE.SetCenter(ui.Point{
		X: center.X,
		Y: center.Y - 0.03,
	})
	TE.SetImage(T)

	xE := ui.NewElement()
	xE.SetSize(xSize)
	xE.SetCenter(ui.Point{
		X: center.X + 0.225,
		Y: center.Y + 0.002,
	})
	xE.SetImage(x)

	restE := ui.NewElement()
	restE.SetSize(restSize)
	restE.SetCenter(ui.Point{
		X: center.X - 0.025,
		Y: center.Y,
	})
	restE.SetImage(rest)

	logo := &TitleLogo{
		s:          sE,
		T:          TE,
		x:          xE,
		rest:       restE,
		cornerAnim: beats.NewPulseAnimation(1.5, 0.01),
		centerAnim: beats.NewPulseAnimation(1.5, 0.01),
	}

	return logo
}

func (t *TitleLogo) Update() {
	t.cornerAnim.Update()
	t.centerAnim.Update()

	cornerScale := t.cornerAnim.GetScale()
	centerScale := t.centerAnim.GetScale()

	t.s.SetImageScale(cornerScale)
	t.T.SetImageScale(centerScale)
	t.x.SetImageScale(cornerScale)
}

func (t *TitleLogo) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	cornerC := &ebiten.DrawImageOptions{}
	cornerC.ColorScale.ScaleWithColor(ui.CornerTrackColor())

	centerC := &ebiten.DrawImageOptions{}
	centerC.ColorScale.ScaleWithColor(ui.CenterTrackColor())
	t.s.Draw(screen, cornerC)
	t.T.Draw(screen, centerC)
	t.x.Draw(screen, cornerC)
	t.rest.Draw(screen, opts)
}

type Title struct {
	types.BaseGameState
	logo     *TitleLogo
	group    *ui.UIGroup
	bmager   *beats.Manager
	thanks   *ui.UIGroup
	namesOne *ui.Element
	namesTwo *ui.Element
	anim     *beats.PulseAnimation
}

func NewTitleState() *Title {
	bgm := audio.GetBGM()
	if bgm == nil || !bgm.IsPlaying() {
		audio.PlayBGM(audio.BGMTitle)
	}

	state := Title{
		anim: beats.NewPulseAnimation(1.25, 0.005),
	}
	state.logo = NewTitleLogo()
	group := ui.NewUIGroup()
	group.SetCenter(ui.Point{
		X: 0.5,
		Y: 0.65,
	})

	group.SetPaneled(true)

	offset := float64(ui.TextHeight(nil) * 2)
	buttonSize := ui.Point{
		X: 0.1,
		Y: 0.1,
	}
	textScale := 1.5

	// Play
	center := ui.Point{
		X: 0.5,
		Y: 0.53,
	}
	if !external.HasConnection() {
		center.Y += offset / 2
	}
	play := ui.NewElement()
	play.SetCenter(center)
	play.SetSize(buttonSize)
	play.SetText(l.String(l.STATE_PLAY))
	play.SetTextScale(textScale)
	play.SetTrigger(func() {
		if user.Current().Settings.IsNewUser {
			state.SetNextState(types.GameStateHowToPlay, nil)
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

	// How to play
	howToPlay := ui.NewElement()
	howToPlay.SetCenter(center)
	howToPlay.SetSize(buttonSize)
	howToPlay.SetText(l.String(l.STATE_HOW_TO_PLAY))
	howToPlay.SetTextScale(textScale)
	howToPlay.SetTrigger(func() {
		state.SetNextState(types.GameStateHowToPlay, nil)
	})
	group.Add(howToPlay)
	center.Y += offset

	// Editor
	editor := ui.NewElement()
	editor.SetCenter(center)
	editor.SetSize(buttonSize)
	editor.SetText(l.String(l.STATE_EDITOR))
	editor.SetTextScale(textScale)
	editor.SetTrigger(func() {
		state.SetNextState(types.GameStateEditor, &EditorArgs{
			Song: nil, // Start with empty chart
		})
	})
	group.Add(editor)
	center.Y += offset

	// login
	loginState := external.GetLoginState()
	var login *ui.Element
	if external.HasConnection() {
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
			login.SetText(l.String(l.STATE_LOGIN))
			login.SetTrigger(func() {
				state.SetNextState(types.GameStateLogin, nil)
			})
		}
		group.Add(login)
		center.Y += offset
	}

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
	if !external.HasConnection() {
		center.Y += offset * 2
	} else {
		center.Y += offset * 1.7
	}

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
	localeSize := buttonSize.Scale(0.35)
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
		howToPlay.SetText(l.String(l.STATE_HOW_TO_PLAY))
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

	group.SetSize(ui.Point{
		X: 0.18,
		Y: center.Y - 0.53,
	})
	state.group = group

	thanksG := ui.NewUIGroup()
	thanksG.SetPaneled(true)
	thanksG.SetDisabled(true)

	thanksG.SetSize(ui.Point{
		X: 0.2,
		Y: 0.5,
	})
	thanksG.SetCenter(ui.Point{
		X: 0.95,
		Y: 0.75,
	})
	thanksTitle := ui.NewElement()
	thanksTitle.SetText("Special Thanks")
	thanksTitle.SetTextAlign(etxt.Right)
	thanksTitle.SetTextScale(1.5)
	thanksTitle.SetTextColor(ui.CornerTrackColor())
	thanksTitle.SetCenter(ui.Point{
		X: 0.99,
		Y: 0.55,
	})
	thanksG.Add(thanksTitle)

	thanksText := ui.NewElement()
	thanksText.SetText("Hypnogram\nMeebah")
	thanksText.SetTextAlign(etxt.Right)
	thanksText.SetTextColor(ui.CenterTrackColor())
	thanksText.SetCenter(ui.Point{
		X: 0.99,
		Y: 0.6,
	})
	state.namesOne = thanksText

	thanksG.Add(thanksText)

	thanksText = ui.NewElement()
	thanksText.SetText("For allowing me to \ninclude their music!")
	thanksText.SetTextAlign(etxt.Right)
	thanksText.SetTextScale(0.75)
	thanksText.SetCenter(ui.Point{
		X: 0.99,
		Y: 0.67,
	})
	thanksG.Add(thanksText)

	thanksText = ui.NewElement()
	thanksText.SetText("isacubes\nLennaLeFay\nWikiTay")
	thanksText.SetTextAlign(etxt.Right)
	thanksText.SetTextColor(ui.CenterTrackColor())
	thanksText.SetCenter(ui.Point{
		X: 0.99,
		Y: 0.75,
	})
	state.namesTwo = thanksText
	thanksG.Add(thanksText)

	thanksText = ui.NewElement()
	thanksText.SetText("For providing me with\nvaluable feedback during\ndevelopment!")
	thanksText.SetTextAlign(etxt.Right)
	thanksText.SetTextScale(0.75)
	thanksText.SetCenter(ui.Point{
		X: 0.99,
		Y: 0.85,
	})
	thanksG.Add(thanksText)

	thanksText = ui.NewElement()
	thanksText.SetText("Thank you!")
	thanksText.SetTextAlign(etxt.Center)
	thanksText.SetTextColor(ui.CornerTrackColor())
	thanksText.SetTextBold(true)
	thanksText.SetTextScale(1)
	thanksText.SetCenter(ui.Point{
		X: 0.95,
		Y: 0.95,
	})
	thanksG.Add(thanksText)

	// Set up beat triggers
	bmager := beats.NewManager(125, audio.GetBGMPositionMS())
	for i := 0; i < 4; i++ {
		bmager.SetTrigger(beats.BeatPosition{Numerator: i, Denominator: 4}, func() {
			state.logo.centerAnim.Pulse()
		})
		if i == 0 {
			bmager.SetTrigger(beats.BeatPosition{Numerator: i, Denominator: 4}, func() {
				state.anim.Pulse()
				c := state.namesOne.GetTextColor()
				state.namesOne.SetTextColor(state.namesTwo.GetTextColor())
				state.namesTwo.SetTextColor(c)
			})
		}
		if i%2 == 0 {
			bmager.SetTrigger(beats.BeatPosition{Numerator: i, Denominator: 4}, func() {
				state.logo.cornerAnim.Pulse()
			})
		}
	}
	state.bmager = bmager
	state.thanks = thanksG
	return &state
}

func (s *Title) Update() error {
	s.BaseGameState.Update()
	s.bmager.Update(audio.GetBGMPositionMS())
	s.group.Update()
	s.logo.Update()
	s.thanks.Update()
	s.anim.Update()

	scale := s.anim.GetScale()
	s.namesOne.SetRenderTextScale(scale)
	s.namesTwo.SetRenderTextScale(scale)

	if input.K.Is(ebiten.KeyEscape, input.JustPressed) || input.K.Is(ebiten.KeyF1, input.JustPressed)  {
		return ebiten.Termination
	}
	return nil
}

func (s *Title) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	s.group.Draw(screen, opts)
	s.logo.Draw(screen, opts)
	s.thanks.Draw(screen, opts)
}
