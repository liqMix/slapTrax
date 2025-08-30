package shaders

import (
	"embed"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed *.kage
var shaderFS embed.FS

type ShaderManager struct {
	noteShader     *ebiten.Shader
	holdNoteShader *ebiten.Shader
}

var Manager *ShaderManager

func InitManager() error {
	Manager = &ShaderManager{}
	
	if err := Manager.loadShaders(); err != nil {
		return err
	}
	
	return nil
}

func (sm *ShaderManager) loadShaders() error {
	var err error
	
	// Load note shader
	noteSource, err := shaderFS.ReadFile("note.kage")
	if err != nil {
		return err
	}
	
	sm.noteShader, err = ebiten.NewShader(noteSource)
	if err != nil {
		return err
	}
	
	// Load hold note shader
	holdSource, err := shaderFS.ReadFile("holdnote.kage")
	if err != nil {
		return err
	}
	
	sm.holdNoteShader, err = ebiten.NewShader(holdSource)
	if err != nil {
		return err
	}
	
	return nil
}

func (sm *ShaderManager) GetNoteShader() *ebiten.Shader {
	return sm.noteShader
}

func (sm *ShaderManager) GetHoldNoteShader() *ebiten.Shader {
	return sm.holdNoteShader
}