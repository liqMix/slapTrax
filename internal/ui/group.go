package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/audio"
	"github.com/liqmix/ebiten-holiday-2024/internal/input"
)

type UIGroup struct {
	Component

	items          []Interactable
	current        Interactable
	currentIdx     int
	horizontal     bool
	triggerOnHover bool
}

func NewUIGroup() *UIGroup {
	return &UIGroup{
		Component: Component{},
		items:     make([]Interactable, 0),
	}
}

func (g *UIGroup) Add(i Interactable) {
	g.items = append(g.items, i)
	if g.current == nil && !i.IsDisabled() {
		g.Hover(len(g.items) - 1)
	}
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

func (g *UIGroup) SetVertical() {
	g.horizontal = false
}

func (g *UIGroup) SetTriggerOnHover(t bool) {
	g.triggerOnHover = t
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
	if g.triggerOnHover {
		g.current.Trigger()
	}
}

// Find next not disabled item until the end of the list
func (g *UIGroup) getNext() int {
	for i := g.currentIdx + 1; i < len(g.items); i++ {
		if !g.items[i].IsDisabled() {
			return i
		}
	}
	return g.currentIdx
}

func (g *UIGroup) getPrev() int {
	for i := g.currentIdx - 1; i >= 0; i-- {
		if !g.items[i].IsDisabled() {
			return i
		}
	}
	return g.currentIdx
}

func (g *UIGroup) Update() {
	if g.IsDisabled() {
		return
	}

	downKey := ebiten.KeyArrowDown
	upKey := ebiten.KeyArrowUp
	if g.horizontal {
		downKey = ebiten.KeyArrowRight
		upKey = ebiten.KeyArrowLeft
	}

	sfxCode := audio.SFXNone

	// Move through buttons with arrow keys
	if input.K.Is(downKey, input.JustPressed) {
		if g.currentIdx < len(g.items)-1 {
			next := g.getNext()
			if next == g.currentIdx {
				return
			}
			g.Hover(next)
			if g.horizontal {
				sfxCode = audio.SFXSelectLeft
			} else {
				sfxCode = audio.SFXSelectDown
			}
		}
	} else if input.K.Is(upKey, input.JustPressed) {
		if g.currentIdx > 0 {
			prev := g.getPrev()
			if prev == g.currentIdx {
				return
			}
			g.Hover(prev)
			if g.horizontal {
				sfxCode = audio.SFXSelectRight
			} else {
				sfxCode = audio.SFXSelectUp
			}
		}
	} else if input.K.Is(ebiten.KeyEnter, input.JustPressed) {
		if g.current != nil {
			g.current.Trigger()
		}
	}

	if sfxCode != audio.SFXNone {
		audio.PlaySFX(sfxCode)
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
