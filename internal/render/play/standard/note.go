package standard

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/liqmix/ebiten-holiday-2024/internal/song"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
)

// TODO: make user configurable
func getNoteColor(trackName song.TrackName) color.RGBA {
	switch trackName {
	case song.LeftBottom:
		return orange
	case song.LeftTop:
		return orange
	case song.RightBottom:
		return orange
	case song.RightTop:
		return orange
	case song.Center:
		return yellow
	case song.EdgeTop:
		return lightBlue
	case song.EdgeTap1:
		return white
	case song.EdgeTap2:
		return white
	case song.EdgeTap3:
		return white
	}
	return gray
}

func (r *Standard) drawNotes(screen *ebiten.Image) {
	for _, track := range r.state.Tracks {
		if len(track.ActiveNotes) == 0 {
			continue
		}

		if !r.renderEdgeTracks {
			if track.Name == song.EdgeTop || track.Name == song.EdgeTap1 || track.Name == song.EdgeTap2 || track.Name == song.EdgeTap3 {
				continue
			}
		}
		// If we have the lane config in the map, draw the notes
		laneConfig, ok := laneConfigs[track.Name]
		if !ok {
			continue
		}

		color := getNoteColor(track.Name)
		for _, note := range track.ActiveNotes {
			r.drawNote(screen, note, laneConfig, color)
		}
	}
}

func (r *Standard) drawNote(screen *ebiten.Image, note *song.Note, config *LaneConfig, noteColor color.RGBA) {
	if note == nil || note.Progress < 0 || note.Progress > 1 {
		return
	}

	left := config.Left.RenderPoint()
	right := config.Right.RenderPoint()
	center := config.Center.RenderPoint()
	vanishingPoint := config.VanishingPoint.RenderPoint()

	// To achieve constant screen-space velocity, we need to use inverse interpolation
	// The formula is: p = 1 / (1 + (1/t - 1) * progress) where t is close to 0 (like 0.01)
	// This gives us perfect constant speed in screen space
	t := float32(0.01) // Small value for the starting point relative to vanishing point
	p := float32(note.Progress)
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

	// Draw the note
	var notePath vector.Path
	notePath.MoveTo(leftCurrent.X, leftCurrent.Y)
	notePath.LineTo(centerCurrent.X, centerCurrent.Y)
	notePath.LineTo(rightCurrent.X, rightCurrent.Y)

	// Create the vertices and indices for the note path with dynamic width
	vs, is := notePath.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{
		Width: 3,
	})

	c := noteColor
	if progress < 0.05 {
		fadeInFactor := progress
		c.A = uint8(float32(noteColor.A) * fadeInFactor)
	}

	ui.ColorVertices(vs, c)
	screen.DrawTriangles(vs, is, ui.BaseTriImg, nil)
}
