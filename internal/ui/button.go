package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/input"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
)

type UICheckKind int

const (
	UICheckNone UICheckKind = iota
	UICheckHover
	UICheckPress
	UICheckRelease
)

type Button struct {
	Element

	pressed       bool
	underlinePath *VectorPath
}

func NewButton() *Button {
	return &Button{
		Element: *NewElement(),
	}
}

func (c *Button) Update() {
	if c.Check(UICheckHover) {
		c.SetHovered(true)
		if c.Check(UICheckPress) {
			c.pressed = true
		} else if c.Check(UICheckRelease) && c.pressed {
			c.Element.Trigger()
			c.pressed = false
		}
	} else {
		c.SetHovered(false)
	}
	c.Element.Update()
}

func (c *Button) Check(kind UICheckKind) bool {
	if c.disabled || kind == UICheckNone {
		return false
	}

	cX, cY := c.center.ToRender()
	w, h := c.size.ToRender()
	x, y := cX-w/2, cY-h/2
	mouseInBounds := input.M.InBounds(x, y, w, h)

	switch kind {
	case UICheckHover:
		return mouseInBounds
	case UICheckPress:
		return mouseInBounds && input.M.Is(ebiten.MouseButtonLeft, input.JustPressed)
	case UICheckRelease:
		return mouseInBounds && input.M.Is(ebiten.MouseButtonLeft, input.JustReleased)
	}

	return false
}

func (c *Button) updateUnderlinePath() {
	text := c.GetText()
	if text == "" {
		c.underlinePath = nil
		return
	}
	center := c.GetCenter()
	size := c.GetSize()

	if center == nil || size == nil {
		c.underlinePath = nil
		return
	}

	x, y := center.V()
	w, h := size.V()
	bottom := y + h/2
	c.underlinePath = GetVectorPath([]*Point{
		{X: x - w/2, Y: bottom},
		{X: x + w/2, Y: bottom},
	})
}

func (c *Button) SetText(text string) {
	c.Element.SetText(text)
	c.updateUnderlinePath()
}

func (c *Button) SetCenter(center Point) {
	c.Element.SetCenter(center)
	c.updateUnderlinePath()
}

func (c *Button) SetSize(size Point) {
	c.Element.SetSize(size)
	c.updateUnderlinePath()
}

func (c *Button) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	// don't draw hover on element text...
	hovered := c.IsHovered()
	c.Element.SetHovered(false)
	c.Element.Draw(screen, opts)
	c.Element.SetHovered(hovered)

	text := c.GetText()
	if text != "" {
		if !c.disabled && c.trigger != nil && c.underlinePath != nil {
			color := types.LightBlue
			if c.IsHovered() {
				color = types.Yellow
			} else if c.pressed {
				color = types.Orange
			}
			c.underlinePath.Draw(screen, 1.0, color)
		}
	}
}
