package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func DrawImageAt(screen *ebiten.Image, img *ebiten.Image, center *Point, scale float64, opts *ebiten.DrawImageOptions) {
	if img == nil || screen == nil {
		return
	}
	x, y := center.ToRender()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(x), float64(y))
	op.GeoM.Translate(-float64(img.Bounds().Dx())/2*scale, -float64(img.Bounds().Dy())/2*scale)
	if opts != nil {
		op.ColorScale = opts.ColorScale
	}
	screen.DrawImage(img, op)
}
