package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/l"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/tinne26/etxt"
)

type Element struct {
	Component

	image          *ebiten.Image
	text           string
	textColor      color.RGBA
	textBold       bool
	scale          float64
	forceTextColor bool
}

func NewElement() *Element {
	e := Element{scale: 1.0, textColor: types.White.C()}
	return &e
}

func (e *Element) SetText(text string) {
	if text == "" {
		text = l.String(l.UNKNOWN)
	}
	e.text = text
	w := TextWidth(nil, text)
	h := TextHeight(nil)
	e.SetSize(Point{float64(w), float64(h)})
}

func (e *Element) SetTextColor(color color.RGBA) {
	e.textColor = color
}

func (e *Element) ForceTextColor() {
	e.forceTextColor = true
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

const hoveredMarkerLeft = "> "
const hoveredMarkerRight = " <"

func (e *Element) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	if e.hidden {
		return
	}
	t := e.text
	if e.image == nil && len(e.text) == 0 {
		t = l.String(l.UNKNOWN)
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
		textClr := e.textColor
		if e.disabled {
			textClr = color.RGBA{
				R: textClr.R / 2,
				G: textClr.G / 2,
				B: textClr.B / 2,
				A: textClr.A,
			}
		} else if e.hovered {
			hoverTextOpts := &TextOptions{
				Align: etxt.Center,
				Scale: scale,
				Color: CornerTrackColor(),
			}
			if !e.forceTextColor {
				textClr = CenterTrackColor()
			}
			width := TextWidth(nil, t)*scale + TextWidth(nil, hoveredMarkerLeft)

			DrawTextAt(
				screen,
				hoveredMarkerLeft,
				&Point{
					X: center.X - float64(width/2),
					Y: center.Y,
				},
				hoverTextOpts,
				opts,
			)
			DrawTextAt(
				screen,
				hoveredMarkerRight,
				&Point{
					X: center.X + float64(width/2),
					Y: center.Y,
				},
				hoverTextOpts,
				opts,
			)
		}

		textOpts := &TextOptions{
			Align: etxt.Center,
			Scale: scale,
			Color: textClr,
		}
		DrawTextAt(screen, t, center, textOpts, opts)
		if e.textBold {
			DrawTextAt(screen, t, &Point{X: center.X + 0.001, Y: center.Y}, textOpts, opts)
		}
	}
}
