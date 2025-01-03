package play

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
)

func (r *Play) drawNotes(screen *ebiten.Image) {
	for _, track := range r.state.Tracks {
		if len(track.ActiveNotes) == 0 {
			continue
		}

		color := track.Name.NoteColor()
		pts := GetNotePoints(track.Name)
		for _, note := range track.ActiveNotes {
			drawNote(screen, note, pts, color)
		}
	}
}

const (
	fadeInThreshold  = 0.02 // Start fade at
	fadeOutThreshold = 0.05 // Complete fade at
	maxAlpha         = 255
	minT             = 0.01 // Small value for vanishing point calculation
)

func drawNote(screen *ebiten.Image, note *types.Note, vec []*ui.Point, color color.RGBA) {
	if note == nil || note.Progress < 0 || note.Progress > 1 {
		return
	}

	// Calculate the constant screen-space velocity progress
	progress := (minT / (minT + (1-minT)*(1-float32(note.Progress))))

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

	leftX, leftY := vec[0].ToRender32()
	centerX, centerY := vec[1].ToRender32()
	rightX, rightY := vec[2].ToRender32()
	cX, cY := playCenterPoint.ToRender32()

	leftCurrentX := cX + (leftX-cX)*progress
	leftCurrentY := cY + (leftY-cY)*progress
	centerCurrentX := cX + (centerX-cX)*progress
	centerCurrentY := cY + (centerY-cY)*progress
	rightCurrentX := cX + (rightX-cX)*progress
	rightCurrentY := cY + (rightY-cY)*progress

	notePath := vector.Path{}
	notePath.MoveTo(leftCurrentX, leftCurrentY)
	notePath.LineTo(centerCurrentX, centerCurrentY)
	notePath.LineTo(rightCurrentX, rightCurrentY)

	width := noteWidth * float32(note.Progress)
	if !note.Solo {
		width *= 3.0
	}

	vs, is := notePath.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{
		Width:      width,
		LineCap:    vector.LineCapRound,
		LineJoin:   vector.LineJoinRound,
		MiterLimit: 1,
	})

	ui.ColorVertices(vs, color)
	screen.DrawTriangles(vs, is, ui.BaseTriImg, nil)

}
