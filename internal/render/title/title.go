package title

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/state"
	title "github.com/liqmix/ebiten-holiday-2024/internal/state/title"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
)

type TitleRenderer interface {
	Draw(screen *ebiten.Image)
}

type Renderer struct {
	state    title.State
	renderer TitleRenderer
}

func (r *Renderer) Init(s state.State, t types.Theme) {
	// Get current theme from settings
	switch t {
	case types.ThemeDefault:
		r.renderer = &Default{}
	case types.ThemeLeftBehind:
		r.renderer = &LeftBehind{}
	}

	r.state = *s.(*title.State)
}

func (r *Renderer) Draw(screen *ebiten.Image) {
	r.renderer.Draw(screen)
}
