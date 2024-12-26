package internal

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/liqmix/ebiten-holiday-2024/internal/config"
	"github.com/liqmix/ebiten-holiday-2024/internal/render"
	"github.com/liqmix/ebiten-holiday-2024/internal/state"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
)

type RenderState struct {
	state.State
	Renderer render.IRenderer
}

type Game struct {
	renderScale  float64
	offsetX      float64
	offsetY      float64
	stateStack   []*RenderState
	currentState *RenderState
}

func getState(gs types.GameState, arg interface{}) *RenderState {
	state := state.New(gs, arg)
	return &RenderState{
		State:    state,
		Renderer: render.GetRenderer(gs, state),
	}

}
func NewGame() *Game {
	return &Game{
		currentState: getState(types.GameStateTitle, nil),
		stateStack:   []*RenderState{},
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	panic("nope")
}
func (g *Game) LayoutF(logicWinWidth, logicWinHeight float64) (float64, float64) {
	// Calculate scale based on window vs canvas ratio
	scaleX := logicWinWidth / float64(config.CANVAS_WIDTH)
	scaleY := logicWinHeight / float64(config.CANVAS_HEIGHT)
	g.renderScale = math.Min(scaleX, scaleY)

	// Calculate centering offsets
	g.offsetX = (logicWinWidth - float64(config.CANVAS_WIDTH)*g.renderScale) / 2
	g.offsetY = (logicWinHeight - float64(config.CANVAS_HEIGHT)*g.renderScale) / 2

	return logicWinWidth, logicWinHeight
}

// TODO: some sort of time step
func (g *Game) Update() error {
	// audio.Update()

	if g.currentState != nil {
		nextState, arg, err := g.currentState.Update()
		if err != nil {
			return err
		}
		if nextState != nil {
			g.stateStack = append(g.stateStack, g.currentState)
			g.currentState = getState(*nextState, arg)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()

	// Create transform for centered rendering
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(g.renderScale, g.renderScale)
	op.GeoM.Translate(g.offsetX, g.offsetY)

	// Create canvas at base resolution
	canvas := ebiten.NewImage(config.CANVAS_WIDTH, config.CANVAS_HEIGHT)

	if g.currentState != nil && g.currentState.Renderer != nil {
		g.currentState.Renderer.Draw(canvas)
	}

	if config.DEBUG {
		g.DrawDebug(canvas)
	}
	screen.DrawImage(canvas, op)
}

// Draw debug information at top left
func (g *Game) DrawDebug(screen *ebiten.Image) {
	// Draw FPS
	y := 0
	offset := 20
	debugPrints := []string{
		fmt.Sprintf("FPS: %.2f", ebiten.ActualFPS()),
		fmt.Sprintf("TPS: %.2f", ebiten.ActualTPS()),
	}
	for i, s := range debugPrints {
		ebitenutil.DebugPrintAt(screen, s, offset, y+((i+1)*offset))
	}
}
