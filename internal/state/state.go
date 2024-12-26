package state

import (
	"github.com/liqmix/ebiten-holiday-2024/internal/state/play"
	"github.com/liqmix/ebiten-holiday-2024/internal/state/title"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
)

var (
	Play  play.State
	Title title.State
)

type State interface {
	// If state should change, return new state and arg for that state
	Update() (*types.GameState, interface{}, error)
}

func New(s types.GameState, arg interface{}) State {
	switch s {
	case types.GameStatePlay:
		return play.New(arg)
	case types.GameStateTitle:
		return title.New(arg)
	}
	return nil
}
