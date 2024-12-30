package standard

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/liqmix/ebiten-holiday-2024/internal/config"
	"github.com/liqmix/ebiten-holiday-2024/internal/song"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
)

func (r *Standard) drawMeasureMarkers(screen *ebiten.Image) {
	// Get time info from the song/state
	bpm := r.state.Song.BPM
	currentTime := r.state.CurrentTime()

	// Calculate which measures should be visible
	// This depends on how far ahead notes are visible in your game
	msPerMeasure := 60000 * 4 / bpm // assuming 4/4 time
	currentMeasure := currentTime / int64(msPerMeasure)

	// Draw markers for next few visible measures
	for i := 0; i < 4; i++ { // Show next 4 measures
		measureTime := (currentMeasure + int64(i)) * int64(msPerMeasure)
		progress := 1.0 - ((measureTime - currentTime) / config.ActualTravelTimeInt64)

		if progress >= 0 && progress <= 1 {
			// Draw marker across all main tracks using a semi-transparent color
			markerColor := white // Very transparent white
			for _, track := range song.MainTracks {
				if config, ok := laneConfigs[track]; ok {
					r.drawMeasureMarker(screen, float32(progress), config, markerColor)
				}
			}
		}
	}
}

func (r *Standard) drawMeasureMarker(screen *ebiten.Image, p float32, config *LaneConfig, markerColor color.RGBA) {
	left := config.Left.RenderPoint()
	right := config.Right.RenderPoint()
	center := config.Center.RenderPoint()
	vanishingPoint := config.VanishingPoint.RenderPoint()

	// To achieve constant screen-space velocity, we need to use inverse interpolation
	// The formula is: p = 1 / (1 + (1/t - 1) * progress) where t is close to 0 (like 0.01)
	// This gives us perfect constant speed in screen space
	t := float32(0.01) // Small value for the starting point relative to vanishing point
	progress := (t / (t + (1-t)*(1-p)))

	var ratio float32 = 0.0
	// Calculate start points near vanishing point
	leftStart := ui.Point{
		X: vanishingPoint.X + (left.X-vanishingPoint.X)*ratio,
		Y: vanishingPoint.Y + (left.Y-vanishingPoint.Y)*ratio,
	}
	centerStart := ui.Point{
		X: vanishingPoint.X + (center.X-vanishingPoint.X)*ratio,
		Y: vanishingPoint.Y + (center.Y-vanishingPoint.Y)*ratio,
	}
	rightStart := ui.Point{
		X: vanishingPoint.X + (right.X-vanishingPoint.X)*ratio,
		Y: vanishingPoint.Y + (right.Y-vanishingPoint.Y)*ratio,
	}

	// Simply use progress directly for linear movement
	leftCurrent := ui.Point{
		X: leftStart.X + (left.X-leftStart.X)*progress,
		Y: leftStart.Y + (left.Y-leftStart.Y)*progress,
	}
	centerCurrent := ui.Point{
		X: centerStart.X + (center.X-centerStart.X)*progress,
		Y: centerStart.Y + (center.Y-centerStart.Y)*progress,
	}
	rightCurrent := ui.Point{
		X: rightStart.X + (right.X-rightStart.X)*progress,
		Y: rightStart.Y + (right.Y-rightStart.Y)*progress,
	}
	var markerPath vector.Path
	markerPath.MoveTo(leftCurrent.X, leftCurrent.Y)
	markerPath.LineTo(centerCurrent.X, centerCurrent.Y)
	markerPath.LineTo(rightCurrent.X, rightCurrent.Y)

	vs, is := markerPath.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{
		Width: 1, // Thinner than regular notes
	})

	ui.ColorVertices(vs, markerColor)
	screen.DrawTriangles(vs, is, ui.BaseTriImg, nil)
}
