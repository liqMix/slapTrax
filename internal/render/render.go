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
	Draw(screen *ebiten.Image)
}

func GetRenderer(gs types.GameState, s state.State) IRenderer {
	t := user.Settings().Theme

	switch gs {
	case types.GameStatePlay:
		return play.GetRenderer(s, t)
	case types.GameStateTitle:
		return title.GetRenderer(s, t)
	}

	panic("No renderer found for game state" + gs.String())
}
