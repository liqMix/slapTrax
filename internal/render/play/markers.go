package play

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
)

func (r *Play) drawMeasureMarkers(screen *ebiten.Image) {
	// Get time info from the song/state
	bpm := r.state.Song.BPM
	currentTime := r.state.CurrentTime()

	// Calculate which measures should be visible
	// This depends on how far ahead notes are visible in your game
	msPerMeasure := 60000 * 4 / bpm // assuming 4/4 time
	currentMeasure := currentTime / int64(msPerMeasure)

	markerPoints := GetMeasureMarkerPoints()
	color := types.Gray
	for i := 0; i < 2; i++ { // Show next 2 markers
		measureTime := (currentMeasure + int64(i)) * int64(msPerMeasure)
		progress := 1.0 - float64(measureTime-currentTime)/r.state.TravelTime

		if progress >= 0 && progress <= 1 {
			// Draw marker across all main tracks using a semi-transparent color
			drawMarker(screen, progress, markerPoints, color)
		}
	}
}

func drawMarker(screen *ebiten.Image, p float64, vec []*ui.Point, color color.RGBA) {
	// Calculate the constant screen-space velocity progress
	progress := (minT / (minT + (1-minT)*(1-float32(p))))

	// Early exit if note is not visible yet
	if progress < fadeInThreshold {
		return
	}

	// Calculate alpha based on position in fade range
	switch {
	case progress < fadeInThreshold:
		color.A = 0
	case progress < fadeOutThreshold:
		// Smooth interpolation between thresholds
		fadeProgress := (progress - fadeInThreshold) / (fadeOutThreshold - fadeInThreshold)
		color.A = uint8(float32(maxAlpha) * fadeProgress)
	default:
		color.A = maxAlpha
	}

	// Skip rendering if fully transparent
	if color.A == 0 {
		return
	}

	x1, y1 := vec[0].ToRender32()
	x2, y2 := vec[1].ToRender32()
	x3, y3 := vec[2].ToRender32()
	x4, y4 := vec[3].ToRender32()
	cX, cY := playCenterPoint.ToRender32()

	// leftCurrentX := cX + (leftX-cX)*progress
	// leftCurrentY := cY + (leftY-cY)*progress
	// centerCurrentX := cX + (centerX-cX)*progress
	// centerCurrentY := cY + (centerY-cY)*progress
	// rightCurrentX := cX + (rightX-cX)*progress
	// rightCurrentY := cY + (rightY-cY)*progress

	x1Current := cX + (x1-cX)*progress
	y1Current := cY + (y1-cY)*progress
	x2Current := cX + (x2-cX)*progress
	y2Current := cY + (y2-cY)*progress
	x3Current := cX + (x3-cX)*progress
	y3Current := cY + (y3-cY)*progress
	x4Current := cX + (x4-cX)*progress
	y4Current := cY + (y4-cY)*progress

	markerPath := vector.Path{}
	markerPath.MoveTo(x1Current, y1Current)
	markerPath.LineTo(x2Current, y2Current)
	markerPath.LineTo(x3Current, y3Current)
	markerPath.LineTo(x4Current, y4Current)
	markerPath.Close()

	var width float32 = 8 * progress

	vs, is := markerPath.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{
		Width: width,
	})

	ui.ColorVertices(vs, color)
	screen.DrawTriangles(vs, is, ui.BaseTriImg, nil)

}
