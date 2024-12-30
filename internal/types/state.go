package types

type GameState string

const (
	GameStateTitle  GameState = L_STATE_TITLE
	GameStateMenu   GameState = L_STATE_MENU
	GameStatePlay   GameState = L_STATE_PLAY
	GameStateEditor GameState = L_STATE_EDITOR
)
