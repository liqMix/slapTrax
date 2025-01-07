package play

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
)

func (r *Play) drawMeasureMarkers(screen *ebiten.Image) {
	currentTime := r.state.CurrentTime()
	beatInterval := r.state.Song.GetBeatInterval()
	measureInterval := beatInterval * 4
	color := types.Gray

	// Draw beat markers
	for i := int64(0); i < 8; i++ {
		beatTime := ((currentTime / beatInterval) + i) * beatInterval
		color.A = beatMarkerAlpha
		progress := types.GetTrackProgress(beatTime, currentTime, r.state.GetTravelTime())
		drawMarker(screen, progress, measureMarkerPoints, color)
	}

	// Draw measure markers
	for i := int64(0); i < 2; i++ {
		measureTime := ((currentTime / measureInterval) + i) * measureInterval
		color.A = beatMarkerAlpha * 2
		progress := types.GetTrackProgress(measureTime, currentTime, r.state.GetTravelTime())
		drawMarker(screen, progress, measureMarkerPoints, color)
	}
}

// Draw marker across all main tracks using a semi-transparent color
func drawMarker(screen *ebiten.Image, p float64, vec []*ui.Point, color types.GameColor) {
	if p < 0 || p > 1 {
		return
	}
	progress := SmoothProgress(p)

	color.A = GetFadeAlpha(progress, color.A)

	// Skip rendering if fully transparent
	if color.A == 0 {
		return
	}

	markerPath := vector.Path{}
	cX, cY := playCenterPoint.ToRender32()

	startX, startY := vec[0].ToRender32()
	x, y := cX+(startX-cX)*progress, cY+(startY-cY)*progress
	markerPath.MoveTo(x, y)

	for i := 1; i < len(vec); i++ {
		if vec[i] == nil {
			continue
		}
		x, y = vec[i].ToRender32()
		x, y = cX+(x-cX)*progress, cY+(y-cY)*progress
		markerPath.LineTo(x, y)
	}
	markerPath.Close()

	var width float32 = markerWidth * progress

	vs, is := markerPath.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{
		Width: width,
	})

	ui.ColorVertices(vs, color.C())
	screen.DrawTriangles(vs, is, ui.BaseTriImg, nil)

}
