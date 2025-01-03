package types

type GameState string

const (
	GameStateNone                GameState = ""
	GameStateTitle               GameState = L_STATE_TITLE
	GameStatePlay                GameState = L_STATE_PLAY
	GameStateEditor              GameState = L_STATE_EDITOR
	GameStateOffset              GameState = L_STATE_OFFSET
	GameStatePause               GameState = L_STATE_PLAY_PAUSE
	GameStateSettings            GameState = L_STATE_SETTINGS
	GameStateSongSelection       GameState = L_STATE_SONG_SELECTION
	GameStateDifficultySelection GameState = L_STATE_DIFFICULTY_SELECTION
	GameStateBack                GameState = L_BACK
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
