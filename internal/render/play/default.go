package play

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/liqmix/ebiten-holiday-2024/internal/cache"
	"github.com/liqmix/ebiten-holiday-2024/internal/config"
	"github.com/liqmix/ebiten-holiday-2024/internal/l"
	"github.com/liqmix/ebiten-holiday-2024/internal/song"
	play "github.com/liqmix/ebiten-holiday-2024/internal/state/play"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
	"github.com/liqmix/ebiten-holiday-2024/internal/user"
)

// Normalized 2D point (0-1 range)
type Point struct {
	X, Y float32
}
type NormPoint struct {
	X, Y float64
}

// Convert normalized point to screen coordinates
func (p NormPoint) RenderPoint() Point {
	s := user.Settings()
	return Point{
		X: float32(p.X * float64(s.RenderWidth)),
		Y: float32(p.Y * float64(s.RenderHeight)),
	}
}

// Create a 1x1 white image as the base texture
var baseImg = ebiten.NewImage(1, 1)
var laneConfigs = map[song.TrackName]*LaneConfig{
	song.LeftBottom: {
		CurveAmount:    cornerCurve,
		VanishingPoint: mainCenter,
		Left: NormPoint{
			X: mainLeft,
			Y: mainCenter.Y + spacing,
		},
		Center: NormPoint{
			X: mainLeft,
			Y: mainBottom,
		},
		Right: NormPoint{
			X: mainCenter.X - mainBottomSpacing,
			Y: mainBottom,
		},
	},
	song.LeftTop: {
		CurveAmount:    cornerCurve,
		VanishingPoint: mainCenter,
		Left: NormPoint{
			X: mainCenter.X - mainTopSpacing,
			Y: mainTop,
		},
		Center: NormPoint{
			X: mainLeft,
			Y: mainTop,
		},
		Right: NormPoint{
			X: mainLeft,
			Y: mainCenter.Y - spacing,
		},
	},

	song.RightBottom: {
		CurveAmount:    cornerCurve,
		VanishingPoint: mainCenter,
		Left: NormPoint{
			X: mainCenter.X + mainBottomSpacing,
			Y: mainBottom,
		},
		Center: NormPoint{
			X: mainRight,
			Y: mainBottom,
		},
		Right: NormPoint{
			X: mainRight,
			Y: mainCenter.Y + spacing,
		},
	},
	song.RightTop: {
		CurveAmount:    cornerCurve,
		VanishingPoint: mainCenter,
		Left: NormPoint{
			X: mainRight,
			Y: mainCenter.Y - spacing,
		},
		Center: NormPoint{
			X: mainRight,
			Y: mainTop,
		},
		Right: NormPoint{
			X: mainCenter.X + mainTopSpacing,
			Y: mainTop,
		},
	},

	song.Center: {
		CurveAmount:    0,
		VanishingPoint: mainCenter,
		Left: NormPoint{
			X: mainCenter.X - centerTrackWidth/2,
			Y: mainBottom,
		},
		Center: NormPoint{
			X: mainCenter.X,
			Y: mainBottom,
		},
		Right: NormPoint{
			X: mainCenter.X + centerTrackWidth/2,
			Y: mainBottom,
		},
	},

	song.EdgeTop: {
		CurveAmount:    0,
		VanishingPoint: edgeCenter,
		Left: NormPoint{
			X: edgeLeft,
			Y: edgeTop,
		},
		Center: NormPoint{
			X: edgeCenter.X,
			Y: edgeTop,
		},
		Right: NormPoint{
			X: edgeRight,
			Y: edgeTop,
		},
	},
	song.EdgeTap1: {
		CurveAmount:    0,
		VanishingPoint: edgeCenter,
		Left: NormPoint{
			X: edgeLeft,
			Y: edgeBottom,
		},
		Center: NormPoint{
			X: (edgeLeft + edgeTapWidth/2),
			Y: edgeBottom,
		},
		Right: NormPoint{
			X: edgeLeft + edgeTapWidth,
			Y: edgeBottom,
		},
	},

	song.EdgeTap2: {
		CurveAmount:    0,
		VanishingPoint: edgeCenter,
		Left: NormPoint{
			X: edgeLeft + edgeTapWidth,
			Y: edgeBottom,
		},
		Center: NormPoint{
			X: (edgeLeft + edgeTapWidth + edgeTapWidth/2),
			Y: edgeBottom,
		},
		Right: NormPoint{
			X: edgeLeft + edgeTapWidth*2,
			Y: edgeBottom,
		},
	},

	song.EdgeTap3: {
		CurveAmount:    0,
		VanishingPoint: edgeCenter,
		Left: NormPoint{
			X: edgeLeft + edgeTapWidth*2,
			Y: edgeBottom,
		},
		Center: NormPoint{
			X: edgeLeft + edgeTapWidth*2 + edgeTapWidth/2,
			Y: edgeBottom,
		},
		Right: NormPoint{
			X: edgeLeft + edgeTapWidth*3,
			Y: edgeBottom,
		},
	},
}

var mainTracks = []song.TrackName{
	song.LeftBottom,
	song.LeftTop,
	song.RightBottom,
	song.RightTop,
	song.Center,
}

// type Animation struct {
// 	startTime int64
// }

// func (a *Animation) Draw() float32 {
// 	return 0
// }

// The default renderer for the play state.
type Default struct {
	state *play.State
	// animations map[string]*Animation
}

func (r Default) New(s *play.State) PlayRenderer {
	baseImg.Fill(white)
	return &Default{
		state: s,
		// animations: map[string]*Animation{},
	}
}

func (r *Default) Draw(screen *ebiten.Image) {
	r.drawBackground(screen)
	r.drawTracks(screen)
	r.drawProfile(screen)
	r.drawSongInfo(screen)
	r.drawScore(screen)
}

// TODO: later after tracks and notes
func (r *Default) drawProfile(screen *ebiten.Image)  {}
func (r *Default) drawSongInfo(screen *ebiten.Image) {}
func (r *Default) drawScore(screen *ebiten.Image) {
	// Draw the score at the top of the screen
	score := r.state.Score

	// Draw the score at the top of the screen
	s := user.Settings()
	x := 0.95 * float64(s.RenderWidth)
	y := 0.05 * float64(s.RenderHeight)

	perfectText := fmt.Sprintf(l.String(l.HIT_PERFECT)+": %d", score.Perfect)
	ui.DrawTextAt(screen, perfectText, int(x), int(y), 1)

	y += 20
	goodText := fmt.Sprintf(l.String(l.HIT_GOOD)+": %d", score.Good)
	ui.DrawTextAt(screen, goodText, int(x), int(y), 1)

	y += 20
	badText := fmt.Sprintf(l.String(l.HIT_BAD)+": %d", score.Bad)
	ui.DrawTextAt(screen, badText, int(x), int(y), 1)

	y += 20
	missText := fmt.Sprintf(l.String(l.HIT_MISS)+": %d", score.Miss)
	ui.DrawTextAt(screen, missText, int(x), int(y), 1)
}

func (r *Default) drawBackground(screen *ebiten.Image) {
	// If we've already created the background, or the render size hasn't changed
	s := user.Settings()

	bg, ok := cache.GetImage("play.background")
	if !ok {
		// Create the background image
		bg = ebiten.NewImage(s.RenderWidth, s.RenderHeight)
		// TODO: actually make some sort of background
		bg.Fill(color.Gray16{0x0000})
		cache.SetImage("play.background", bg)
	}
	screen.DrawImage(bg, nil)
}

func (r *Default) drawTracks(screen *ebiten.Image) {
	r.drawMainTracks(screen)
	r.drawEdgeTracks(screen)
	// r.drawMeasureMarkers(screen)
	r.drawNotes(screen)
}

var (
	offset      = 0.05
	spacing     = 0.0
	cornerCurve = 0.7

	// Main
	mainHeight = 0.60
	mainWidth  = 0.60
	mainLeft   = offset
	mainRight  = offset + mainWidth
	mainTop    = 1 - (offset + mainHeight)
	mainBottom = 1 - offset
	mainCenter = NormPoint{
		X: mainRight - (mainWidth / 2),
		Y: mainBottom - (mainHeight / 2),
	}
	mainTopSpacing    = spacing / 2
	mainBottomSpacing = (centerTrackWidth / 2) + spacing
	centerTrackWidth  = 0.3
	mainScoreScale    = 0.035

	// Edge
	edgeLeft   = mainRight + offset
	edgeRight  = 1 - offset
	edgeWidth  = edgeRight - edgeLeft
	edgeHeight = 0.8 * mainHeight
	edgeTop    = 1 - (offset + edgeHeight)
	edgeBottom = 1 - offset
	edgeCenter = NormPoint{
		X: edgeLeft + (edgeWidth / 2),
		Y: mainCenter.Y,
	}
	edgeTapWidth = edgeWidth / 3

	// Colors
	black     = color.RGBA{0, 0, 0, 255}
	gray      = color.RGBA{100, 100, 100, 255}
	white     = color.RGBA{200, 200, 200, 255}
	orange    = color.RGBA{255, 165, 0, 255}
	blue      = color.RGBA{0, 0, 255, 255}
	lightBlue = color.RGBA{173, 216, 230, 255}
	yellow    = color.RGBA{255, 255, 0, 255}

	// Note Width (as ratio to track width)
	targetThick = 4
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

func (r *Default) drawMainTracks(screen *ebiten.Image) {
	var laneConfig *LaneConfig
	var img *ebiten.Image
	var ok bool

	for _, track := range mainTracks {
		laneConfig = laneConfigs[track]

		img, ok = cache.GetImage(string(track))
		if !ok {
			img = r.createLaneBackground(laneConfig)
			cache.SetImage(string(track), img)
		}
		screen.DrawImage(img, nil)

		// Judgement line
		// Enlarge if track is active
		r.drawJudgementLine(screen, track)
	}

	// Then draw the area for combo / hit display
	// It's a box with rounded corners at the very center of main tracks
	// It will obscure some of the vanishing point lines
	s := user.Settings()
	img, ok = cache.GetImage("play.maincombo")
	if !ok {
		img = ebiten.NewImage(s.RenderWidth, s.RenderHeight)

		// Draw the box
		x := (mainCenter.X - (mainScoreScale / 2)) * float64(s.RenderWidth)
		y := (mainCenter.Y - (mainScoreScale / 2)) * float64(s.RenderHeight)
		width := mainScoreScale * float64(s.RenderWidth)
		height := mainScoreScale * float64(s.RenderHeight)

		// Fill
		border := 0.005 * float64(s.RenderHeight)

		// Border
		vector.DrawFilledRect(
			img,
			float32(x),
			float32(y),
			float32(width),
			float32(height),
			gray,
			true,
		)

		vector.DrawFilledRect(
			img,
			float32(x+border/2),
			float32(y+border/2),
			float32(width-border),
			float32(height-border),
			black,
			true,
		)
		cache.SetImage("play.maincombo", img)
	}
	screen.DrawImage(img, nil)

	// Draw the combo text in the center of the combo box
	combo := r.state.Score.Combo
	if combo > 0 {
		comboText := fmt.Sprintf("%d", r.state.Score.Combo)
		ui.DrawTextAt(screen, comboText, int(mainCenter.X*float64(s.RenderWidth)), int(mainCenter.Y*float64(s.RenderHeight)), 1)
	}
}

func (r *Default) drawEdgeTracks(screen *ebiten.Image) {
	// Draw the top track
	img, ok := cache.GetImage(string(song.EdgeTop))
	if !ok {
		s := user.Settings()
		img = ebiten.NewImage(s.RenderWidth, s.RenderHeight)
		leftHalf := r.createLaneBackground(&LaneConfig{
			CurveAmount:    0,
			VanishingPoint: edgeCenter,
			Left: NormPoint{
				X: edgeCenter.X,
				Y: edgeTop,
			},
			Center: NormPoint{
				X: edgeLeft,
				Y: edgeTop,
			},
			Right: NormPoint{
				X: edgeLeft,
				Y: edgeCenter.Y,
			},
		})
		rightHalf := r.createLaneBackground(&LaneConfig{
			CurveAmount:    0,
			VanishingPoint: edgeCenter,
			Left: NormPoint{
				X: edgeRight,
				Y: edgeCenter.Y,
			},
			Center: NormPoint{
				X: edgeRight,
				Y: edgeTop,
			},
			Right: NormPoint{
				X: edgeCenter.X,
				Y: edgeTop,
			},
		})
		img.DrawImage(leftHalf, nil)
		img.DrawImage(rightHalf, nil)
		cache.SetImage(string(song.EdgeTop), img)
	}
	screen.DrawImage(img, nil)
	r.drawJudgementLine(screen, song.EdgeTop)

	// Draw the bottom track
	img, ok = cache.GetImage(string(song.EdgeTap1))
	if !ok {
		s := user.Settings()
		img = ebiten.NewImage(s.RenderWidth, s.RenderHeight)
		img.DrawImage(r.createLaneBackground(laneConfigs[song.EdgeTap1]), nil)
		img.DrawImage(r.createLaneBackground(laneConfigs[song.EdgeTap2]), nil)
		img.DrawImage(r.createLaneBackground(laneConfigs[song.EdgeTap3]), nil)

		cache.SetImage(string(song.EdgeTap1), img)
	}
	screen.DrawImage(img, nil)
	r.drawJudgementLine(screen, song.EdgeTap1)
	r.drawJudgementLine(screen, song.EdgeTap2)
	r.drawJudgementLine(screen, song.EdgeTap3)
}

type LaneConfig struct {
	Left           NormPoint // Bottom left point (0-1 space)
	Center         NormPoint // Bottom center point (0-1 space)
	Right          NormPoint // Bottom right point (0-1 space)
	VanishingPoint NormPoint // Convergence point (0-1 space)
	CurveAmount    float64   // 0 = right angle at center, 1 = straight line
	NoteWidth      float32   // Width of notes as a ratio of track width
	NoteColor      color.RGBA
}

func getVectorPath(points []Point, curveAmount float64) vector.Path {
	path := vector.Path{}
	if len(points) < 1 {
		return path
	}

	start := points[0]
	path.MoveTo(start.X, start.Y)

	for i := 1; i < len(points); i++ {
		next := points[i]
		path.LineTo(next.X, next.Y)
	}
	return path
}

// Get the rendering of the judgement line
func (r *Default) getJudgementLine(config *LaneConfig) vector.Path {
	left := config.Left.RenderPoint()
	right := config.Right.RenderPoint()
	center := config.Center.RenderPoint()

	if left.X == right.X || left.Y == right.Y {
		return getVectorPath([]Point{
			left,
			right,
		}, config.CurveAmount)
	}

	return getVectorPath([]Point{
		left,
		center,
		right,
	}, config.CurveAmount)
}

func (r *Default) drawJudgementLine(screen *ebiten.Image, trackName song.TrackName) {
	config, ok := laneConfigs[trackName]
	if !ok {
		return
	}
	img := r.getJudgementLine(config)
	var width float32 = 2.0
	if r.state.IsTrackPressed(trackName) {
		width = 8.0
	}
	vs, is := img.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{
		Width: width,
	})
	colorVertices(vs, getNoteColor(trackName))
	screen.DrawTriangles(vs, is, baseImg, nil)
}

func (r *Default) createLaneBackground(config *LaneConfig) *ebiten.Image {
	s := user.Settings()
	img := ebiten.NewImage(s.RenderWidth, s.RenderHeight)

	// left := config.Left.RenderPoint()
	// right := config.Right.RenderPoint()
	center := config.Center.RenderPoint()
	vanishingPoint := config.VanishingPoint.RenderPoint()

	// Draw the target line
	// targetLine := r.getJudgementLine(config)

	// Draw the guide lines to vanishing point
	// guideLineLeft := getVectorPath([]Point{
	// 	left,
	// 	vanishingPoint,
	// }, 0)

	// guideLineRight := getVectorPath([]Point{
	// 	right,
	// 	vanishingPoint,
	// }, 0)

	var guideColor = color.RGBA{
		R: gray.R,
		G: gray.G,
		B: gray.B,
		A: 150,
	}

	r.drawDashedLine(img,
		center,
		vanishingPoint,
		10,
		10,
		guideColor,
	)

	// vs, is = guideLineLeft.AppendVerticesAndIndicesForStroke(nil, nil, &opts)
	// colorVertices(vs, guideColor)
	// img.DrawTriangles(vs, is, baseImg, nil)

	// vs, is = guideLineRight.AppendVerticesAndIndicesForStroke(vs, is, &opts)
	// colorVertices(vs, guideColor)
	// img.DrawTriangles(vs, is, baseImg, nil)
	return img
}

func colorVertices(vs []ebiten.Vertex, color color.RGBA) {
	for i := range vs {
		vs[i].ColorR = float32(color.R) / 255
		vs[i].ColorG = float32(color.G) / 255
		vs[i].ColorB = float32(color.B) / 255
		vs[i].ColorA = float32(color.A) / 255
	}
}

func (r *Default) drawDashedLine(img *ebiten.Image, start Point, end Point, dashLength float32, gapLength float32, color color.RGBA) {
	dx := end.X - start.X
	dy := end.Y - start.Y
	length := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	// Normalize direction vector
	nx := dx / length
	ny := dy / length

	// Calculate number of segments
	totalLength := dashLength + gapLength
	segments := int(length / totalLength)

	// Draw dash segments
	for i := 0; i < segments; i++ {
		dashStart := Point{
			X: start.X + nx*float32(i)*totalLength,
			Y: start.Y + ny*float32(i)*totalLength,
		}
		dashEnd := Point{
			X: dashStart.X + nx*dashLength,
			Y: dashStart.Y + ny*dashLength,
		}
		path := getVectorPath([]Point{dashStart, dashEnd}, 0)

		vs, is := path.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{
			Width: 1,
		})

		colorVertices(vs, color)
		img.DrawTriangles(vs, is, baseImg, nil)
	}

	// Draw final dash if there's remaining space
	remainingLength := length - float32(segments)*totalLength
	if remainingLength > 0 && remainingLength > dashLength {
		finalStart := Point{
			X: start.X + nx*float32(segments)*totalLength,
			Y: start.Y + ny*float32(segments)*totalLength,
		}
		finalEnd := Point{
			X: finalStart.X + nx*dashLength,
			Y: finalStart.Y + ny*dashLength,
		}

		path := getVectorPath([]Point{finalStart, finalEnd}, 0)

		vs, is := path.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{
			Width: 1,
		})

		colorVertices(vs, color)
		img.DrawTriangles(vs, is, baseImg, nil)
	}
}

func (r *Default) drawNotes(screen *ebiten.Image) {
	for _, track := range r.state.Tracks {
		if len(track.ActiveNotes) == 0 {
			continue
		}

		// If we have the lane config in the map, draw the notes
		laneConfig, ok := laneConfigs[track.Name]
		if !ok {
			continue
		}

		color := getNoteColor(track.Name)
		for _, note := range track.ActiveNotes {
			r.drawNote(screen, track.Name, note, laneConfig, color)
		}
	}
}

func (r *Default) drawNote(screen *ebiten.Image, trackName song.TrackName, note *song.Note, config *LaneConfig, noteColor color.RGBA) {
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
	leftStart := Point{
		X: vanishingPoint.X + (left.X-vanishingPoint.X)*ratio,
		Y: vanishingPoint.Y + (left.Y-vanishingPoint.Y)*ratio,
	}
	centerStart := Point{
		X: vanishingPoint.X + (center.X-vanishingPoint.X)*ratio,
		Y: vanishingPoint.Y + (center.Y-vanishingPoint.Y)*ratio,
	}
	rightStart := Point{
		X: vanishingPoint.X + (right.X-vanishingPoint.X)*ratio,
		Y: vanishingPoint.Y + (right.Y-vanishingPoint.Y)*ratio,
	}

	// Simply use progress directly for linear movement
	leftCurrent := Point{
		X: leftStart.X + (left.X-leftStart.X)*progress,
		Y: leftStart.Y + (left.Y-leftStart.Y)*progress,
	}
	centerCurrent := Point{
		X: centerStart.X + (center.X-centerStart.X)*progress,
		Y: centerStart.Y + (center.Y-centerStart.Y)*progress,
	}
	rightCurrent := Point{
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

	colorVertices(vs, c)
	screen.DrawTriangles(vs, is, baseImg, nil)
}

func (r *Default) drawMeasureMarkers(screen *ebiten.Image) {
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
			for _, track := range mainTracks {
				if config, ok := laneConfigs[track]; ok {
					r.drawMeasureMarker(screen, float32(progress), config, markerColor)
				}
			}
		}
	}
}

func (r *Default) drawMeasureMarker(screen *ebiten.Image, p float32, config *LaneConfig, markerColor color.RGBA) {
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
	leftStart := Point{
		X: vanishingPoint.X + (left.X-vanishingPoint.X)*ratio,
		Y: vanishingPoint.Y + (left.Y-vanishingPoint.Y)*ratio,
	}
	centerStart := Point{
		X: vanishingPoint.X + (center.X-vanishingPoint.X)*ratio,
		Y: vanishingPoint.Y + (center.Y-vanishingPoint.Y)*ratio,
	}
	rightStart := Point{
		X: vanishingPoint.X + (right.X-vanishingPoint.X)*ratio,
		Y: vanishingPoint.Y + (right.Y-vanishingPoint.Y)*ratio,
	}

	// Simply use progress directly for linear movement
	leftCurrent := Point{
		X: leftStart.X + (left.X-leftStart.X)*progress,
		Y: leftStart.Y + (left.Y-leftStart.Y)*progress,
	}
	centerCurrent := Point{
		X: centerStart.X + (center.X-centerStart.X)*progress,
		Y: centerStart.Y + (center.Y-centerStart.Y)*progress,
	}
	rightCurrent := Point{
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

	colorVertices(vs, markerColor)
	screen.DrawTriangles(vs, is, baseImg, nil)
}
