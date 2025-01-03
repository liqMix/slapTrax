package types

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/cache"
)

// Allowed render sizes
type RenderSize string

var OffsetX float64
var OffsetY float64
var RenderScale float64

const (
	RenderSizeTiny   RenderSize = L_RENDERSIZE_TINY
	RenderSizeSmall  RenderSize = L_RENDERSIZE_SMALL
	RenderSizeMedium RenderSize = L_RENDERSIZE_MEDIUM
	RenderSizeLarge  RenderSize = L_RENDERSIZE_LARGE
	RenderSizeMax    RenderSize = L_RENDERSIZE_MAX
)

func (r RenderSize) Value() (int, int) {
	switch r {
	case RenderSizeTiny:
		return 640, 360
	case RenderSizeSmall:
		return 960, 540
	case RenderSizeMedium:
		return 1280, 720
	case RenderSizeLarge:
		return 1920, 1080
	// Max allowed by screen
	case RenderSizeMax:
		w, _ := ebiten.Monitor().Size()
		// 16:9 aspect ratio
		return w, w * 9 / 16
	}
	return 1280, 720
}

var Window window = window{
	displayWidth:  1280,
	displayHeight: 720,

	renderWidth:  1280,
	renderHeight: 720,

	offsetX: 0,
	offsetY: 0,

	renderScale: 1,
}

type window struct {
	displayWidth  float64
	displayHeight float64

	renderWidth  int
	renderHeight int

	offsetX float64
	offsetY float64

	renderScale float64
}

func (w *window) SetOffset(x, y float64) {
	w.offsetX = x
	w.offsetY = y
}

func (w *window) SetScale(s float64) {
	w.renderScale = s
}

func (w *window) SetRenderSize(width, height int) {
	w.renderWidth = width
	w.renderHeight = height
	cache.ClearImageCache()
}
func (w *window) ScaleByRender(n float32) float32 {
	return float32(w.renderScale) * n
}

func (w *window) SetDisplaySize(width, height float64) {
	w.displayWidth = width
	w.displayHeight = height
}

func (w *window) Offset() (float64, float64) {
	return w.offsetX, w.offsetY
}

func (w *window) RenderScale() float64 {
	return w.renderScale
}

func (w *window) RenderSize() (int, int) {
	return w.renderWidth, w.renderHeight
}

func (w *window) DisplaySize() (float64, float64) {
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
