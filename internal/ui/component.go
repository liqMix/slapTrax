package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
)

type Interactable interface {
	Draw(*ebiten.Image, *ebiten.DrawImageOptions)
	Update()

	GetCenter() *Point
	SetCenter(Point)

	GetSize() *Point
	SetSize(Point)

	IsHovered() bool
	SetHovered(bool)

	IsDisabled() bool
	SetDisabled(bool)

	IsPaneled() bool
	SetPaneled(bool)

	SetTrigger(func())
	Trigger()
}

type Component struct {
	center *Point
	size   *Point

	disabled bool
	hovered  bool
	paneled  bool
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

func (c *Component) SetHovered(h bool) { c.hovered = h }
func (c *Component) IsHovered() bool   { return c.hovered }

func (c *Component) SetPaneled(p bool) { c.paneled = p }
func (c *Component) IsPaneled() bool   { return c.paneled }

func (c *Component) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	if c.paneled {
		// Panel
		border := 0.02
		size := c.GetSize()
		center := c.GetCenter()
		borderSize := &Point{
			X: size.X + border*2,
			Y: size.Y + border*2,
		}
		DrawFilledRect(screen, center, borderSize, types.Gray)
		DrawFilledRect(screen, center, size, types.Black)
	}
}
func (c *Component) Update() {}
