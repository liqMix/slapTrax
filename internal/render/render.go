package render

import (
	"github.com/liqmix/slaptrax/internal/display"
	"github.com/liqmix/slaptrax/internal/state"
	"github.com/liqmix/slaptrax/internal/types"
)

func GetRenderer(gs types.GameState, s state.State) display.Renderer {
	switch gs {
	case types.GameStatePlay:
		return NewPlayRender(s)
	case types.GameStateOffset:
		return NewOffsetRender(s)
	case types.GameStateEditor:
		return NewEditorRenderer(s.(*state.EditorState))
	}
	return nil
}
