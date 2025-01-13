package types

import (
	"github.com/liqmix/ebiten-holiday-2024/internal/input"
	"github.com/liqmix/ebiten-holiday-2024/internal/l"
)

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
	GameStateExit                GameState = l.EXIT
	GameStateHowToPlay           GameState = "howtoplay"
	GameStateKeyConfig           GameState = l.SETTINGS_GAME_KEY_CONFIG
)

func (gs GameState) String() string {
	return string(gs)
}

type BaseGameState struct {
	NextState        GameState
	NextStateArgs    interface{}
	Active           bool
	AvailableActions []input.Action
	actions          map[input.Action]func()
	floats           bool
	notNavigable     bool
}

func (s *BaseGameState) SetActive(active bool) {
	s.Active = active
}

func (s *BaseGameState) SetAction(action input.Action, f func()) {
	if s.actions == nil {
		s.actions = make(map[input.Action]func())
		s.AvailableActions = append(s.AvailableActions, action)
	}
	s.actions[action] = f
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

func (s *BaseGameState) CheckActions() input.Action {
	action := input.ActionUnknown
	for _, a := range s.AvailableActions {
		if input.JustActioned(a) {
			if f, ok := s.actions[a]; ok {
				f()
				if action == input.ActionUnknown {
					action = a
				}
			}
		}
	}
	return action
}

func (s *BaseGameState) SetNotNavigable() {
	s.notNavigable = true
}

func (s *BaseGameState) IsNavigable() bool {
	return !s.notNavigable
}

func (s *BaseGameState) Update() error {
	s.CheckActions()
	return nil
}
