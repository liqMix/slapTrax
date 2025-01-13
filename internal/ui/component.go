package ui

import (
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/input"
)

type Positionable interface {
	GetCenter() *Point
	SetCenter(Point)

	GetSize() *Point
	SetSize(Point)
}

type Textable interface {
	GetText() string
	SetText(string)
}

type Pressable interface {
	IsPressed() bool
	SetPressed(bool)
}

type Focusable interface {
	IsFocused() bool
	SetFocused(bool)
}

type Componentable interface {
	Positionable
	Textable
	Pressable
	Focusable

	Draw(*ebiten.Image, *ebiten.DrawImageOptions)
	Update()

	IsHovered() bool
	SetHovered(bool)

	IsDisabled() bool
	SetDisabled(bool)

	IsHidden() bool
	SetHidden(bool)

	GetOpacity() float64
	SetOpacity(float64)

	IsPaneled() bool
	SetPaneled(bool)

	SetTrigger(func())
	Trigger()
}

type UICheckKind int

const (
	UICheckNone UICheckKind = iota
	UICheckHover
	UICheckPress
	UICheckRelease
)

type Clickable interface {
	Componentable
	Check(UICheckKind) bool
}

type Component struct {
	id     string
	center *Point
	size   *Point

	disabled bool
	hidden   bool
	hovered  bool
	pressed  bool
	focused  bool

	paneled bool
	opacity float64
	trigger func()
}

func (c *Component) GetId() string {
	if c.id == "" {
		c.id = uuid.New().String()
	}
	return c.id
}
func (c *Component) SetCenter(center Point) { c.center = &center }
func (c *Component) GetCenter() *Point      { return c.center }

func (c *Component) SetSize(size Point) { c.size = &size }
func (c *Component) GetSize() *Point    { return c.size }

func (c *Component) SetTrigger(trigger func()) { c.trigger = trigger }
func (c *Component) Trigger() {
	if c.trigger != nil {
		c.trigger()
	}
}

func (c *Component) SetDisabled(d bool) { c.disabled = d }
func (c *Component) IsDisabled() bool   { return c.disabled }

func (c *Component) SetHidden(h bool) {
	c.hidden = h
	c.disabled = h
}
func (c *Component) IsHidden() bool { return c.hidden }

func (c *Component) SetHovered(h bool) { c.hovered = h }
func (c *Component) IsHovered() bool   { return c.hovered }

func (c *Component) SetPressed(p bool) { c.pressed = p }
func (c *Component) IsPressed() bool   { return c.pressed }

func (c *Component) SetFocused(f bool) { c.focused = f }
func (c *Component) IsFocused() bool   { return c.focused }

func (c *Component) SetPaneled(p bool) { c.paneled = p }
func (c *Component) IsPaneled() bool   { return c.paneled }

func (c *Component) SetOpacity(o float64) { c.opacity = o }
func (c *Component) GetOpacity() float64  { return c.opacity }

func (c *Component) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	if c.hidden || c.size == nil || c.center == nil {
		return
	}

	if c.paneled {
		// Panel
		size := c.GetSize()
		center := c.GetCenter()
		if c.focused {
			DrawInvertedNoteThemedRect(screen, center, size)
		} else {
			DrawNoteThemedRect(screen, center, size)
		}
	}
}

func (c *Component) Update() {}

func (c *Component) Check(kind UICheckKind) bool {
	if c.disabled || kind == UICheckNone || c.center == nil || c.size == nil {
		return false
	}
	cX, cY := c.center.ToRender()
	w, h := c.size.ToRender()
	x, y := cX-w/2, cY-h/2
	mouseInBounds := input.M.InBounds(x, y, w, h)

	switch kind {
	case UICheckHover:
		return mouseInBounds || c.IsFocused()
	case UICheckPress:
		if !c.hovered {
			return false
		}
		return mouseInBounds && input.M.Is(ebiten.MouseButtonLeft, input.JustPressed)
	case UICheckRelease:
		if !c.pressed {
			return false
		}
		return mouseInBounds && input.M.Is(ebiten.MouseButtonLeft, input.JustReleased)
	}

	return false
}
