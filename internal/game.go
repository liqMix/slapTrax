package internal

import (
	"errors"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/assets"
	"github.com/liqmix/slaptrax/internal/audio"
	"github.com/liqmix/slaptrax/internal/cache"
	"github.com/liqmix/slaptrax/internal/debug"
	"github.com/liqmix/slaptrax/internal/display"
	"github.com/liqmix/slaptrax/internal/input"
	"github.com/liqmix/slaptrax/internal/logger"
	"github.com/liqmix/slaptrax/internal/types"
	"github.com/liqmix/slaptrax/internal/ui"
	"github.com/tinne26/etxt"
)

const (
	maxStateStackSize = 10
	startingTicks     = 200
)

func getStartingState() *RenderState {
	// startingState := types.GameStatePlay
	// song := types.GetAllSongs()[0]
	// diff := song.GetDifficulties()[0]
	// startingArgs := &state.PlayArgs{
	// Song:       song,
	// Difficulty: diff,
	// }

	// startingState := types.GameStateOffset
	startingState := types.GameStateTitle
	startingArgs := interface{}(nil)
	// startingState := types.GameStateResult
	// startingArgs := &state.ResultStateArgs{
	// 	Score: &types.Score{
	// 		TotalScore: types.MaxScore,
	// 		MaxCombo:   100,
	// 		Difficulty: 7,
	// 		Song:       types.GetAllSongs()[0],
	// 		Rating:     types.RatingS,
	// 		TotalNotes: 100,
	// 		Perfect:    1000, Good: 0, Bad: 0, Miss: 0,
	// 		Combo:      100,
	// 		HitRecords: []*types.HitRecord{},
	// 	},
	// }

	return GetState(startingState, startingArgs)
}

type Game struct {
	debugster *debug.Debugster

	started         bool
	startTicks      int64
	stateStack      []*RenderState
	currentState    *RenderState
	navText         *ui.NavText
	loadingTextOpts *ui.TextOptions
	userHeader      *ui.UserProfile
	background      *ebiten.Image
}

func NewGame() *Game {
	return &Game{
		started:      false,
		startTicks:   0,
		currentState: getStartingState(),
		debugster:    debug.NewDebugster(),
		navText:      ui.NewNavText(),
		userHeader:   ui.NewUserProfile(),
		loadingTextOpts: &ui.TextOptions{
			Align: etxt.Center,
			Scale: 1.5,
			Color: types.White.C(),
		},
		background: assets.GetImage("background.png"),
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	panic("nope")
}

func (g *Game) LayoutF(displayWidth, displayHeight float64) (float64, float64) {
	display.Window.SetDisplaySize(displayWidth, displayHeight)
	return displayWidth, displayHeight
}

func (g *Game) handleStateTransition(nextState types.GameState, nextArgs interface{}) error {
	if nextState == types.GameStateBack {
		audio.PlaySFX(audio.SFXBack)
		return g.popState()
	}

	next := GetState(nextState, nextArgs)
	if next.state.Floats() {
		if len(g.stateStack) >= maxStateStackSize {
			return fmt.Errorf("state stack overflow: max size %d reached", maxStateStackSize)
		}
		g.currentState.Freeze()
		g.stateStack = append(g.stateStack, g.currentState)
	} else {
		g.stateStack = nil
	}
	g.currentState = next
	return nil
}

func (g *Game) popState() error {
	if len(g.stateStack) == 0 {
		return errors.New("cannot pop state: stack is empty")
	}
	g.currentState = g.stateStack[len(g.stateStack)-1]
	g.stateStack = g.stateStack[:len(g.stateStack)-1]
	g.currentState.Unfreeze()
	return nil
}

func (g *Game) Update() error {
	if !g.started {
		g.startTicks++
		if g.startTicks >= startingTicks {
			g.started = true
		}
	}

	audio.Update()
	input.Update()

	if input.JustActioned(input.ActionToggleDebug) {
		g.debugster.Toggle()
	}

	if g.currentState == nil {
		return nil
	}

	gs := g.currentState.state
	if gs == nil {
		return nil
	}

	action := gs.CheckActions()
	if action != input.ActionUnknown {
		sfx := audio.ActionSFX(action)
		if sfx != audio.SFXNone {
			audio.PlaySFX(sfx)
		}
	}

	if err := gs.Update(); err != nil {
		return err
	}

	if gs.HasNextState() {
		nextState, nextArgs := gs.GetNextState()
		gs.SetNextState(types.GameStateNone, nil)

		if err := g.handleStateTransition(nextState, nextArgs); err != nil {
			return fmt.Errorf("state transition failed: %w", err)
		}
	}

	g.userHeader.Update()
	g.debugster.Update()
	return nil
}

// Create canvas at render size
func (g *Game) GetCanvasImage() *ebiten.Image {
	canvas, ok := cache.Image.Get("canvas")
	if !ok {
		canvas = display.NewRenderImage()
		cache.Image.Set("canvas", canvas)
	}
	canvas.Clear()
	return canvas
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	screen.Fill(color.Black)
	canvas := g.GetCanvasImage()

	bgOpts := &ebiten.DrawImageOptions{}
	
	// Scale background to fit render size
	renderWidth, renderHeight := display.Window.RenderSize()
	if g.background != nil {
		scaleX := float64(renderWidth) / float64(g.background.Bounds().Dx())
		scaleY := float64(renderHeight) / float64(g.background.Bounds().Dy())
		bgOpts.GeoM.Scale(scaleX, scaleY)
	}
	
	bgOpts.ColorScale.Scale(0.25, 0.25, 0.25, 0.25)
	canvas.DrawImage(g.background, bgOpts)

	if g.currentState != nil {
		for i, s := range g.stateStack {
			opts := &ebiten.DrawImageOptions{}
			a := float32(0.25) * float32(i+1)
			opts.ColorScale.Scale(a, a, a, a)
			s.Draw(canvas, opts)
		}
		g.currentState.Draw(canvas, nil)

		// Draw nav action bar if navigable
		if g.currentState.state.IsNavigable() {
			g.navText.Draw(canvas, nil)
		}
	}

	g.userHeader.Draw(canvas, nil)
	opts := display.Window.GetScreenDrawOptions()
	if !g.started {
		scale := float32(g.startTicks) / float32(startingTicks)
		opts.ColorScale.ScaleAlpha(scale)
	}
	screen.DrawImage(canvas, opts)

	if logger.IsDebugEnabled() {
		g.debugster.Draw(screen)
	}
}
