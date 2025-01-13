package internal

import (
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/cache"
	"github.com/liqmix/ebiten-holiday-2024/internal/display"
	"github.com/liqmix/ebiten-holiday-2024/internal/render"
	"github.com/liqmix/ebiten-holiday-2024/internal/state"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
)

type RenderState struct {
	id     string
	frozen bool

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
		id:       uuid.New().String(),
		state:    state,
		renderer: render.GetRenderer(gs, state),
	}
}

func (r *RenderState) Freeze() {
	r.frozen = true
}

func (r *RenderState) Unfreeze() {
	r.frozen = false
}
