package render

import (
	"github.com/liqmix/ebiten-holiday-2024/internal/state"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
)

func GetRenderer(gs types.GameState, s state.State) types.Renderer {
	switch gs {
	case types.GameStatePlay:
		return NewPlayRender(s)
	case types.GameStateOffset:
		return NewOffsetRender(s)
	}
	return nil
}
