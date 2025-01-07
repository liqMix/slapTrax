package display

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/cache"
	"github.com/liqmix/ebiten-holiday-2024/internal/user"
)

var Window window

type window struct {
	displayWidth  float64
	displayHeight float64
	fullScreen    bool

	offsetX float64
	offsetY float64

	renderWidth      int
	renderHeight     int
	renderScale      float64
	fixedRenderScale bool
}

func InitWindow() {
	Window = window{
		displayWidth:  float64(user.S.ScreenWidth),
		displayHeight: float64(user.S.ScreenHeight),

		renderWidth:  user.S.RenderWidth,
		renderHeight: user.S.RenderHeight,

		fixedRenderScale: true,
	}
}

func (w *window) ClearCaches() {
	cache.Image.Clear(w.renderWidth, w.renderHeight)
	cache.Path.Clear(w.renderWidth, w.renderHeight)
}

func (w *window) Refresh() {
	// Update the offset and the render scale
	width, height := w.DisplaySize()
	x, y := w.RenderSize()

	w.offsetX = (width - float64(x)*w.renderScale) / 2
	w.offsetY = (height - float64(y)*w.renderScale) / 2

	scaleX := width / float64(x)
	scaleY := height / float64(y)
	renderScale := math.Min(scaleX, scaleY)
	w.SetRenderScale(renderScale)
}

func (w *window) SetDisplaySize(width, height float64) {
	w.displayWidth = width
	w.displayHeight = height
	w.Refresh()
}

func (w *window) SetFixedRenderScale(fixed bool) {
	w.fixedRenderScale = fixed
}

func (w *window) IsFixedRenderScale() bool {
	return w.fixedRenderScale
}

func (w *window) IsFullscreen() bool {
	return ebiten.IsFullscreen()
}

func (w *window) SetFullscreen(fullscreen bool) {
	ebiten.SetFullscreen(fullscreen)
	w.fullScreen = fullscreen
	w.Refresh()
}

func (w *window) GetMonitorSize() (float64, float64) {
	width, height := ebiten.Monitor().Size()
	return float64(width), float64(height)
}

func (w *window) GetScreenDrawOptions() *ebiten.DrawImageOptions {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(w.offsetX, w.offsetY)

	if !w.fixedRenderScale || !w.fullScreen {
		opts.GeoM.Scale(w.renderScale, w.renderScale)
	} else {
		width, height := w.DisplaySize()
		left := (width - float64(w.renderWidth)) / 2
		top := (height - float64(w.renderHeight)) / 2
		opts.GeoM.Translate(left, top)
	}
	return opts
}

func (w *window) SetRenderScale(s float64) {
	w.renderScale = s
}

func (w *window) SetRenderSize(width, height int) {
	if width == w.renderWidth && height == w.renderHeight {
		return
	}

	w.renderWidth = width
	w.renderHeight = height
	w.ClearCaches()
	w.SetDisplaySize(w.displayWidth, w.displayHeight)
}

func (w *window) ScaleByRender(n float32) float32 {
	return float32(w.renderScale) * n
}

// func (w *window) SetRefreshRate(rate int) {
// 	w.refreshRate = rate
// }

// func (w *window) RefreshRate() int {
// 	return w.refreshRate
// }

func (w *window) Offset() (float64, float64) {
	return w.offsetX, w.offsetY
}

func (w *window) RenderScale() float64 {
	return w.renderScale
}

func (w *window) RenderSize() (int, int) {
	width, height := w.DisplaySize()
	return int(math.Min(float64(w.renderWidth), width)), int(math.Min(float64(w.renderHeight), height))
}

func (w *window) DisplaySize() (float64, float64) {
	if w.IsFullscreen() {
		return w.GetMonitorSize()
	}
	return w.displayWidth, w.displayHeight
}

func (w *window) CanvasPosition(windowX, windowY float64) (float64, float64) {
	canvasX := (windowX - w.offsetX) / w.renderScale
	canvasY := (windowY - w.offsetY) / w.renderScale
	return canvasX, canvasY
}

func (w *window) WindowPosition(canvasX, canvasY float64) (float64, float64) {
	windowX := canvasX*w.renderScale + w.offsetX
	windowY := canvasY*w.renderScale + w.offsetY
	return windowX, windowY
}
