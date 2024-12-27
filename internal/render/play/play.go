package play

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/state"
	play "github.com/liqmix/ebiten-holiday-2024/internal/state/play"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
)

type PlayRenderer interface {
	New(*play.State) PlayRenderer
	Draw(screen *ebiten.Image)

	drawBackground(screen *ebiten.Image)
	drawProfile(screen *ebiten.Image)
	drawSongInfo(screen *ebiten.Image)
	drawScore(screen *ebiten.Image)
	drawTracks(screen *ebiten.Image)
}

func GetRenderer(s state.State, t types.Theme) PlayRenderer {
	state := s.(*play.State)

	switch t {
	case types.ThemeDefault:
		return Default{}.New(state)
	default:
		return Default{}.New(state)
		// case types.ThemeLeftBehind:
		// 	r.renderer = &LeftBehind{}
	}
}
