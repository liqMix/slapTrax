package state

import "github.com/liqmix/ebiten-holiday-2024/internal/types"

type Result struct {
	types.BaseGameState
}

func NewResultState() *Result {
	return &Result{}
}
