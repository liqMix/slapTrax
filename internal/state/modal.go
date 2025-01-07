package state

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
)

type Modal struct {
	types.BaseGameState

	update func(func(types.GameState, interface{}))
	draw   func(i *ebiten.Image, o *ebiten.DrawImageOptions)
}

type ModalStateArgs struct {
	Update func(func(types.GameState, interface{}))
	Draw   func(i *ebiten.Image, o *ebiten.DrawImageOptions)
}

func NewModalState(args *ModalStateArgs) *Modal {
	return &Modal{
		update: args.Update,
		draw:   args.Draw,
	}
}

func (m *Modal) Update() error {
	m.update(m.SetNextState)
	return nil
}

func (m *Modal) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	m.draw(screen, opts)
}
