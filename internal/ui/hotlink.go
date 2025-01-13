package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/external"
	"github.com/liqmix/slaptrax/internal/types"
)

type HotLink struct {
	Element

	url           string
	underlinePath *VectorPath
}

func NewHotLink() *HotLink {
	h := HotLink{
		Element:       *NewElement(),
		underlinePath: nil,
	}

	h.SetTrigger(func() {
		external.OpenURL(h.GetURL())
	})

	return &h
}

func (h *HotLink) SetURL(url string) {
	h.url = url
}

func (h *HotLink) GetURL() string {
	return h.url
}

func (h *HotLink) updateUnderlinePath() {
	if h.text == "" {
		h.underlinePath = nil
		return
	}
	center := h.GetCenter()
	width, height := TextWidth(h.textOptions, h.text), TextHeight(h.textOptions)

	if center == nil {
		h.underlinePath = nil
		return
	}

	x, y := center.V()
	bottom := y + height/2
	h.underlinePath = GetVectorPath([]*Point{
		{X: x - width/2, Y: bottom},
		{X: x + width/2, Y: bottom},
	})
}

func (h *HotLink) SetText(text string) {
	h.Element.SetText(text)
	h.updateUnderlinePath()
}

func (h *HotLink) SetCenter(center Point) {
	h.Element.SetCenter(center)
	h.updateUnderlinePath()
}

func (h *HotLink) SetSize(size Point) {
	h.Element.SetSize(size)
	h.updateUnderlinePath()
}

func (h *HotLink) Update() {
	h.Element.Update()

	if h.url == "" {
		return
	}
	if h.Check(UICheckHover) {
		h.SetHovered(true)

		if h.Check(UICheckPress) {
			h.SetPressed(true)

		} else if h.Check(UICheckRelease) && h.IsPressed() {
			h.Trigger()
			h.SetPressed(false)
		}
	} else {
		h.SetHovered(false)
	}
}

func (h *HotLink) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	h.Element.Draw(screen, opts)

	if h.url != "" && h.underlinePath != nil {
		color := types.LightBlue
		if h.IsHovered() {
			color = types.Yellow
		} else if h.pressed {
			color = types.Orange
		}
		h.underlinePath.Draw(screen, 2.0, color.C())
	}
}
