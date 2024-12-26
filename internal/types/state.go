package types

import "github.com/liqmix/ebiten-holiday-2024/internal/l"

type GameState string

const (
	GameStateTitle  GameState = l.STATE_TITLE
	GameStateMenu   GameState = l.STATE_MENU
	GameStatePlay   GameState = l.STATE_PLAY
	GameStateEditor GameState = l.STATE_EDITOR
)

func (gs GameState) String() string {
	return l.String(string(gs))
}
