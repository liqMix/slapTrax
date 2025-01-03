package render

import (
	"github.com/liqmix/ebiten-holiday-2024/internal/render/play"
	"github.com/liqmix/ebiten-holiday-2024/internal/state"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
)

func NewPlayRender(s state.State) types.Renderer {
	return play.NewPlayRender(s)
}
