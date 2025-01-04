package play

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
)

func (r *Play) drawMeasureMarkers(screen *ebiten.Image) {
	// Add marker every quarter note
	interval := r.state.Song.GetQuarterNoteInterval()
	currentTime := r.state.CurrentTime()
	currentMeasure := currentTime / interval
	color := types.Gray

	var i int64
	for i = 0; i < 8; i++ { // Show next 8 markers
		measureTime := (currentMeasure + i) * interval
		progress := types.GetTrackProgress(measureTime, currentTime, r.state.TravelTime)
		drawMarker(screen, progress, measureMarkerPoints, color)
	}
}

// Draw marker across all main tracks using a semi-transparent color
func drawMarker(screen *ebiten.Image, p float64, vec []*ui.Point, color color.RGBA) {
	// Calculate the constant screen-space velocity progress
	progress := SmoothProgress(p)
	color.A = GetFadeAlpha(progress)

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

	var width float32 = 8 * progress

	vs, is := markerPath.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{
		Width: width,
	})

	ui.ColorVertices(vs, color)
	screen.DrawTriangles(vs, is, ui.BaseTriImg, nil)

}
