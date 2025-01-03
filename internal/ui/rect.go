package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
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
	bW, bH := Point{X: size.X - borderSize, Y: size.Y - borderSize}.ToRender32()

	bX := x - bW/2
	bY := y - bH/2

	x = x - w/2
	y = y - h/2

	vector.DrawFilledRect(screen, bX, bY, bW, bH, borderColor, true)
	vector.DrawFilledRect(screen, x, y, w, h, color, true)
}
