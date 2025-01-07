package types

import "github.com/liqmix/ebiten-holiday-2024/internal/l"

type GameState string

const (
	GameStateNone                GameState = ""
	GameStateTitle               GameState = l.STATE_TITLE
	GameStatePlay                GameState = l.STATE_PLAY
	GameStateEditor              GameState = l.STATE_EDITOR
	GameStateOffset              GameState = l.STATE_OFFSET
	GameStatePause               GameState = l.STATE_PLAY_PAUSE
	GameStateSettings            GameState = l.STATE_SETTINGS
	GameStateSongSelection       GameState = l.STATE_SONG_SELECTION
	GameStateDifficultySelection GameState = l.STATE_DIFFICULTY_SELECTION
	GameStateResult              GameState = l.STATE_RESULT
	GameStateLogin               GameState = l.STATE_LOGIN
	GameStateModal               GameState = "modal"
	GameStateBack                GameState = l.BACK
)

func (gs GameState) String() string {
	return string(gs)
}

type BaseGameState struct {
	NextState     GameState
	NextStateArgs interface{}
	Active        bool
	floats        bool
}

func (s *BaseGameState) SetActive(active bool) {
	s.Active = active
}

func (s *BaseGameState) Floats() bool {
	return s.floats
}

func (s *BaseGameState) IsDoBeFloating() {
	s.floats = true
}

func (s *BaseGameState) SetNextState(nextState GameState, args interface{}) {
	s.NextState = nextState
	s.NextStateArgs = args
}

func (s *BaseGameState) HasNextState() bool {
	return s.NextState != GameStateNone
}

func (s *BaseGameState) GetNextState() (GameState, interface{}) {
	return s.NextState, s.NextStateArgs
}
