package play

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/render/shaders"
	"github.com/liqmix/slaptrax/internal/types"
)

// addNotePathShader adds note paths using shader rendering
func (r *Play) addNotePathShader(track *types.Track, screen *ebiten.Image) {
	if len(track.ActiveNotes) == 0 {
		return
	}
	
	// Get track points and center point for this track
	trackPoints := notePoints[track.Name]
	if len(trackPoints) == 0 {
		return
	}
	
	// Render each note using shaders
	for _, note := range track.ActiveNotes {
		if note.IsHoldNote() {
			// Render hold note using hold shader
			shaders.Renderer.RenderHoldNote(screen, track.Name, note, trackPoints, &playCenterPoint)
		} else {
			// Render regular note using note shader
			shaders.Renderer.RenderNote(screen, track.Name, note, trackPoints, &playCenterPoint)
		}
	}
}

// ShaderRenderingEnabled flag to toggle between shader and vertex rendering
var ShaderRenderingEnabled = true

// EnableShaderRendering enables shader-based note rendering
func EnableShaderRendering() {
	ShaderRenderingEnabled = true
}

// DisableShaderRendering disables shader-based note rendering
func DisableShaderRendering() {
	ShaderRenderingEnabled = false
}