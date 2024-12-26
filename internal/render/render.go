package render

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/render/play"
	"github.com/liqmix/ebiten-holiday-2024/internal/render/title"
	"github.com/liqmix/ebiten-holiday-2024/internal/state"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/liqmix/ebiten-holiday-2024/internal/user"
)

// TODO: Clean up this mess of renderer renrderd rennder
type IRenderer interface {
	Init(s state.State, t types.Theme)
	Draw(screen *ebiten.Image)
}

func GetRenderer(gs types.GameState, s state.State) IRenderer {
	var r IRenderer
	switch gs {
	case types.GameStatePlay:
		r = &play.Renderer{}
	case types.GameStateTitle:
		r = &title.Renderer{}
	}

	// Get current theme from settings
	r.Init(s, user.Current.Settings.Theme)
	return r
}
