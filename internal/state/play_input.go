package state

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/l"
)

type PlayAction string

const (
	PauseAction   PlayAction = l.STATE_PLAY_PAUSE
	RestartAction PlayAction = l.STATE_PLAY_RESTART
)

var PlayActions = map[PlayAction][]ebiten.Key{
	RestartAction: {
		ebiten.KeyF5,
	},
	PauseAction: {
		ebiten.KeyEscape,
	},
}
