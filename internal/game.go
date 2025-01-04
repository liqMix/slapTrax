package internal

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/liqmix/ebiten-holiday-2024/internal/assets"
	"github.com/liqmix/ebiten-holiday-2024/internal/cache"
	"github.com/liqmix/ebiten-holiday-2024/internal/input"
	"github.com/liqmix/ebiten-holiday-2024/internal/logger"
	"github.com/liqmix/ebiten-holiday-2024/internal/render"
	"github.com/liqmix/ebiten-holiday-2024/internal/state"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
)

func getStartingState() *RenderState {
	// startingState := types.GameStatePlay
	// startingArgs := &state.PlayArgs{
	// 	Song:       resource.GetSongByTitle("another"),
	// 	Difficulty: 7,
	// }
	// startingState := types.GameStateOffset
	startingState := types.GameStateTitle
	startingArgs := interface{}(nil)
	return getState(startingState, startingArgs)
}

type RenderState struct {
	state    state.State
	renderer types.Renderer
}

func (r *RenderState) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	if r.renderer != nil {
		r.renderer.Draw(screen, opts)
	} else {
		r.state.Draw(screen, opts)
	}
}

type Game struct {
	debugLog     *ui.DebugLog
	renderScale  float64
	offsetX      float64
	offsetY      float64
	stateStack   []*RenderState
	stackImages  []*ebiten.Image
	currentState *RenderState

	// debug
	updateTimes []time.Duration
	drawTimes   []time.Duration
}

func getState(gs types.GameState, arg interface{}) *RenderState {
	state := state.New(gs, arg)
	return &RenderState{
		state:    state,
		renderer: render.GetRenderer(gs, state),
	}
}

func NewGame() *Game {
	return &Game{
		currentState: getStartingState(),
		debugLog:     ui.NewDebugLog(),
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	panic("nope")
}

func (g *Game) LayoutF(displayWidth, displayHeight float64) (float64, float64) {
	x, y := types.Window.RenderSize()

	// Calculate scale based on window vs canvas ratio
	scaleX := displayWidth / float64(x)
	scaleY := displayHeight / float64(y)
	prev := g.renderScale
	g.renderScale = math.Min(scaleX, scaleY)

	if prev != g.renderScale {
		// Clear cache
		cache.ClearImageCache()
	}
	// Calculate centering offsets
	g.offsetX = (displayWidth - float64(x)*g.renderScale) / 2
	g.offsetY = (displayHeight - float64(y)*g.renderScale) / 2

	// Set window offsets and scale
	types.Window.SetOffset(g.offsetX, g.offsetY)
	types.Window.SetScale(g.renderScale)
	types.Window.SetRenderSize(x, y)
	types.Window.SetDisplaySize(displayWidth, displayHeight)

	return displayWidth, displayHeight
}

func handleGlobalKeybinds() {
	// Global keybinds
	if input.K.Is(ebiten.KeyF2, input.JustPressed) {
		fmt.Println("Toggling debug")
		logger.ToggleDebug()
	}
}

func (g *Game) Update() error {
	start := time.Now()
	assets.Update()
	input.Update()
	handleGlobalKeybinds()

	if logger.IsDebugEnabled() {
		g.debugLog.Update()
		if input.M.Is(ebiten.MouseButtonLeft, input.JustPressed) {
			clicks++
		} else if input.M.Is(ebiten.MouseButtonLeft, input.JustReleased) {
			releases++
		}
	}

	if g.currentState != nil {
		gs := g.currentState.state
		gs.Update()
		if gs.HasNextState() {
			nextState, nextArgs := gs.GetNextState()
			gs.SetNextState(types.GameStateNone, nil)

			if nextState == types.GameStateBack {
				if len(g.stateStack) == 0 {
					// uh oh, exit?
					return nil
				}
				g.currentState = g.stateStack[len(g.stateStack)-1]
				g.stateStack = g.stateStack[:len(g.stateStack)-1]
				g.stackImages = g.stackImages[:len(g.stackImages)-1]
			} else {
				next := getState(nextState, nextArgs)
				if next.state.Floats() {
					// Freeze current state render to image
					img := types.NewRenderImage()
					opts := &ebiten.DrawImageOptions{}
					a := float32(0.25)
					opts.ColorScale.Scale(a, a, a, a)
					g.currentState.Draw(img, opts)
					g.stateStack = append(g.stateStack, g.currentState)
					g.stackImages = append(g.stackImages, img)
				} else {
					g.stateStack = nil
					g.stackImages = nil
				}
				g.currentState = next
			}
		}
	}

	g.updateTimes = append(g.updateTimes, time.Since(start))
	if len(g.updateTimes) > 60 {
		g.updateTimes = g.updateTimes[1:]
	}
	return nil
}

// Create canvas at render size
func (g *Game) GetCanvasImage() *ebiten.Image {
	canvas, ok := cache.GetImage("canvas")
	if !ok {
		canvas = types.NewRenderImage()
		cache.SetImage("canvas", canvas)
	}
	canvas.Clear()
	return canvas
}

func (g *Game) Draw(screen *ebiten.Image) {
	start := time.Now()
	screen.Clear()
	screen.Fill(color.Black)

	// Create transform for centered rendering
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(g.renderScale, g.renderScale)
	op.GeoM.Translate(g.offsetX, g.offsetY)
	canvas := g.GetCanvasImage()

	if g.currentState != nil {
		for _, img := range g.stackImages {
			// Draw previous states dimmed out
			canvas.DrawImage(img, nil)
		}
		g.currentState.Draw(canvas, nil)
	}

	if logger.IsDebugEnabled() {
		g.DrawDebug(canvas)
	}
	screen.DrawImage(canvas, op)

	g.drawTimes = append(g.drawTimes, time.Since(start))
	if len(g.drawTimes) > 60 {
		g.drawTimes = g.drawTimes[1:]
	}
}

var clicks = 0
var releases = 0

// Draw debug information at top left
func (g *Game) DrawDebug(screen *ebiten.Image) {
	// Draw FPS
	rX, rY := types.Window.RenderSize()
	wX, wY := ebiten.WindowSize()
	offset := 20
	y := offset
	debugPrints := []string{
		fmt.Sprintf("Render Size: %v, %v", rX, rY),
		fmt.Sprintf("Window Size: %v, %v", wX, wY),
		fmt.Sprintf("Scale: %.2f", g.renderScale),
		fmt.Sprintf("FPS: %.2f", ebiten.ActualFPS()),
		fmt.Sprintf("TPS: %.2f", ebiten.ActualTPS()),
	}
	for _, s := range debugPrints {
		ebitenutil.DebugPrintAt(screen, s, offset, y)
		y += offset
	}

	// Mouse buttons
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Clicks: %d", clicks), offset, y)
	y += offset
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Releases: %d", releases), offset, y)
	y += offset * 2

	// Draw mouse position
	mX, mY := ebiten.CursorPosition()
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Absolute Pos: %v, %v", mX, mY), offset, y)
	y += offset

	oX, oY := input.M.Position()
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Offset Pos: %.2f, %.2f", oX, oY), offset, y)
	y += offset

	nX, nY := ui.PointFromRender(oX, oY).V()
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Normalized Pos: %.2f, %.2f", nX, nY), offset, y)

	y += offset * 2

	// Draw average update and draw times
	avgUpdate := time.Duration(0)
	avgDraw := time.Duration(0)
	for _, t := range g.updateTimes {
		avgUpdate += t
	}
	for _, t := range g.drawTimes {
		avgDraw += t
	}
	if (len(g.updateTimes) != 0) && (len(g.drawTimes) != 0) {
		avgUpdate /= time.Duration(len(g.updateTimes))
		avgDraw /= time.Duration(len(g.drawTimes))
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Avg Update: %s", avgUpdate), offset, y)
		y += offset
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Avg Draw: %s", avgDraw), offset, y)
		y += offset * 2

	}

	// Draw pressed keys
	pressed := inpututil.AppendPressedKeys([]ebiten.Key{})
	ebitenutil.DebugPrintAt(screen, "Pressed keys:", offset, y)
	for i, key := range pressed {
		ebitenutil.DebugPrintAt(screen, key.String(), offset, y+((i+1)*offset))
	}
	// Draw debug log
	g.debugLog.Draw(screen)
}
