package play

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/liqmix/ebiten-holiday-2024/internal/cache"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
)

var (
	guideColor = color.RGBA{
		R: types.Gray.R,
		G: types.Gray.G,
		B: types.Gray.B,
		A: 150,
	}
	trackActiveAlpha   = uint8(15)
	trackInactiveAlpha = uint8(10)
)

func (r *Play) addJudgementPath(track *types.Track) {
	width := judgementWidth
	if track.IsPressed() {
		width *= 2
	}
	path := GetJudgementLinePath(track.Name, track.IsPressed())
	r.vectorCollection.AddPath(path)
}

func CreateJudgementPath(track types.TrackName, pressed bool) *cache.CachedPath {
	return CreateNotePath(track, 1, &NotePathOpts{
		lineWidth:       judgementWidth,
		largeWidthRatio: judgementPressedRatio,
		isLarge:         pressed,
		color:           track.NoteColor(),
		alpha:           GetNoteFadeAlpha(1),
		solo:            true,
	})
}

func (r *Play) addTrackPath(track *types.Track) {
	notePoints := notePoints[track.Name]
	centerX, centerY := playCenterPoint.ToRender32()

	// Create two paths - one for each side of the lane
	trackPath := vector.Path{}

	// Get starting point dimensions
	startX, startY := notePoints[0].ToRender32()
	width := float32(2) // Adjust this for lane width

	// Calculate perpendicular vector for width
	dx := startX - centerX
	dy := startY - centerY
	length := float32(math.Sqrt(float64(dx*dx + dy*dy)))
	normalX := -dy / length * width
	normalY := dx / length * width

	// Start with the right edge
	trackPath.MoveTo(centerX+(startX-centerX)+normalX, centerY+(startY-centerY)+normalY)

	// Add right edge line segments
	for i := 1; i < len(notePoints); i++ {
		if notePoints[i] == nil {
			continue
		}
		x, y := notePoints[i].ToRender32()
		dx = x - centerX
		dy = y - centerY
		length = float32(math.Sqrt(float64(dx*dx + dy*dy)))
		normalX = -dy / length * width
		normalY = dx / length * width

		trackPath.LineTo(centerX+(x-centerX)+normalX, centerY+(y-centerY)+normalY)
	}

	// Move to center point
	trackPath.LineTo(centerX, centerY)

	// Now add the left edge in reverse order
	for i := len(notePoints) - 1; i >= 0; i-- {
		if notePoints[i] == nil {
			continue
		}
		x, y := notePoints[i].ToRender32()
		dx = x - centerX
		dy = y - centerY
		length = float32(math.Sqrt(float64(dx*dx + dy*dy)))
		normalX = -dy / length * width
		normalY = dx / length * width

		trackPath.LineTo(centerX+(x-centerX)-normalX, centerY+(y-centerY)-normalY)
	}

	// Close the path
	trackPath.Close()

	// Generate vertices and indices for filling
	vs, is := trackPath.AppendVerticesAndIndicesForFilling(nil, nil)

	// Set color with transparency
	color := track.Name.NoteColor()
	if track.IsPressed() {
		color.A = trackActiveAlpha
	} else {
		color.A = trackInactiveAlpha
	}
	ui.ColorVertices(vs, color)

	// Add the path to the vector collection
	r.vectorCollection.Add(vs, is)
}
