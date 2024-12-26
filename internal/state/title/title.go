package title

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/l"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
)

type State struct {
	text string
	flag *ebiten.Image
}

func New(arg interface{}) *State {
	return &State{
		text: l.String(l.DEBUG_GREETING),
		flag: l.Flag(),
	}
}

func (s *State) Update() (*types.GameState, interface{}, error) {
	return nil, nil, nil
}
