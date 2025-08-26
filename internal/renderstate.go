package internal

import (
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/cache"
	"github.com/liqmix/slaptrax/internal/display"
	"github.com/liqmix/slaptrax/internal/render"
	"github.com/liqmix/slaptrax/internal/state"
	"github.com/liqmix/slaptrax/internal/types"
)

type RenderState struct {
	id        string
	frozen    bool
	stateType types.GameState

	state    state.State
	renderer display.Renderer
}

func (r *RenderState) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	if r.frozen {
		img, ok := cache.Image.Get(r.id)
		if ok {
			screen.DrawImage(img, opts)
		}
	}
	draw := r.state.Draw
	if r.renderer != nil {
		draw = r.renderer.Draw
	}

	if !r.frozen {
		draw(screen, opts)
		return
	}

	image := display.NewRenderImage()
	draw(image, opts)
	cache.Image.Set(r.id, image)
	screen.DrawImage(image, opts)
}

func GetState(gs types.GameState, arg interface{}) *RenderState {
	state := state.New(gs, arg)
	return &RenderState{
		id:        uuid.New().String(),
		stateType: gs,
		state:     state,
		renderer:  render.GetRenderer(gs, state),
	}
}

func (r *RenderState) Freeze() {
	r.frozen = true
}

func (r *RenderState) Unfreeze() {
	r.frozen = false
}
