package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tinne26/etxt"
	"golang.org/x/image/font/sfnt"
)

var text = initRenderer()

const defaultAlign = etxt.Center
const defaultScale = 1.0

func initRenderer() *etxt.Renderer {
	r := etxt.NewRenderer()
	r.Utils().SetCache8MiB()

	resetRenderer(r)
	return r
}

func resetRenderer(t *etxt.Renderer) {
	t.SetScale(defaultScale)
	t.SetAlign(defaultAlign)
}

func DrawTextCenterAt(screen *ebiten.Image, s string, x, y int, scale float64) {
	text.SetScale(scale)
	text.Draw(screen, s, x, y)

	// Reset
	resetRenderer(text)
}

func DrawTextRightAt(screen *ebiten.Image, s string, x, y int, scale float64) {
	text.SetScale(scale)
	text.SetAlign(etxt.Right)
	text.Draw(screen, s, x, y)

	// Reset
	resetRenderer(text)
}

func SetFont(f *sfnt.Font) {
	text.SetFont(f)
}
