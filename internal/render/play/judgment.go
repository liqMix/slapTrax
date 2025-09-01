package play

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/liqmix/slaptrax/internal/display"
	"github.com/liqmix/slaptrax/internal/types"
	"github.com/liqmix/slaptrax/internal/ui"
	"github.com/liqmix/slaptrax/internal/user"
)

type NotePathOpts struct {
	lineWidth       float32
	isLarge         bool
	largeWidthRatio float32
	color           color.RGBA
	alpha           uint8
	solo            bool
}

// GetJudgementLinePath creates a path for the judgment line
func GetJudgementLinePath(track types.TrackName, pressed bool) *ui.CachedPath {
	return CreateJudgementPath(track, pressed)
}

func CreateNotePath(track types.TrackName, progress float32, opts *NotePathOpts) *ui.CachedPath {
	return CreateNotePathFromPoints(notePoints[track], progress, opts)
}

func CreateNotePathFromPoints(pts []*ui.Point, progress float32, opts *NotePathOpts) *ui.CachedPath {
	if len(pts) == 0 {
		return nil
	}

	notePath := vector.Path{}

	x, y := getPointPosition(pts[0], progress)
	notePath.MoveTo(x, y)

	for _, pt := range pts[1:] {
		if pt == nil {
			continue
		}
		x, y := getPointPosition(pt, progress)
		notePath.LineTo(x, y)
	}

	width := (opts.lineWidth * progress) / float32(display.Window.RenderScale()) * 1.5
	width *= user.S().NoteWidth
	if opts.isLarge {
		width *= opts.largeWidthRatio
	}

	vertices, indices := notePath.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{
		Width:    width,
		LineCap:  vector.LineCapRound,
		LineJoin: vector.LineJoinRound,
	})

	if !opts.solo {
		opts.color = types.White.C()
	}
	opts.color.A = opts.alpha

	ui.ColorVertices(vertices, opts.color)
	return &ui.CachedPath{
		Vertices: vertices,
		Indices:  indices,
	}
}

func getPointPosition(point *ui.Point, progress float32) (float32, float32) {
	pX, pY := point.ToRender32()
	cX, cY := playCenterPoint.ToRender32()
	x, y := cX+(pX-cX)*progress, cY+(pY-cY)*progress
	return x, y
}

// GetNotePath creates a note path for effects (simplified version)
func GetNotePath(track types.TrackName, note *types.Note, isEffect bool) *ui.CachedPath {
	progress := float32(note.Progress)
	alpha := GetNoteFadeAlpha(progress)
	if isEffect {
		alpha = uint8(200 * progress)
	}
	
	opts := &NotePathOpts{
		lineWidth:       getNoteWidth(),
		isLarge:         !note.Solo,
		largeWidthRatio: noteComboRatio,
		color:           track.NoteColor(),
		alpha:           alpha,
		solo:            note.Solo,
	}
	
	return CreateNotePath(track, progress, opts)
}

