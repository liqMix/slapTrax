package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Interactable interface {
	Draw(*ebiten.Image, *ebiten.DrawImageOptions)
	Update()

	GetCenter() *Point
	SetCenter(Point)

	GetSize() *Point
	SetSize(Point)

	GetText() string
	SetText(string)

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

type Component struct {
	center *Point
	size   *Point

	disabled bool
	hidden   bool
	hovered  bool
	paneled  bool
	opacity  float64
	trigger  func()
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
		DrawNoteThemedRect(screen, center, size)
	}
}

func (c *Component) Update() {}
