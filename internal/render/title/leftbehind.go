package title

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/state/title"
)

type LeftBehind struct {
	state *title.State
}

func (r LeftBehind) New(s *title.State) TitleRenderer {
	return &LeftBehind{
		state: s,
	}
}
func (r *LeftBehind) Draw(screen *ebiten.Image) {
}
