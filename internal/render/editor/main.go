package editor

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/liqmix/slaptrax/internal/display"
	"github.com/liqmix/slaptrax/internal/render/play"
	"github.com/liqmix/slaptrax/internal/state"
	"github.com/liqmix/slaptrax/internal/types"
	"github.com/liqmix/slaptrax/internal/ui"
)

type EditorRenderer struct {
	display.BaseRenderer
	state            *state.EditorState
	vectorCollection *ui.VectorCollection
	
	// Track layout from play renderer
	notePoints  [][]*ui.Point
	playArea    playAreaLayout
}

type playAreaLayout struct {
	centerX      float64
	centerY      float64
	centerPoint  ui.Point
	left         float64
	right        float64
	top          float64
	bottom       float64
	width        float64
	height       float64
	noteLength   float64
}

func NewEditorRenderer(s *state.EditorState) *EditorRenderer {
	renderer := &EditorRenderer{
		state:            s,
		vectorCollection: ui.NewVectorCollection(),
	}
	renderer.initLayout()
	renderer.BaseRenderer.Init(renderer.static)
	
	// Initialize vector cache for note rendering
	play.RebuildVectorCache()
	
	return renderer
}

func (r *EditorRenderer) initLayout() {
	// Use EXACT same layout calculations as play renderer
	windowOffset := 0.025
	minX := 0.0 + windowOffset
	maxY := 1.0 - (windowOffset * 2)
	maxX := 1.0 - (windowOffset * 2)
	centerX := minX + ((maxX - minX) / 2)
	
	// Exact same play area calculations as play renderer
	playWidth := maxX * 0.75
	playCenterX := centerX
	playHeight := 0.6
	playBottom := maxY - windowOffset
	playTop := playBottom - playHeight
	playLeft := centerX - (playWidth / 2)
	playRight := centerX + (playWidth / 2)
	playCenterY := playTop + (playHeight / 2)
	noteLength := playWidth * 0.25
	
	r.playArea = playAreaLayout{
		centerX:     playCenterX,
		centerY:     playCenterY,
		centerPoint: ui.Point{X: playCenterX, Y: playCenterY},
		left:        playLeft,
		right:       playRight,
		top:         playTop,
		bottom:      playBottom,
		width:       playWidth,
		height:      playHeight,
		noteLength:  noteLength,
	}
	
	// Initialize note points for each track (same as play renderer)
	r.notePoints = r.createNotePoints(noteLength)
}

func (r *EditorRenderer) createNotePoints(length float64) [][]*ui.Point {
	centerLength := length * 0.5 // centerNoteLengthRatio
	
	// Order must match types.TrackNames() order:
	// TrackLeftBottom, TrackLeftTop, TrackRightBottom, TrackRightTop, TrackCenterBottom, TrackCenterTop
	return [][]*ui.Point{
		// TrackLeftBottom (index 0)
		{
			&ui.Point{X: r.playArea.left, Y: r.playArea.bottom - length},
			&ui.Point{X: r.playArea.left, Y: r.playArea.bottom},
			&ui.Point{X: r.playArea.left + length, Y: r.playArea.bottom},
		},
		// TrackLeftTop (index 1)
		{
			&ui.Point{X: r.playArea.left, Y: r.playArea.top + length},
			&ui.Point{X: r.playArea.left, Y: r.playArea.top},
			&ui.Point{X: r.playArea.left + length, Y: r.playArea.top},
		},
		// TrackRightBottom (index 2)
		{
			&ui.Point{X: r.playArea.right - length, Y: r.playArea.bottom},
			&ui.Point{X: r.playArea.right, Y: r.playArea.bottom},
			&ui.Point{X: r.playArea.right, Y: r.playArea.bottom - length},
		},
		// TrackRightTop (index 3)
		{
			&ui.Point{X: r.playArea.right - length, Y: r.playArea.top},
			&ui.Point{X: r.playArea.right, Y: r.playArea.top},
			&ui.Point{X: r.playArea.right, Y: r.playArea.top + length},
		},
		// TrackCenterBottom (index 4)
		{
			&ui.Point{X: r.playArea.centerX - centerLength, Y: r.playArea.bottom},
			&ui.Point{X: r.playArea.centerX, Y: r.playArea.bottom},
			&ui.Point{X: r.playArea.centerX + centerLength, Y: r.playArea.bottom},
		},
		// TrackCenterTop (index 5)
		{
			&ui.Point{X: r.playArea.centerX - centerLength, Y: r.playArea.top},
			&ui.Point{X: r.playArea.centerX, Y: r.playArea.top},
			&ui.Point{X: r.playArea.centerX + centerLength, Y: r.playArea.top},
		},
	}
}

func (r *EditorRenderer) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	r.BaseRenderer.Draw(screen, opts)
	
	// Draw help text at top
	r.drawHelpText(screen)
	
	// Draw play area boundary
	r.drawPlayAreaBounds(screen)
	
	// Draw track lanes
	r.drawTrackLanes(screen)
	
	// Draw notes using same system as play renderer
	r.drawNotes(screen)
	
	// Draw cursor and editing indicators
	r.drawEditingIndicators(screen)
	
	// Draw vector collection (notes and effects)
	r.vectorCollection.Draw(screen)
	r.vectorCollection.Clear()
}

func (r *EditorRenderer) static(img *ebiten.Image, opts *ebiten.DrawImageOptions) {
	// Draw background
	img.Fill(types.Black.C())
}

func (r *EditorRenderer) drawHelpText(screen *ebiten.Image) {
	helpText := "QWE: Top tracks | ASD: Bottom tracks | P: Play | ←→: Time | ↑↓: Division | Alt+↑↓: Speed | Alt±: Audio Offset | K: Audio Start | I: Countdown Beat | B: Beginning | Shift+B: Last Note | G: Grid | Shift±: BPM | Ctrl+S: Save | Ctrl+Shift+O: Import Audio"
	ebitenutil.DebugPrintAt(screen, helpText, 10, 10)
	
	// Display BPM, time division, lane speed, audio offset, and current position info
	infoText := fmt.Sprintf("BPM: %.1f | Division: 1/%d | Speed: %.1fx | Audio Offset: %+dms | Position: %s", 
		r.state.GetBPM(), 
		r.state.GetTimeDivision(),
		r.state.GetLaneSpeed(), 
		r.state.GetAudioOffset(),
		r.state.GetCurrentTimePosition())
	ebitenutil.DebugPrintAt(screen, infoText, 10, 30)
}

func (r *EditorRenderer) drawPlayAreaBounds(screen *ebiten.Image) {
	// Draw play area boundary
	leftPt := ui.Point{X: r.playArea.left, Y: r.playArea.top}
	rightPt := ui.Point{X: r.playArea.right, Y: r.playArea.bottom}
	topLeftX, topLeftY := leftPt.ToRender()
	bottomRightX, bottomRightY := rightPt.ToRender()
	
	width := float64(bottomRightX - topLeftX)
	height := float64(bottomRightY - topLeftY)
	
	// Draw boundary rectangle
	boundaryColor := color.RGBA{128, 128, 128, 100}
	ebitenutil.DrawRect(screen, float64(topLeftX), float64(topLeftY), 2, height, boundaryColor)
	ebitenutil.DrawRect(screen, float64(bottomRightX-2), float64(topLeftY), 2, height, boundaryColor)
	ebitenutil.DrawRect(screen, float64(topLeftX), float64(topLeftY), width, 2, boundaryColor)
	ebitenutil.DrawRect(screen, float64(topLeftX), float64(bottomRightY-2), width, 2, boundaryColor)
}

func (r *EditorRenderer) drawTrackLanes(screen *ebiten.Image) {
	tracks := types.TrackNames()
	
	for i, trackName := range tracks {
		points := r.notePoints[i]
		isSelected := trackName == r.state.SelectedTrack()
		
		r.addTrackLanePath(trackName, points, isSelected)
		// Removed addJudgementLinePath - no longer needed with QWE/ASD direct note toggling
	}
}

func (r *EditorRenderer) addTrackLanePath(trackName types.TrackName, points []*ui.Point, isSelected bool) {
	if len(points) == 0 {
		return
	}
	
	centerX, centerY := r.playArea.centerPoint.ToRender32()
	trackPath := vector.Path{}
	
	// Get starting point dimensions
	startX, startY := points[0].ToRender32()
	width := float32(2) // Lane width
	
	if isSelected {
		width = float32(4) // Wider for selected track
	}
	
	// Calculate perpendicular vector for width
	dx := startX - centerX
	dy := startY - centerY
	length := float32(distance32(dx, dy))
	normalX := -dy / length * width
	normalY := dx / length * width
	
	// Create lane outline
	trackPath.MoveTo(centerX+(startX-centerX)+normalX, centerY+(startY-centerY)+normalY)
	
	// Add right edge line segments
	for _, pt := range points[1:] {
		if pt == nil {
			continue
		}
		x, y := pt.ToRender32()
		dx = x - centerX
		dy = y - centerY
		length = float32(distance32(dx, dy))
		normalX = -dy / length * width
		normalY = dx / length * width
		
		trackPath.LineTo(centerX+(x-centerX)+normalX, centerY+(y-centerY)+normalY)
	}
	
	// Move to center
	trackPath.LineTo(centerX, centerY)
	
	// Add left edge in reverse
	for i := len(points) - 1; i >= 0; i-- {
		if points[i] == nil {
			continue
		}
		x, y := points[i].ToRender32()
		dx = x - centerX
		dy = y - centerY
		length = float32(distance32(dx, dy))
		normalX = -dy / length * width
		normalY = dx / length * width
		
		trackPath.LineTo(centerX+(x-centerX)-normalX, centerY+(y-centerY)-normalY)
	}
	
	trackPath.Close()
	
	// Generate vertices and indices
	vs, is := trackPath.AppendVerticesAndIndicesForFilling(nil, nil)
	
	// Set color
	trackColor := trackName.NoteColor()
	if isSelected {
		trackColor.A = 30
	} else {
		trackColor.A = 15
	}
	
	ui.ColorVertices(vs, trackColor)
	r.vectorCollection.Add(vs, is)
}

func (r *EditorRenderer) addJudgementLinePath(trackName types.TrackName, isSelected bool) {
	points := r.notePoints[int(trackName)]
	if len(points) == 0 {
		return
	}
	
	// Draw judgement line at the end of the track
	judgementPath := vector.Path{}
	
	// Use the first point (target area) for judgement line
	x, y := points[0].ToRender32()
	centerX, centerY := r.playArea.centerPoint.ToRender32()
	
	// Calculate direction vector
	dx := x - centerX
	dy := y - centerY
	length := float32(distance32(dx, dy))
	normalX := -dy / length
	normalY := dx / length
	
	width := float32(20)
	if isSelected {
		width = float32(30)
	}
	
	// Create judgement line
	judgementPath.MoveTo(x+normalX*width, y+normalY*width)
	judgementPath.LineTo(x-normalX*width, y-normalY*width)
	
	// Generate stroke
	vs, is := judgementPath.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{
		Width:    float32(6),
		LineCap:  vector.LineCapRound,
		LineJoin: vector.LineJoinRound,
	})
	
	// Set color
	lineColor := trackName.NoteColor()
	if isSelected {
		lineColor.A = 200
	} else {
		lineColor.A = 100
	}
	
	ui.ColorVertices(vs, lineColor)
	r.vectorCollection.Add(vs, is)
}

func (r *EditorRenderer) drawNotes(screen *ebiten.Image) {
	currentTime := r.state.CurrentTime()
	// Apply lane speed to travel time (exactly matching play renderer)
	baseTravelTime := 10000.0 
	travelTime := int64(baseTravelTime / r.state.GetLaneSpeed())
	
	tracks := types.TrackNames()
	for _, trackName := range tracks {
		notes := r.state.GetNotesForTrack(trackName)
		
		for _, note := range notes {
			// Update note progress exactly like gameplay (types.Note.Update)
			note.Progress = math.Max(0, 1-float64(note.Target-currentTime)/float64(travelTime))
			if note.IsHoldNote() {
				note.ReleaseProgress = math.Max(0, 1-float64(note.TargetRelease-currentTime)/float64(travelTime))
			}
			
			// Use exact same visibility check as play renderer
			if note.Progress < 0 || note.Progress > 1 {
				continue
			}
			
			// Use cached path system exactly like gameplay
			path := play.GetNotePath(trackName, note, false)
			if path != nil {
				r.vectorCollection.AddPath(path)
			}
		}
	}
}


func (r *EditorRenderer) getPointPosition(point *ui.Point, progress float32) (float32, float32) {
	pX, pY := point.ToRender32()
	cX, cY := r.playArea.centerPoint.ToRender32()
	x, y := cX+(pX-cX)*progress, cY+(pY-cY)*progress
	return x, y
}

func (r *EditorRenderer) drawEditingIndicators(screen *ebiten.Image) {
	currentTime := r.state.CurrentTime()
	
	// Draw current time indicator (vertical line across play area)
	r.drawTimeIndicator(screen, currentTime, color.RGBA{255, 0, 0, 200})
	
	// Draw grid lines if enabled
	if r.state.ShowGrid() {
		r.drawGrid(screen, currentTime)
	}
	
	// Draw event markers (only when not playing - hide during gameplay simulation)
	if !r.state.IsPlaying() {
		r.drawEventMarkers(screen, currentTime)
	}
	
	// Draw hold note preview if creating one
	if r.state.IsHolding() && r.state.HoldStartTime() > 0 {
		r.drawHoldPreview(screen, currentTime)
	}
}

func (r *EditorRenderer) drawTimeIndicator(screen *ebiten.Image, time int64, clr color.RGBA) {
	// Draw a simple vertical line across the entire play area
	indicatorPath := vector.Path{}
	
	// Draw from top to bottom of play area at center
	centerX := float32(r.playArea.centerX)
	topY := float32(r.playArea.top)
	bottomY := float32(r.playArea.bottom)
	
	indicatorPath.MoveTo(centerX, topY)
	indicatorPath.LineTo(centerX, bottomY)
	
	vs, is := indicatorPath.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{
		Width:    float32(4),
		LineCap:  vector.LineCapRound,
		LineJoin: vector.LineJoinRound,
	})
	
	ui.ColorVertices(vs, clr)
	r.vectorCollection.Add(vs, is)
}

func (r *EditorRenderer) drawGrid(screen *ebiten.Image, currentTime int64) {
	// Apply lane speed to travel time (exactly matching play renderer)
	baseTravelTime := 10000.0
	travelTime := int64(baseTravelTime / r.state.GetLaneSpeed())
	
	// Convert current time to beats for precise alignment
	currentBeat := r.state.MsIntToBeats(currentTime)
	
	// Draw eighth note markers (faintest) - 1/8 notes
	for i := int64(0); i < 32; i++ {
		// Calculate beat position precisely - start from nearest eighth beat boundary
		eighthBeat := math.Floor(currentBeat*2)/2 + float64(i)/2
		eighthTimeMs := r.state.BeatsToMsInt(eighthBeat)
		
		progress := math.Max(0, 1-float64(eighthTimeMs-currentTime)/float64(travelTime))
		if progress > 0 && progress <= 1.0 {
			r.drawMeasureMarker(screen, progress, color.RGBA{32, 32, 32, 50}) // Faintest
		}
	}
	
	// Draw quarter note markers (medium) - 1/4 notes
	for i := int64(0); i < 16; i++ {
		// Calculate beat position precisely - start from nearest quarter beat boundary
		quarterBeat := math.Floor(currentBeat) + float64(i)
		quarterTimeMs := r.state.BeatsToMsInt(quarterBeat)
		
		progress := math.Max(0, 1-float64(quarterTimeMs-currentTime)/float64(travelTime))
		if progress > 0 && progress <= 1.0 {
			r.drawMeasureMarker(screen, progress, color.RGBA{64, 64, 64, 100}) // Medium visibility
		}
	}
	
	// Draw whole note/measure markers (brightest) - 1/1 notes  
	for i := int64(0); i < 4; i++ {
		// Calculate measure position precisely - start from nearest measure boundary
		wholeBeat := math.Floor(currentBeat/4)*4 + float64(i*4)
		wholeTimeMs := r.state.BeatsToMsInt(wholeBeat)
		
		progress := math.Max(0, 1-float64(wholeTimeMs-currentTime)/float64(travelTime))
		if progress > 0 && progress <= 1.0 {
			r.drawMeasureMarker(screen, progress, color.RGBA{128, 128, 128, 200}) // Most prominent
		}
	}
}

func (r *EditorRenderer) drawMeasureMarker(screen *ebiten.Image, p float64, clr color.RGBA) {
	if p < 0 || p > 1 {
		return
	}
	
	// Apply smooth progress (exactly matching play renderer)
	progress := play.SmoothProgress(p)
	
	// Create marker rectangle points (matching play area bounds)
	markerPoints := []*ui.Point{
		{X: r.playArea.left, Y: r.playArea.top},     // Top left
		{X: r.playArea.right, Y: r.playArea.top},    // Top right  
		{X: r.playArea.right, Y: r.playArea.bottom}, // Bottom right
		{X: r.playArea.left, Y: r.playArea.bottom},  // Bottom left
	}
	
	// Draw marker using center-point perspective (matching play renderer)
	markerPath := vector.Path{}
	cX, cY := r.playArea.centerPoint.ToRender32()
	
	// Start with first point
	startX, startY := markerPoints[0].ToRender32()
	x, y := cX+(startX-cX)*progress, cY+(startY-cY)*progress
	markerPath.MoveTo(x, y)
	
	// Add remaining points
	for i := 1; i < len(markerPoints); i++ {
		px, py := markerPoints[i].ToRender32()
		x, y := cX+(px-cX)*progress, cY+(py-cY)*progress
		markerPath.LineTo(x, y)
	}
	markerPath.Close()
	
	// Line width scales with progress (matching play renderer)
	width := float32(8) * progress
	
	vs, is := markerPath.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{
		Width: width,
	})
	
	ui.ColorVertices(vs, clr)
	r.vectorCollection.Add(vs, is)
}

func (r *EditorRenderer) drawHoldPreview(screen *ebiten.Image, currentTime int64) {
	if currentTime <= r.state.HoldStartTime() {
		return
	}
	
	// Draw preview of hold note being created
	previewColor := color.RGBA{255, 255, 0, 100}
	
	// Draw hold start and end indicators
	r.drawTimeIndicator(screen, r.state.HoldStartTime(), previewColor)
	r.drawTimeIndicator(screen, currentTime, previewColor)
	
	// Could add a connecting line here if needed
}

func (r *EditorRenderer) drawEventMarkers(screen *ebiten.Image, currentTime int64) {
	// Get all events from the editor state
	events := r.state.GetAllEvents()
	if len(events) == 0 {
		return
	}
	
	// Apply lane speed to travel time (exactly matching play renderer)
	baseTravelTime := 10000.0
	travelTime := int64(baseTravelTime / r.state.GetLaneSpeed())
	
	for _, event := range events {
		eventTime := event.GetTime()
		progress := math.Max(0, 1-float64(eventTime-currentTime)/float64(travelTime))
		
		// Only draw events that are visible (within travel range)
		if progress > 0 && progress <= 1.0 {
			var markerColor color.RGBA
			
			switch event.GetType() {
			case "audio_start":
				markerColor = color.RGBA{255, 0, 0, 180} // Red for audio start
			case "countdown_beat":
				markerColor = color.RGBA{255, 255, 0, 150} // Light yellow for countdown beat
			default:
				continue // Skip unknown event types
			}
			
			r.drawEventMarker(screen, progress, markerColor)
		}
	}
}

func (r *EditorRenderer) drawEventMarker(screen *ebiten.Image, p float64, clr color.RGBA) {
	if p < 0 || p > 1 {
		return
	}
	
	// Apply smooth progress (exactly matching play renderer)
	progress := play.SmoothProgress(p)
	
	// Create marker rectangle points (matching play area bounds - exactly like measure markers)
	markerPoints := []*ui.Point{
		{X: r.playArea.left, Y: r.playArea.top},     // Top left
		{X: r.playArea.right, Y: r.playArea.top},    // Top right  
		{X: r.playArea.right, Y: r.playArea.bottom}, // Bottom right
		{X: r.playArea.left, Y: r.playArea.bottom},  // Bottom left
	}
	
	// Draw marker using center-point perspective (matching play renderer)
	markerPath := vector.Path{}
	cX, cY := r.playArea.centerPoint.ToRender32()
	
	// Start with first point
	startX, startY := markerPoints[0].ToRender32()
	x, y := cX+(startX-cX)*progress, cY+(startY-cY)*progress
	markerPath.MoveTo(x, y)
	
	// Add remaining points
	for i := 1; i < len(markerPoints); i++ {
		px, py := markerPoints[i].ToRender32()
		x, y := cX+(px-cX)*progress, cY+(py-cY)*progress
		markerPath.LineTo(x, y)
	}
	markerPath.Close()
	
	// Line width scales with progress (exactly matching measure markers)
	width := float32(8) * progress
	
	vs, is := markerPath.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{
		Width:    width,
		LineCap:  vector.LineCapRound,
		LineJoin: vector.LineJoinRound,
	})
	
	ui.ColorVertices(vs, clr)
	r.vectorCollection.Add(vs, is)
}

// Helper function to calculate distance between two 32-bit float points
func distance32(dx, dy float32) float32 {
	return float32(math.Sqrt(float64(dx*dx + dy*dy)))
}


// getNoteWidth calculates the base note width with proper scaling (matching play renderer)
func (r *EditorRenderer) getNoteWidth() float32 {
	noteWidth := float32(30) // Base note width
	renderWidth, _ := display.Window.RenderSize()
	return noteWidth * (float32(renderWidth) / 1280) // Scale based on render width
}