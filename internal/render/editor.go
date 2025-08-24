package render

import (
	"github.com/liqmix/slaptrax/internal/display"
	"github.com/liqmix/slaptrax/internal/render/editor"
	"github.com/liqmix/slaptrax/internal/state"
)

func NewEditorRenderer(s *state.EditorState) display.Renderer {
	return editor.NewEditorRenderer(s)
}