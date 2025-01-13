package state

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/audio"
	"github.com/liqmix/slaptrax/internal/input"
	"github.com/liqmix/slaptrax/internal/l"
	"github.com/liqmix/slaptrax/internal/types"
	"github.com/liqmix/slaptrax/internal/ui"
	"github.com/liqmix/slaptrax/internal/user"
	"github.com/tinne26/etxt"
)

var (
	titleCenter = ui.Point{
		X: 0.5,
		Y: 0.25,
	}
	defaultCenterY = 0.45
	reducedCenterY = 0.65
	textCenterX    = 0.33
	imageCenterX   = 0.63

	configTextOpts = &ui.TextOptions{
		Align: etxt.Center,
		Scale: 1.5,
		Color: types.White.C(),
	}
	selectedTextOpts = &ui.TextOptions{
		Align: etxt.Center,
		Scale: 2.5,
		Color: ui.CornerTrackColor(),
	}
	descTextOpts = &ui.TextOptions{
		Align: etxt.Center,
		Scale: 1,
		Color: types.White.C(),
	}
)

type KeyConfigState struct {
	types.BaseGameState

	cb      func()
	title   *ui.Element
	current input.TrackKeyConfig
	group   *ui.UIGroup
}

func NewKeyConfigState(args *FloatStateArgs) *KeyConfigState {
	kcs := KeyConfigState{}
	if args != nil && args.Cb != nil {
		kcs.cb = args.Cb
	}

	title := ui.NewElement()
	title.SetSize(ui.Point{X: 0.5, Y: 0.1})
	title.SetCenter(titleCenter)
	title.SetText(l.String(l.SETTINGS_GAME_KEY_CONFIG))
	title.SetTextBold(true)
	title.SetTextScale(2)
	kcs.title = title

	g := ui.NewUIGroup()
	g.SetPaneled(true)
	g.SetCenter(ui.Point{X: 0.5, Y: 0.5})
	g.SetSize(ui.Point{X: 0.6, Y: 0.6})
	g.SetTriggerOnHover(true)

	// Default
	e := ui.NewElement()
	e.SetSize(ui.Point{X: 0.3, Y: 0.1})
	e.SetCenter(ui.Point{X: imageCenterX, Y: defaultCenterY})
	e.SetImage(input.TrackKeyConfigDefault.Image())
	e.SetTrigger(func() {
		kcs.current = input.TrackKeyConfigDefault
	})
	def := e

	// Reduced
	e = ui.NewElement()
	e.SetSize(ui.Point{X: 0.25, Y: 0.1})
	e.SetCenter(ui.Point{X: imageCenterX, Y: reducedCenterY})
	e.SetImage(input.TrackKeyConfigReduced.Image())
	e.SetTrigger(func() {
		kcs.current = input.TrackKeyConfigReduced
	})
	red := e
	if user.S().KeyConfig == 0 {
		g.Add(def)
		g.Add(red)
	} else {
		g.Add(red)
		g.Add(def)
	}

	kcs.group = g
	return &kcs
}

func (kcs *KeyConfigState) SaveAndExit() {
	user.S().KeyConfig = int(kcs.current)
	user.Save()
	input.SetTrackKeys(kcs.current)
	kcs.Exit()
}

func (kcs *KeyConfigState) Exit() {
	if kcs.cb != nil {
		kcs.cb()
	}
	audio.PlaySFX(audio.SFXBack)
	kcs.SetNextState(types.GameStateBack, nil)
}
func (kcs *KeyConfigState) Update() error {
	kcs.BaseGameState.Update()

	kcs.group.Update()

	if input.JustActioned(input.ActionBack) {
		kcs.Exit()
	} else if input.JustActioned(input.ActionSelect) {
		kcs.SaveAndExit()
	}
	return nil
}

func (kcs *KeyConfigState) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	kcs.group.Draw(screen, opts)
	kcs.title.Draw(screen, opts)

	if kcs.current == 0 {
		ui.DrawTextAt(screen, l.String(l.KEY_CONFIG_DEFAULT), &ui.Point{X: textCenterX, Y: defaultCenterY}, selectedTextOpts, opts)
		ui.DrawTextAt(screen, l.String(l.KEY_CONFIG_REDUCED), &ui.Point{X: textCenterX, Y: reducedCenterY}, configTextOpts, opts)
	} else {
		ui.DrawTextAt(screen, l.String(l.KEY_CONFIG_DEFAULT), &ui.Point{X: textCenterX, Y: defaultCenterY}, configTextOpts, opts)
		ui.DrawTextAt(screen, l.String(l.KEY_CONFIG_REDUCED), &ui.Point{X: textCenterX, Y: reducedCenterY}, selectedTextOpts, opts)
	}

	ui.DrawTextAt(screen, l.String(l.KEY_CONFIG_DEFAULT_DESC), &ui.Point{X: textCenterX, Y: defaultCenterY + 0.05}, descTextOpts, opts)
	ui.DrawTextAt(screen, l.String(l.KEY_CONFIG_REDUCED_DESC), &ui.Point{X: textCenterX, Y: reducedCenterY + 0.05}, descTextOpts, opts)

}
