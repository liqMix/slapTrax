package play

// func (r *Play) drawNotes(screen *ebiten.Image) {
// 	vectorCollection := ui.NewVectorCollection()

// }

// func (r *Play) drawCachedNote(screen *ebiten.Image, track types.TrackName, note *types.Note, color color.RGBA) *cachedPath {
// 	cachedPath := r.vectorCache.GetNotePath(track, note)
// 	if cachedPath == nil {
// 		return nil
// 	}
// 	return cachedPath
// }

// func drawNote(screen *ebiten.Image, note *types.Note, vec []*ui.Point, color color.RGBA) {
// 	if note == nil || note.Progress < 0 || note.Progress > 1 {
// 		return
// 	}

// 	progress := SmoothProgress(note.Progress)
// 	color.A = GetFadeAlpha(progress)

// 	// Skip rendering if fully transparent
// 	if color.A == 0 {
// 		return
// 	}

// 	notePath := vector.Path{}
// 	cX, cY := playCenterPoint.ToRender32()

// 	startX, startY := vec[0].ToRender32()
// 	x, y := cX+(startX-cX)*progress, cY+(startY-cY)*progress
// 	notePath.MoveTo(x, y)

// 	for i := 1; i < len(vec); i++ {
// 		if vec[i] == nil {
// 			continue
// 		}
// 		x, y = vec[i].ToRender32()
// 		x, y = cX+(x-cX)*progress, cY+(y-cY)*progress
// 		notePath.LineTo(x, y)
// 	}

// 	width := noteWidth * float32(note.Progress)
// 	if !note.Solo {
// 		width *= 3.0
// 	}

// 	vs, is := notePath.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{
// 		Width:    width,
// 		LineCap:  vector.LineCapRound,
// 		LineJoin: vector.LineJoinRound,
// 	})

// 	ui.ColorVertices(vs, color)
// 	screen.DrawTriangles(vs, is, ui.BaseTriImg, nil)

// }
