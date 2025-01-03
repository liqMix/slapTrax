package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/locale"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/tinne26/etxt"
)

type Element struct {
	Component

	image     *ebiten.Image
	text      string
	textColor color.RGBA
	textBold  bool
	scale     float64
}

func NewElement() *Element {
	e := Element{scale: 1.0, textColor: types.White}
	return &e
}

func (e *Element) SetText(text string) {
	if text == "" {
		text = locale.String(types.L_UNKNOWN)
	}
	e.text = text
	w := TextWidth(text)
	h := TextHeight()
	e.SetSize(Point{float64(w), float64(h)})
}

func (e *Element) SetTextColor(color color.RGBA) {
	e.textColor = color
}

func (e *Element) GetText() string {
	return e.text
}

func (e *Element) SetTextBold(b bool) {
	e.textBold = b
}

func (e *Element) SetImage(img *ebiten.Image) {
	e.image = img
	if img != nil {
		w, h := img.Bounds().Dx(), img.Bounds().Dy()

		size := e.GetSize()
		if size == nil {
			e.SetSize(Point{float64(w), float64(h)})
		} else {
			// Scale image to fit size
			eW, eH := size.ToRender()
			scaleW := eW / float64(w)
			scaleH := eH / float64(h)
			if scaleW < scaleH {
				e.SetScale(scaleW)
			} else {
				e.SetScale(scaleH)
			}
		}
	}
}

func (e *Element) SetScale(scale float64) {
	e.scale = scale
}

func (e *Element) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	t := e.text
	if e.image == nil && len(e.text) == 0 {
		t = locale.String(types.L_UNKNOWN)
	}

	center := e.GetCenter()
	if center == nil {
		return
	}

	scale := e.scale
	if e.IsHovered() {
		scale = 1.5 * scale
	}

	if e.image != nil {
		DrawImageAt(screen, e.image, center, scale, opts)
	}

	if len(t) > 0 && e.scale > 0 {
		if e.hovered {
			t = "> " + t + " <"
		}
		clr := e.textColor
		if e.disabled {
			clr = color.RGBA{
				R: clr.R / 2,
				G: clr.G / 2,
				B: clr.B / 2,
				A: clr.A,
			}
		}
		DrawTextAt(screen, t, center, &TextOptions{
			Align: etxt.Center,
			Scale: scale,
			Color: clr,
		})
		if e.textBold {
			DrawTextAt(screen, t, &Point{X: center.X + 0.001, Y: center.Y}, &TextOptions{
				Align: etxt.Center,
				Scale: scale,
				Color: clr,
			})
		}
	}
}
