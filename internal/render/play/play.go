package play

import (
	"github.com/liqmix/ebiten-holiday-2024/internal/render/play/def"
	"github.com/liqmix/ebiten-holiday-2024/internal/render/play/standard"
	"github.com/liqmix/ebiten-holiday-2024/internal/state"
	play "github.com/liqmix/ebiten-holiday-2024/internal/state/play"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
)

func GetRenderer(s state.State, t types.Theme) def.PlayRenderer {
	var renderer def.PlayRenderer
	state := s.(*play.State)

	switch t {
	case types.ThemeStandard:
		renderer = &standard.Standard{}
	default:
		renderer = &standard.Standard{}
		// case types.ThemeLeftBehind:
		// 	r.renderer = &LeftBehind{}
	}

	renderer.Init(state)
	return renderer
}
