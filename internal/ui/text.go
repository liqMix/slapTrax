package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/assets"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/tinne26/etxt"
)

var text *etxt.Renderer

func TextHeight() float64 {
	_, y := types.Window.RenderSize()
	h := text.Measure(" ").Height().ToFloat64()
	return h / float64(y)
}

func TextWidth(s string) float64 {
	x, _ := types.Window.RenderSize()
	w := text.Measure(s).Width().ToFloat64()
	return w / float64(x)
}

type TextOptions struct {
	Align etxt.Align
	Scale float64
	Color color.RGBA
}

var DefaultOptions = TextOptions{
	Align: etxt.Center,
	Scale: 1,
	Color: types.White,
}

func InitTextRenderer() {
	opts := &DefaultOptions

	r := etxt.NewRenderer()
	r.Utils().SetCache8MiB()
	r.SetFont(assets.Font())
	r.SetAlign(opts.Align)
	r.SetScale(opts.Scale)
	r.SetColor(opts.Color)
	text = r
}

func resetRenderer() {
	setRenderer(&DefaultOptions)
}

func setRenderer(opts *TextOptions) {
	if opts == nil {
		return
	}
	text.SetAlign(opts.Align)
	text.SetScale(opts.Scale)
	text.SetColor(opts.Color)
}

func DrawTextAt(screen *ebiten.Image, txt string, center *Point, opts *TextOptions) {
	if center == nil {
		return
	}
	if opts == nil {
		opts = &DefaultOptions
	}

	x, y := center.ToRender()

	setRenderer(opts)
	text.Draw(screen, txt, int(x), int(y))
	resetRenderer()
}

func DrawTextBlockAt(screen *ebiten.Image, s []string, p *Point, opts *TextOptions) {
	if p == nil || len(s) == 0 {
		return
	}
	if opts == nil {
		opts = &DefaultOptions
	}

	setRenderer(opts)

	h := text.Measure(s[0]).IntHeight()

	x, y := p.ToRender()
	for i, line := range s {
		text.Draw(screen, line, int(x), int(y)+h*i)
	}

	resetRenderer()
}
