package state

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/types"
)

type Exit struct {
	types.BaseGameState
}

func NewExitState() *Exit {
	return &Exit{}
}

func (e *Exit) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {}

func (e *Exit) Update() error {
	return ebiten.Termination
}
