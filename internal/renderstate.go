package internal

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/display"
	"github.com/liqmix/ebiten-holiday-2024/internal/render"
	"github.com/liqmix/ebiten-holiday-2024/internal/state"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
)

type RenderState struct {
	state       state.State
	renderer    display.Renderer
	frozenImage *ebiten.Image
}

func (r *RenderState) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	if r.frozenImage != nil {
		screen.DrawImage(r.frozenImage, opts)
	} else if r.renderer != nil {
		r.renderer.Draw(screen, opts)
	} else {
		r.state.Draw(screen, opts)
	}
}

func getState(gs types.GameState, arg interface{}) *RenderState {
	state := state.New(gs, arg)
	return &RenderState{
		state:    state,
		renderer: render.GetRenderer(gs, state),
	}
}

func (r *RenderState) Freeze() {
	img := display.NewRenderImage()
	r.Draw(img, nil)
	r.frozenImage = img
}

func (r *RenderState) Unfreeze() {
	r.frozenImage = nil
}
