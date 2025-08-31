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

// Freezable interface for renderers that need to know when they're being frozen
type Freezable interface {
	SetFrozen(bool)
}

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
	// Notify renderer if it implements Freezable interface
	if freezable, ok := r.renderer.(Freezable); ok {
		freezable.SetFrozen(true)
	}
}

func (r *RenderState) Unfreeze() {
	r.frozen = false
	// Notify renderer if it implements Freezable interface
	if freezable, ok := r.renderer.(Freezable); ok {
		freezable.SetFrozen(false)
	}
}
