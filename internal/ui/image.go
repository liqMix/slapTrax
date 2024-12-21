package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

func DrawImageAtRaw(screen *ebiten.Image, img *ebiten.Image, x, y int, scale float64) {
	if img == nil || screen == nil {
		fmt.Println("img or screen is nil")
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.GeoM.Scale(scale, scale)
	screen.DrawImage(img, op)
}

// Draws image centered on x, y
func DrawImageAt(screen *ebiten.Image, img *ebiten.Image, x, y int, scale float64) {
	if img == nil || screen == nil {
		fmt.Println("img or screen is nil")
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(x), float64(y))
	op.GeoM.Translate(-float64(img.Bounds().Dx())/2*scale, -float64(img.Bounds().Dy())/2*scale)
	screen.DrawImage(img, op)
}
