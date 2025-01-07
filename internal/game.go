package internal

import (
	"errors"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/audio"
	"github.com/liqmix/ebiten-holiday-2024/internal/debug"
	"github.com/liqmix/ebiten-holiday-2024/internal/display"
	"github.com/liqmix/ebiten-holiday-2024/internal/input"
	"github.com/liqmix/ebiten-holiday-2024/internal/logger"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
)

const (
	maxStateStackSize = 10
)

func getStartingState() *RenderState {
	// startingState := types.GameStatePlay
	// startingArgs := &state.PlayArgs{
	// 	Song:       assert,
	// 	Difficulty: 7,
	// }
	// startingState := types.GameStateOffset
	startingState := types.GameStateTitle
	startingArgs := interface{}(nil)
	return getState(startingState, startingArgs)
}

type Game struct {
	debugster *debug.Debugster

	stateStack   []*RenderState
	currentState *RenderState
}

func NewGame() *Game {
	return &Game{
		currentState: getStartingState(),
		debugster:    debug.NewDebugster(),
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
		return g.popState()
	}

	next := getState(nextState, nextArgs)
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
	audio.Update()
	input.Update()

	if input.K.Is(ebiten.KeyF2, input.JustPressed) {
		g.debugster.Toggle()
	}

	if g.currentState == nil {
		return nil
	}

	gs := g.currentState.state
	if err := gs.Update(); err != nil {
		return fmt.Errorf("state update failed: %w", err)
	}

	if gs.HasNextState() {
		nextState, nextArgs := gs.GetNextState()
		gs.SetNextState(types.GameStateNone, nil)

		if err := g.handleStateTransition(nextState, nextArgs); err != nil {
			return fmt.Errorf("state transition failed: %w", err)
		}
	}

	g.debugster.Update()
	return nil
}

// Create canvas at render size
func (g *Game) GetCanvasImage() *ebiten.Image {
	canvas, ok := display.GetCachedImage("canvas")
	if !ok {
		canvas = display.NewRenderImage()
		display.SetCachedImage("canvas", canvas)
	}
	canvas.Clear()
	return canvas
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	screen.Fill(color.Black)
	canvas := g.GetCanvasImage()

	if g.currentState != nil {
		for i, s := range g.stateStack {
			opts := &ebiten.DrawImageOptions{}
			a := float32(0.25) * float32(i+1)
			opts.ColorScale.Scale(a, a, a, a)
			s.Draw(canvas, opts)
		}
		g.currentState.Draw(canvas, nil)
	}

	screen.DrawImage(canvas, display.Window.GetScreenDrawOptions())

	if logger.IsDebugEnabled() {
		g.debugster.Draw(screen)
	}
}
