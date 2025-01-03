package types

import (
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/cache"
)

func NewRenderImage() *ebiten.Image {
	x, y := Window.RenderSize()
	return ebiten.NewImage(x, y)
}

type Renderer interface {
	Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions)

	// children can split static rendering into
	// the return image which is cached,
	// still drawing dynamic items to the screen
	static(screen *ebiten.Image)
}

type BaseRenderer struct {
	id           string
	StaticRender func(*ebiten.Image)
}

func (r *BaseRenderer) Init(render func(*ebiten.Image)) {
	r.id = uuid.New().String()
	r.StaticRender = render
}

func (r *BaseRenderer) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	img, ok := cache.GetImage(r.id)
	if !ok {
		img = NewRenderImage()
		r.static(img)
		if img != nil {
			cache.SetImage(r.id, img)
		}
	}
	if img != nil {
		screen.DrawImage(img, opts)
	}
}

func (r *BaseRenderer) static(screen *ebiten.Image) {
	r.StaticRender(screen)
}
