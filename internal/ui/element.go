package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/cache"
	"github.com/liqmix/ebiten-holiday-2024/internal/l"
	"github.com/tinne26/etxt"
)

type Element struct {
	Component

	text           string
	baseImage      *ebiten.Image
	baseImageScale float64
	imageScale     float64
	imageWidth     float64

	textOptions      *TextOptions
	renderTextScale  float64
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
		textOptions:     textOpts,
		renderTextScale: 1,
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
func (e *Element) SetRenderTextScale(scale float64) {
	e.renderTextScale = scale
}

func (e *Element) SetTextColor(color color.RGBA) {
	e.textOptions.Color = color
}
func (e *Element) SetTextAlign(align etxt.Align) {
	e.textOptions.Align = align
}

func (e *Element) InvertHoverColor() {
	e.invertHoverColor = true
}

func (e *Element) GetText() string {
	return e.text
}

func (e *Element) GetImageSize() (int, int) {
	if e.baseImage == nil {
		return 0, 0
	}
	return e.baseImage.Bounds().Dx(), e.baseImage.Bounds().Dy()
}
func (e *Element) SetImageScale(scale float64) {
	if e.baseImage == nil {
		return
	}
	e.imageScale = e.baseImageScale * scale
}

func (e *Element) SetTextBold(b bool) {
	e.textBold = b
}

func (e *Element) refreshImage() *ebiten.Image {
	img := ebiten.NewImageFromImage(e.baseImage)
	width, _ := e.GetSize().ToRender()
	scale := width / float64(img.Bounds().Dx())

	e.baseImageScale = scale
	e.imageScale = scale
	e.imageWidth = PointFromRender(float64(img.Bounds().Dx()), 0).X
	cache.Image.Set(e.GetId(), img)
	return img
}

func (e *Element) SetImage(img *ebiten.Image) {
	if img == nil {
		return
	}
	e.baseImage = img
	e.refreshImage()
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
	if e.baseImage != nil {
		width += e.imageWidth * e.imageScale
	} else {
		width += TextWidth(textOpts, e.text)
	}

	DrawHoverMarkersCenteredAt(screen, e.center, &Point{X: width, Y: 0}, textOpts, opts)
}

func (e *Element) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	if e.hidden || (e.baseImage == nil && len(e.text) == 0) {
		return
	}

	center := e.GetCenter()

	if e.baseImage != nil {
		img, ok := cache.Image.Get(e.GetId())
		if !ok {
			img = e.refreshImage()
		}
		DrawImageAt(screen, img, center, e.imageScale, opts)
	} else {
		renderScale := e.renderTextScale
		if e.hovered {
			renderScale *= hoverScale
		}

		// Plain draw
		textOpts := &TextOptions{
			Align: e.textOptions.Align,
			Scale: e.textOptions.Scale * renderScale,
			Color: e.textOptions.Color,
		}
		DrawTextAt(screen, e.text, center, textOpts, opts)
		if e.textBold {
			DrawTextAt(screen, e.text, &Point{X: center.X + 0.001, Y: center.Y}, textOpts, opts)
		}
	}
	e.DrawMarkers(screen, opts)
}
