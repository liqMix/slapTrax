package play

import (
	"sort"
	
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
	
	// Sort notes by progress for proper depth ordering (back to front)
	// Notes with lower progress (further from judgment line) should render first
	sortedNotes := make([]*types.Note, len(track.ActiveNotes))
	copy(sortedNotes, track.ActiveNotes)
	sort.Slice(sortedNotes, func(i, j int) bool {
		return sortedNotes[i].Progress < sortedNotes[j].Progress
	})
	
	// Render each note using shaders in depth order
	for _, note := range sortedNotes {
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