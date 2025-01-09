package play

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/liqmix/ebiten-holiday-2024/internal/cache"
	"github.com/liqmix/ebiten-holiday-2024/internal/display"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
)

func (r *Play) addNotePath(track *types.Track) {
	if len(track.ActiveNotes) == 0 {
		return
	}

	for _, note := range track.ActiveNotes {
		if note.IsHoldNote() {
			r.addHoldPaths(track.Name, note)
		}
		path := GetNotePath(track.Name, note, false)
		if path != nil {
			r.vectorCollection.AddPath(path)
		}
	}
}

func getPointPosition(point *ui.Point, progress float32) (float32, float32) {
	pX, pY := point.ToRender32()
	cX, cY := playCenterPoint.ToRender32()
	x, y := cX+(pX-cX)*progress, cY+(pY-cY)*progress
	return x, y
}

type NotePathOpts struct {
	lineWidth       float32
	isLarge         bool
	largeWidthRatio float32
	color           color.RGBA
	alpha           uint8
	solo            bool
}

func CreateNotePath(track types.TrackName, progress float32, opts *NotePathOpts) *cache.CachedPath {
	return CreateNotePathFromPoints(notePoints[track], progress, opts)
}

func CreateNotePathFromPoints(pts []*ui.Point, progress float32, opts *NotePathOpts) *cache.CachedPath {
	if len(pts) == 0 {
		return nil
	}

	// Create initial path for base shape
	notePath := vector.Path{}

	// Store original points for depth extrusion
	originalPoints := make([][2]float32, len(pts))

	// Initialize starting point
	x, y := getPointPosition(pts[0], progress)
	notePath.MoveTo(x, y)
	originalPoints[0] = [2]float32{x, y}

	// Add line segments and store original points
	for i, pt := range pts[1:] {
		if pt == nil {
			continue
		}
		x, y := getPointPosition(pt, progress)
		notePath.LineTo(x, y)
		originalPoints[i+1] = [2]float32{x, y}
	}

	// Calculate line width with depth consideration
	width := (opts.lineWidth * progress) / float32(display.Window.RenderScale()) * 1.5
	if opts.isLarge {
		width *= opts.largeWidthRatio
	}

	// Create base vertices and indices
	vertices, indices := notePath.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{
		Width:    width,
		LineCap:  vector.LineCapRound,
		LineJoin: vector.LineJoinRound,
	})

	if !opts.solo {
		opts.color = types.White.C()
	}
	// Set color with transparency
	opts.color.A = opts.alpha

	ui.ColorVertices(vertices, opts.color)
	cachedPath := &cache.CachedPath{
		Vertices: vertices,
		Indices:  indices,
	}
	return cachedPath
}

func (r *Play) addHoldPaths(track types.TrackName, note *types.Note) {
	notePoints := notePoints[track]
	holdPaths := [3]vector.Path{}

	if note.Progress < 0.25 {
		return
	}
	progress := SmoothProgress(note.Progress)
	releaseProgress := SmoothProgress(note.ReleaseProgress)

	// Get last point positions
	var x, y float32
	var endX, endY float32

	for i, pt := range notePoints {
		if pt == nil {
			continue
		}
		x, y = getPointPosition(notePoints[i], progress)
		endX, endY = getPointPosition(notePoints[i], releaseProgress)
		holdPaths[i].MoveTo(x, y)
		holdPaths[i].LineTo(endX, endY)
		if i > 0 {
			holdPaths[i].LineTo(getPointPosition(notePoints[i-1], releaseProgress))
			holdPaths[i].LineTo(getPointPosition(notePoints[i-1], progress))
		}
	}

	// Generate vertices and indices for stroke
	for _, holdPath := range holdPaths {
		vs, is := holdPath.AppendVerticesAndIndicesForFilling(nil, nil)

		// Set color with transparency
		color := track.NoteColor()
		if note.WasHit() {
			if note.WasReleased() {
				color.A = 50
			} else {
				color.A = 200
			}
		} else {
			color.A = 100
		}
		ui.ColorVertices(vs, color)
		r.vectorCollection.Add(vs, is)
	}
}
