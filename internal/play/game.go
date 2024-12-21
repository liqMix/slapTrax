package play

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/config"
	"github.com/liqmix/ebiten-holiday-2024/internal/l"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
)

type Game struct {
	renderScale float64
	offsetX     float64
	offsetY     float64
}

func NewGame() *Game {
	return &Game{}
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

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Create transform for centered rendering
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(g.renderScale, g.renderScale)
	op.GeoM.Translate(g.offsetX, g.offsetY)

	// Create canvas at base resolution
	canvas := ebiten.NewImage(config.CANVAS_WIDTH, config.CANVAS_HEIGHT)

	// Draw everything to canvas at base resolution
	centerX := config.CANVAS_WIDTH / 2
	centerY := config.CANVAS_HEIGHT / 2

	ui.DrawTextAt(canvas, l.String(l.DEBUG_GREETING), centerX, centerY, config.FONT_SCALE)
	ui.DrawImageAt(canvas, l.Flag(), centerX, centerY-100, 5.0)

	// Draw scaled canvas to screen
	screen.DrawImage(canvas, op)
}
