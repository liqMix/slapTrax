package debug

import (
	"fmt"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/liqmix/ebiten-holiday-2024/internal/display"
	"github.com/liqmix/ebiten-holiday-2024/internal/input"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
)

const maxTimes = 50

type Debugster struct {
	enabled bool

	debugLog    *ui.DebugLog
	updateTimes []time.Duration
	drawTimes   []time.Duration

	lastUpdateTime time.Time
	lastDrawTime   time.Time
}

func NewDebugster() *Debugster {
	return &Debugster{
		debugLog: ui.NewDebugLog(),
	}
}

func (d *Debugster) Toggle() {
	d.enabled = !d.enabled

	if !d.enabled {
		d.updateTimes = []time.Duration{}
		d.drawTimes = []time.Duration{}
	}
}

func (d *Debugster) Update() {
	if !d.enabled {
		return
	}

	d.debugLog.Update()

	if d.lastUpdateTime.IsZero() {
		d.lastUpdateTime = time.Now()
	} else {
		d.updateTimes = append(d.updateTimes, time.Since(d.lastUpdateTime))
		d.lastUpdateTime = time.Now()
		if len(d.updateTimes) > maxTimes {
			d.updateTimes = d.updateTimes[1:]
		}
	}
}

// Draw debug information at top left
func (d *Debugster) Draw(screen *ebiten.Image) {
	if !d.enabled {
		return
	}

	// Draw FPS
	rX, rY := display.Window.RenderSize()
	wX, wY := ebiten.WindowSize()
	offset := 15
	y := 300
	debugPrints := []string{
		fmt.Sprintf("Render Size: %v, %v", rX, rY),
		fmt.Sprintf("Window Size: %v, %v", wX, wY),
		fmt.Sprintf("Scale: %.2f", display.Window.RenderScale()),
		fmt.Sprintf("FPS: %.2f", ebiten.ActualFPS()),
		fmt.Sprintf("TPS: %.2f", ebiten.ActualTPS()),
	}
	for _, s := range debugPrints {
		ebitenutil.DebugPrintAt(screen, s, offset, y)
		y += offset
	}

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
	for _, t := range d.updateTimes {
		avgUpdate += t
	}
	for _, t := range d.drawTimes {
		avgDraw += t
	}
	if (len(d.updateTimes) != 0) && (len(d.drawTimes) != 0) {
		avgUpdate /= time.Duration(len(d.updateTimes))
		avgDraw /= time.Duration(len(d.drawTimes))
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
	d.debugLog.Draw(screen)

	if d.lastDrawTime.IsZero() {
		d.lastDrawTime = time.Now()
	} else {
		d.drawTimes = append(d.drawTimes, time.Since(d.lastDrawTime))
		d.lastDrawTime = time.Now()
		if len(d.drawTimes) > maxTimes {
			d.drawTimes = d.drawTimes[1:]
		}
	}
}
