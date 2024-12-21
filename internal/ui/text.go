package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tinne26/etxt"
	"golang.org/x/image/font/sfnt"
)

var text = initRenderer()

func initRenderer() *etxt.Renderer {
	r := etxt.NewRenderer()
	r.SetAlign(etxt.Center)
	r.Utils().SetCache8MiB()
	return r
}

func DrawTextAt(screen *ebiten.Image, s string, x, y int, scale float64) {
	text.SetScale(scale)
	text.Draw(screen, s, x, y)
}

func SetFont(f *sfnt.Font) {
	text.SetFont(f)
}
