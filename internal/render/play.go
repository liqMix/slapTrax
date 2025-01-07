package render

import (
	"github.com/liqmix/ebiten-holiday-2024/internal/display"
	"github.com/liqmix/ebiten-holiday-2024/internal/render/play"
	"github.com/liqmix/ebiten-holiday-2024/internal/state"
)

func NewPlayRender(s state.State) display.Renderer {
	return play.NewPlayRender(s)
}
