package render

import (
	"github.com/liqmix/slaptrax/internal/display"
	"github.com/liqmix/slaptrax/internal/render/play"
	"github.com/liqmix/slaptrax/internal/state"
)

func NewPlayRender(s state.State) display.Renderer {
	return play.NewPlayRender(s)
}
