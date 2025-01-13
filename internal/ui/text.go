package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/assets"
	"github.com/liqmix/ebiten-holiday-2024/internal/display"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/tinne26/etxt"
)

type TextOptions struct {
	Align etxt.Align
	Scale float64
	Color color.RGBA
}

var GetDefaultTextOptions = func() *TextOptions {
	return &TextOptions{
		Align: etxt.Center,
		Scale: 1,
		Color: types.White.C(),
	}
}

var defaultWidth, defaultHeight float64 = 1280, 720

func TextHeight(opts *TextOptions) float64 {
	if opts == nil {
		opts = GetDefaultTextOptions()
	}
	text := getTextRenderer(opts)
	_, y := display.Window.RenderSize()
	h := text.Measure(" ").Height().ToFloat64() * getRenderTextScale()
	return h / float64(y)
}

func TextWidth(opts *TextOptions, s string) float64 {
	if opts == nil {
		opts = GetDefaultTextOptions()
	}
	text := getTextRenderer(opts)
	x, _ := display.Window.RenderSize()
	w := text.Measure(s).Width().ToFloat64() * getRenderTextScale()
	return w / float64(x)
}

func getTextRenderer(opts *TextOptions) *etxt.Renderer {
	if opts == nil {
		opts = GetDefaultTextOptions()
	}
	r := etxt.NewRenderer()
	r.Utils().SetCache8MiB()
	r.SetFont(assets.Font())
	r.SetAlign(opts.Align)
	r.SetScale(opts.Scale)
	r.SetColor(opts.Color)
	return r
}

func saferDraw(txt *etxt.Renderer, screen *ebiten.Image, t string, x, y int) {
	defer func() {
		// Try again with default locale
		if r := recover(); r != nil {
			prev := txt.GetFont()
			defTxt, defFont := assets.GetDefaultLocaleString(t)
			txt.SetFont(defFont)
			txt.Draw(screen, defTxt, x, y)
			txt.SetFont(prev)
		}
	}()
	txt.Draw(screen, t, x, y)
}

func getRenderTextScale() float64 {
	renderWidth, _ := display.Window.RenderSize()
	return float64(renderWidth) / defaultWidth
}

func DrawTextAt(screen *ebiten.Image, txt string, center *Point, opts *TextOptions, screenOpts *ebiten.DrawImageOptions) {
	if center == nil {
		return
	}

	if opts == nil {
		opts = GetDefaultTextOptions()
	}
	color := opts.Color
	if screenOpts != nil {
		color = ApplyAlphaScale(color, screenOpts.ColorScale.A())
	}
	text := getTextRenderer(&TextOptions{
		Align: opts.Align,
		Scale: opts.Scale * getRenderTextScale(),
		Color: color,
	})

	x, y := center.ToRender()
	saferDraw(text, screen, txt, int(x), int(y))
}

func DrawTextBlockAt(screen *ebiten.Image, s []string, p *Point, opts *TextOptions, screenOpts *ebiten.DrawImageOptions) {
	if p == nil || len(s) == 0 {
		return
	}
	if opts == nil {
		opts = GetDefaultTextOptions()
	}

	color := opts.Color
	if screenOpts != nil {
		color = ApplyAlphaScale(color, screenOpts.ColorScale.A())
	}

	text := getTextRenderer(&TextOptions{
		Align: opts.Align,
		Scale: opts.Scale * getRenderTextScale(),
		Color: color,
	})

	x, y := p.ToRender()
	h := text.Measure(s[0]).IntHeight()
	totalHeight := float64(len(s) * h)
	y -= totalHeight / 2

	for i, line := range s {
		saferDraw(text, screen, line, int(x), int(y)+h*(i+1))
	}
}

const hoverMarkerLeft = "> "
const hoverMarkerRight = " <"

func DrawHoverMarkersCenteredAt(screen *ebiten.Image, center *Point, size *Point, opts *TextOptions, screenOpts *ebiten.DrawImageOptions) {
	if center == nil {
		return
	}
	if opts == nil {
		opts = GetDefaultTextOptions()
	}
	color := opts.Color
	if screenOpts != nil {
		color = ApplyAlphaScale(color, screenOpts.ColorScale.A())
	}
	text := getTextRenderer(&TextOptions{
		Align: opts.Align,
		Scale: opts.Scale * getRenderTextScale(),
		Color: color,
	})
	w, _ := size.ToRender()
	x, y := center.ToRender()

	// Draw the selection markers
	saferDraw(text, screen, hoverMarkerLeft, int(x-w/2), int(y))
	saferDraw(text, screen, hoverMarkerRight, int(x+w/2), int(y))
}
