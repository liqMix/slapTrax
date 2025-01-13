package state

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/input"
	"github.com/liqmix/slaptrax/internal/types"
)

type FloatStateArgs struct {
	Cb func()
}

type State interface {
	// If state should change, return new state and arg for that state
	Update() error
	Draw(*ebiten.Image, *ebiten.DrawImageOptions)

	// If state should render above the current state
	Floats() bool
	IsDoBeFloating()

	// Set the state as inactive
	SetActive(bool)

	// Are we transitioning to a new state?
	SetNextState(types.GameState, interface{})
	HasNextState() bool
	GetNextState() (types.GameState, interface{})

	CheckActions() input.Action
	IsNavigable() bool
}

var FloatingStates = map[types.GameState]bool{
	types.GameStatePause:               true,
	types.GameStateOffset:              true,
	types.GameStateSettings:            true,
	types.GameStateDifficultySelection: true,
	types.GameStateModal:               true,
	types.GameStateLogin:               true,
	types.GameStateHowToPlay:           true,
	types.GameStateKeyConfig:           true,
}

func New(s types.GameState, arg interface{}) State {
	var state State

	switch s {
	case types.GameStateExit:
		state = NewExitState()
	case types.GameStatePlay:
		state = NewPlayState(arg.(*PlayArgs))
	case types.GameStateTitle:
		state = NewTitleState()
	case types.GameStatePause:
		state = NewPauseState(arg.(*PauseArgs))
	case types.GameStateOffset:
		state = NewOffsetState(arg.(*FloatStateArgs))
	case types.GameStateSettings:
		state = NewSettingsState()
	case types.GameStateSongSelection:
		state = NewSongSelectionState()
	case types.GameStateResult:
		state = NewResultState(arg.(*ResultStateArgs))
	case types.GameStateLogin:
		state = NewLoginState()
	case types.GameStateModal:
		state = NewModalState(arg.(*ModalStateArgs))
	case types.GameStateHowToPlay:
		state = NewHowToPlayState()
	case types.GameStateKeyConfig:
		state = NewKeyConfigState(arg.(*FloatStateArgs))
	}

	if state == nil {
		panic("Invalid state")
	}
	if _, ok := FloatingStates[s]; ok {
		state.IsDoBeFloating()
	}

	return state
}
