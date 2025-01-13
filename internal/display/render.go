package display

import (
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/cache"
	"github.com/liqmix/slaptrax/internal/l"
)

func NewRenderImage() *ebiten.Image {
	x, y := Window.RenderSize()
	return ebiten.NewImage(x, y)
}

// Allowed render sizes
type RenderSize string

const (
	RenderSizeTiny   RenderSize = l.RENDERSIZE_TINY
	RenderSizeSmall  RenderSize = l.RENDERSIZE_SMALL
	RenderSizeMedium RenderSize = l.RENDERSIZE_MEDIUM
	RenderSizeLarge  RenderSize = l.RENDERSIZE_LARGE
	RenderSizeMax    RenderSize = l.RENDERSIZE_MAX
)

func (rs RenderSize) String() string {
	return string(rs)
}
func (rs RenderSize) Value() (int, int) {
	switch rs {
	case RenderSizeTiny:
		return 640, 360
	case RenderSizeSmall:
		return 960, 540
	case RenderSizeMedium:
		return 1280, 720
	case RenderSizeLarge:
		return 1920, 1080
	case RenderSizeMax:
		w, _ := ebiten.Monitor().Size()
		return w, w * 9 / 16
	}
	return 1280, 720
}

type Renderer interface {
	Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions)

	// draw static to separate image instead of screen
	static(image *ebiten.Image, opts *ebiten.DrawImageOptions)
}

type BaseRenderer struct {
	id           string
	StaticRender func(*ebiten.Image, *ebiten.DrawImageOptions)
}

func (r *BaseRenderer) Init(render func(*ebiten.Image, *ebiten.DrawImageOptions)) {
	r.id = uuid.New().String()
	r.StaticRender = render
}

func (r *BaseRenderer) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	img, ok := cache.Image.Get(r.id)
	if !ok {
		img = NewRenderImage()
		r.static(img, opts)
		if img != nil {
			cache.Image.Set(r.id, img)
		}
	}
	if img != nil {
		screen.DrawImage(img, opts)
	}
}

func (r *BaseRenderer) static(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	r.StaticRender(screen, opts)
}
