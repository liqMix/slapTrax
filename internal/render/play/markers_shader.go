package play

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/render/shaders"
	"github.com/liqmix/slaptrax/internal/types"
)

func (r *Play) drawMeasureMarkersShader(screen *ebiten.Image) {
	if shaders.LaneRendererInstance == nil {
		// Fallback to old rendering if shaders not available
		r.drawMeasureMarkers(screen)
		return
	}
	
	currentTime := r.state.CurrentTime()
	beatInterval := r.state.Song.GetBeatInterval()
	measureInterval := beatInterval * 4
	
	// Draw beat markers
	for i := int64(0); i < 8; i++ {
		beatTime := ((currentTime / beatInterval) + i) * beatInterval
		rawProgress := types.GetTrackProgress(beatTime, currentTime, r.state.GetTravelTime())
		
		if rawProgress < 0 || rawProgress > 1 {
			continue
		}
		
		// Apply smooth progress transformation (same as old system)
		progress := float64(SmoothProgress(rawProgress))
		
		// Render shader-based beat marker
		shaders.LaneRendererInstance.RenderMarker(screen, progress, measureMarkerPoints, &playCenterPoint, false)
	}

	// Draw measure markers
	for i := int64(0); i < 2; i++ {
		measureTime := ((currentTime / measureInterval) + i) * measureInterval
		rawProgress := types.GetTrackProgress(measureTime, currentTime, r.state.GetTravelTime())
		
		if rawProgress < 0 || rawProgress > 1 {
			continue
		}
		
		// Apply smooth progress transformation (same as old system)
		progress := float64(SmoothProgress(rawProgress))
		
		// Render shader-based measure marker
		shaders.LaneRendererInstance.RenderMarker(screen, progress, measureMarkerPoints, &playCenterPoint, true)
	}
}