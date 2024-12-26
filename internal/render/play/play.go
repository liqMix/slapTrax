package play

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/song"
	"github.com/liqmix/ebiten-holiday-2024/internal/state"
	play "github.com/liqmix/ebiten-holiday-2024/internal/state/play"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
)

type PlayRenderer interface {
	Draw(screen *ebiten.Image)
	drawTrack(screen *ebiten.Image, t *song.Track)
	drawNote(screen *ebiten.Image, n *song.Note)
}

type Renderer struct {
	state    play.State
	renderer PlayRenderer
}

func (r *Renderer) Init(s state.State, t types.Theme) {
	// Get current theme from settings
	switch t {
	case types.ThemeDefault:
		r.renderer = &Default{}
	case types.ThemeLeftBehind:
		r.renderer = &LeftBehind{}
	}

	r.state = *s.(*play.State)
}

func (r *Renderer) Draw(screen *ebiten.Image) {
	r.renderer.Draw(screen)
}
