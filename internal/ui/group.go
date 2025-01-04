package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/assets"
	"github.com/liqmix/ebiten-holiday-2024/internal/input"
)

type UIGroup struct {
	Component

	items      []Interactable
	current    Interactable
	currentIdx int
	horizontal bool
}

func NewUIGroup() *UIGroup {
	return &UIGroup{
		Component: Component{},
		items:     make([]Interactable, 0),
	}
}

func (g *UIGroup) Add(i Interactable) {
	g.items = append(g.items, i)
	if g.current == nil {
		g.Hover(0)
	}
	g.updateLayout()
}

func (g *UIGroup) updateLayout() {
	// Update size to fit all items
	size := Point{}
	minPoint := Point{}
	maxPoint := Point{}
	for _, item := range g.items {
		itemSize := item.GetSize()
		itemCenter := item.GetCenter()
		if itemSize == nil || itemCenter == nil {
			continue
		}

		top := itemCenter.Y - itemSize.Y/2
		bottom := itemCenter.Y + itemSize.Y/2
		left := itemCenter.X - itemSize.X/2
		right := itemCenter.X + itemSize.X/2

		if top < minPoint.Y {
			minPoint.Y = top
		}
		if bottom > maxPoint.Y {
			maxPoint.Y = bottom
		}
		if left < minPoint.X {
			minPoint.X = left
		}
		if right > maxPoint.X {
			maxPoint.X = right
		}
	}

	size.X = maxPoint.X - minPoint.X
	size.Y = maxPoint.Y - minPoint.Y
	g.SetSize(size)

	g.SetCenter(Point{
		X: minPoint.X + size.X/2,
		Y: minPoint.Y + size.Y/2,
	})
}

func (g *UIGroup) SetDisabled(d bool) {
	g.Component.SetDisabled(d)
	for _, item := range g.items {
		item.SetDisabled(d)
	}
}

func (g *UIGroup) SetHorizontal() {
	g.horizontal = true
}

func (g *UIGroup) Hover(idx int) {
	if g.IsDisabled() || idx < 0 || idx >= len(g.items) {
		return
	}

	if g.current != nil {
		g.current.SetHovered(false)
	}
	g.currentIdx = idx
	g.current = g.items[idx]
	g.current.SetHovered(true)
}

func (g *UIGroup) Update() {
	if g.IsDisabled() {
		return
	}
	downKey := ebiten.KeyArrowDown
	upKey := ebiten.KeyArrowUp
	if g.horizontal {
		downKey = ebiten.KeyArrowLeft
		upKey = ebiten.KeyArrowRight
	}

	// Move through buttons with arrow keys
	if input.K.Is(downKey, input.JustPressed) {
		if g.currentIdx < len(g.items)-1 {
			assets.PlaySFX(assets.SFXSelectDown)
			g.Hover(g.currentIdx + 1)
		}
	} else if input.K.Is(upKey, input.JustPressed) {
		if g.currentIdx > 0 {
			assets.PlaySFX(assets.SFXSelectUp)
			g.Hover(g.currentIdx - 1)
		}
	} else if input.K.Is(ebiten.KeyEnter, input.JustPressed) {
		g.current.Trigger()
	}

	for i, item := range g.items {
		item.Update()
		if item.IsHovered() && i != g.currentIdx {
			g.Hover(i)
		}
	}
}

func (g *UIGroup) Get() Interactable {
	if g.current == nil && len(g.items) > 0 {
		g.Hover(0)
		return g.current
	}
	return g.current
}

func (g *UIGroup) GetIndex() int {
	return g.currentIdx
}

func (g *UIGroup) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	g.Component.Draw(screen, opts)

	if len(g.items) == 0 {
		return
	}

	for i, item := range g.items {
		if i == g.currentIdx {
			continue
		}
		item.Draw(screen, opts)
	}

	// Draw the current item last so it appears on top
	if g.currentIdx >= 0 {
		g.items[g.currentIdx].Draw(screen, opts)
	}
}
