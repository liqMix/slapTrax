package title

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/l"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
)

type State struct {
	Text string
	Flag *ebiten.Image
}

func New(arg interface{}) *State {
	return &State{
		Text: l.String(l.TITLE),
		Flag: l.Flag(),
	}
}

func (s *State) Update() (*types.GameState, interface{}, error) {
	return nil, nil, nil
}
