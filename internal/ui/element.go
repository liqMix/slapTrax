package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/l"
	"github.com/tinne26/etxt"
)

type Element struct {
	Component

	image *ebiten.Image
	text  string

	imageScale float64
	imageWidth float64

	textOptions      *TextOptions
	textBold         bool
	invertHoverColor bool
}

func NewElement() *Element {
	textOpts := GetDefaultTextOptions()
	e := &Element{
		Component: Component{
			center: &Point{},
			size:   &Point{},
		},
		textOptions: textOpts,
	}
	return e
}

func (e *Element) SetText(text string) {
	if text == "" {
		text = l.String(l.UNKNOWN)
	}
	e.text = text
}

func (e *Element) SetTextScale(scale float64) {
	e.textOptions.Scale = scale
}

func (e *Element) SetTextColor(color color.RGBA) {
	e.textOptions.Color = color
}

func (e *Element) InvertHoverColor() {
	e.invertHoverColor = true
}

func (e *Element) GetText() string {
	return e.text
}

func (e *Element) SetTextBold(b bool) {
	e.textBold = b
}

func (e *Element) SetImage(img *ebiten.Image) {
	if img == nil {
		return
	}

	width, _ := e.GetSize().ToRender()
	scale := width / float64(img.Bounds().Dx())
	e.image = img
	e.imageScale = scale
	e.imageWidth = PointFromRender(float64(img.Bounds().Dx()), 0).X
}

const hoverScale = 1.25

func (e *Element) DrawMarkers(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	if !e.hovered {
		return
	}

	textOpts := &TextOptions{
		Align: etxt.Center,
		Scale: e.textOptions.Scale * hoverScale,
		Color: CornerTrackColor(),
	}
	if e.invertHoverColor {
		textOpts.Color = CenterTrackColor()
	}

	width := TextWidth(textOpts, hoverMarkerLeft)
	if e.image != nil {
		width += e.imageWidth * e.imageScale
	} else {
		width += TextWidth(textOpts, e.text)
	}

	DrawHoverMarkersCenteredAt(screen, e.center, &Point{X: width, Y: 0}, textOpts, opts)
}

func (e *Element) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	if e.hidden || (e.image == nil && len(e.text) == 0) {
		return
	}

	center := e.GetCenter()

	if e.image != nil {
		DrawImageAt(screen, e.image, center, e.imageScale, opts)
	} else {
		textScale := e.textOptions.Scale
		textClr := e.textOptions.Color
		if e.hovered {
			textScale *= hoverScale
		}

		// Plain draw
		textOpts := &TextOptions{
			Align: etxt.Center,
			Scale: textScale,
			Color: textClr,
		}
		DrawTextAt(screen, e.text, center, textOpts, opts)
		if e.textBold {
			DrawTextAt(screen, e.text, &Point{X: center.X + 0.001, Y: center.Y}, textOpts, opts)
		}
	}
	e.DrawMarkers(screen, opts)
}
