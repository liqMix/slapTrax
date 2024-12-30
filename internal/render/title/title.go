package title

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/state"
	title "github.com/liqmix/ebiten-holiday-2024/internal/state/title"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
)

type TitleRenderer interface {
	New(*title.State) TitleRenderer
	Draw(screen *ebiten.Image)
}

func GetRenderer(s state.State, t types.Theme) TitleRenderer {
	state := s.(*title.State)

	switch t {
	case types.ThemeStandard:
		return Default{}.New(state)
	default:
		return Default{}.New(state)
		// case types.ThemeLeftBehind:
		// 	return LeftBehind{}.New(state)
	}
}
