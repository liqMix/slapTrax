package state

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
)

type PlayAction string

const (
	RestartAction PlayAction = types.L_STATE_PLAY_RESTART
	PauseAction   PlayAction = types.L_STATE_PLAY_PAUSE
)

var PlayActions = map[PlayAction][]ebiten.Key{
	RestartAction: {
		ebiten.KeyF5,
	},
	PauseAction: {
		ebiten.KeyEscape,
	},
}
