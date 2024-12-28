package internal

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/liqmix/ebiten-holiday-2024/internal/cache"
	"github.com/liqmix/ebiten-holiday-2024/internal/config"
	"github.com/liqmix/ebiten-holiday-2024/internal/render"
	"github.com/liqmix/ebiten-holiday-2024/internal/song"
	"github.com/liqmix/ebiten-holiday-2024/internal/state"
	"github.com/liqmix/ebiten-holiday-2024/internal/state/play"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/liqmix/ebiten-holiday-2024/internal/user"
)

type Position struct {
	X float64
	Y float64
}

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
		// currentState: getState(types.GameStateTitle, nil),
		currentState: getState(types.GameStatePlay, play.PlayArgs{
			Song:       song.GetTestSong(),
			Difficulty: 7,
		}),
		stateStack: []*RenderState{},
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	panic("nope")
}
func (g *Game) LayoutF(displayWidth, displayHeight float64) (float64, float64) {
	s := user.Settings()

	// Calculate scale based on window vs canvas ratio
	scaleX := displayWidth / float64(s.RenderWidth)
	scaleY := displayHeight / float64(s.RenderHeight)
	g.renderScale = math.Min(scaleX, scaleY)

	// Calculate centering offsets
	g.offsetX = (displayWidth - float64(s.RenderWidth)*g.renderScale) / 2
	g.offsetY = (displayHeight - float64(s.RenderHeight)*g.renderScale) / 2

	return displayWidth, displayHeight
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
	screen.Fill(color.Black)

	// Create transform for centered rendering
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(g.renderScale, g.renderScale)
	op.GeoM.Translate(g.offsetX, g.offsetY)

	// Create canvas at base resolution
	canvas, ok := cache.GetImage("canvas")
	if !ok {
		s := user.Settings()
		canvas = ebiten.NewImage(s.RenderWidth, s.RenderHeight)
		cache.SetImage("canvas", canvas)
	}
	canvas.Clear()
	canvas.Fill(color.Black)

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
