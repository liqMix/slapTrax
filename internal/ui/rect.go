package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
)

func DrawFilledRect(screen *ebiten.Image, center *Point, size *Point, color color.RGBA) {
	if center == nil || size == nil {
		return
	}
	x, y := center.ToRender32()
	w, h := size.ToRender32()

	x = x - w/2
	y = y - h/2
	vector.DrawFilledRect(screen, x, y, w, h, color, true)
}

func DrawBorderedFilledRect(screen *ebiten.Image, center *Point, size *Point, color, borderColor color.RGBA, borderSize float64) {
	if center == nil || size == nil {
		return
	}
	x, y := center.ToRender32()
	w, h := size.ToRender32()
	bSize, _ := (&Point{X: borderSize, Y: borderSize}).ToRender32()

	bW := w + bSize
	bH := h + bSize

	x = x - w/2
	y = y - h/2
	bX := x - bSize/2
	bY := y - bSize/2

	vector.DrawFilledRect(screen, bX, bY, bW, bH, borderColor, true)
	vector.DrawFilledRect(screen, x, y, w, h, color, true)
}

const noteThemedRectBorderSize = 0.01

func DrawNoteThemedRect(screen *ebiten.Image, center *Point, size *Point) {
	yBorderSize := &Point{X: size.X, Y: size.Y + noteThemedRectBorderSize*2}
	xBorderSize := &Point{X: size.X + noteThemedRectBorderSize*2, Y: yBorderSize.Y}

	yBorderColor := CenterTrackColor()
	xBorderColor := CornerTrackColor()
	DrawFilledRect(screen, center, xBorderSize, xBorderColor)
	DrawFilledRect(screen, center, yBorderSize, yBorderColor)

	DrawFilledRect(screen, center, size, types.Black.C())
}
