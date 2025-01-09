package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/audio"
	"github.com/liqmix/ebiten-holiday-2024/internal/input"
)

type UIGroup struct {
	Component

	items          []Componentable
	currentIdx     int
	usingKeyboard  bool
	horizontal     bool
	triggerOnHover bool
}

func NewUIGroup() *UIGroup {
	return &UIGroup{
		Component: Component{},
		items:     make([]Componentable, 0),
	}
}

func (g *UIGroup) Add(i Componentable) {
	if g.disabled {
		i.SetDisabled(true)
	}

	g.items = append(g.items, i)
	if g.Get() != nil {
		g.Select(g.currentIdx)
	} else {
		g.currentIdx++
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

func (g *UIGroup) Select(idx int) {
	if idx < 0 || idx >= len(g.items) || len(g.items) == 0 {
		return
	}

	item := g.Get()
	if item != nil {
		item.SetHovered(false)
		item.SetPressed(false)
	}

	g.currentIdx = idx
	item = g.Get()
	if item == nil {
		return
	}

	item.SetHovered(true)
	if g.triggerOnHover {
		item.Trigger()
	}
}

func (g *UIGroup) Get() Componentable {
	if g.currentIdx < 0 || g.currentIdx >= len(g.items) || len(g.items) == 0 {
		return nil
	}
	current := g.items[g.currentIdx]
	if current.IsDisabled() {
		return nil
	}
	return current
}

func (g *UIGroup) GetIndex() int {
	return g.currentIdx
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
	if g.IsDisabled() || len(g.items) == 0 {
		return
	}

	if input.M.DidMove() {
		g.usingKeyboard = false
	}
	for _, item := range g.items {
		if item != nil && item.IsFocused() {
			item.Update()
			return
		}
	}

	// if !g.usingKeyboard {
	// 	for _, item := range g.items {
	// 		item.Update()

	// 		if item.Check(UICheckHover) {
	// 			item.SetHovered(true)

	// 			if item.Check(UICheckPress) {
	// 				item.SetPressed(true)

	// 			} else if item.Check(UICheckRelease) && item.IsPressed() {
	// 				item.Trigger()
	// 				item.SetPressed(false)
	// 				g.Select(g.currentIdx)
	// 			}
	// 		} else {
	// 			item.SetHovered(false)
	// 		}
	// 	}
	// }

	downKey := ebiten.KeyArrowDown
	upKey := ebiten.KeyArrowUp
	if g.horizontal {
		downKey = ebiten.KeyArrowRight
		upKey = ebiten.KeyArrowLeft
	}

	sfxCode := audio.SFXNone

	// Move through buttons with arrow keys
	if input.K.Is(downKey, input.JustPressed) {
		g.usingKeyboard = true
		if g.currentIdx < len(g.items)-1 {
			next := g.getNext()
			if next == g.currentIdx {
				return
			}
			g.Select(next)
			sfxCode = audio.SFXSelectDown
		}
	} else if input.K.Is(upKey, input.JustPressed) {
		g.usingKeyboard = true
		if g.currentIdx > 0 {
			prev := g.getPrev()
			if prev == g.currentIdx {
				return
			}
			g.Select(prev)
			sfxCode = audio.SFXSelectUp
		}
	} else if input.K.Is(ebiten.KeyEnter, input.JustPressed) {
		g.usingKeyboard = true
		g.items[g.currentIdx].Trigger()
	}

	if sfxCode != audio.SFXNone {
		audio.PlaySFX(sfxCode)
	}
}

func (g *UIGroup) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	if len(g.items) == 0 {
		return
	}
	g.Component.Draw(screen, opts)

	for i, item := range g.items {
		if i == g.currentIdx {
			continue
		}
		item.Draw(screen, opts)
	}

	// Draw the current item last so it appears on top
	item := g.Get()
	if item == nil {
		return
	}
	item.Draw(screen, opts)
}
