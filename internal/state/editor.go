package state

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/audio"
	"github.com/liqmix/slaptrax/internal/input"
	"github.com/liqmix/slaptrax/internal/logger"
	"github.com/liqmix/slaptrax/internal/system"
	"github.com/liqmix/slaptrax/internal/types"
	"github.com/liqmix/slaptrax/internal/user"
	"github.com/liqmix/slaptrax/internal/types/schema"
)

type EditorArgs struct {
	Song *types.Song
}

type ChartAction struct {
	Type        string
	TrackName   types.TrackName
	Time        int64
	Note        *types.Note
	PrevNote    *types.Note
}

type EditorChart struct {
	tracks       map[types.TrackName][]*types.Note
	events       []types.Event
	metadata     EditorChartMetadata
	modified     bool
	totalNotes   int
	totalHolds   int
}

type EditorChartMetadata struct {
	Name       string
	Difficulty int
}

type EditorState struct {
	types.BaseGameState
	song           *types.Song
	chart          *EditorChart
	currentTime    int64
	selectedTrack  types.TrackName
	playing        bool
	snapToGrid     bool
	gridSize       int64
	undoStack      []ChartAction
	redoStack      []ChartAction
	startTime      time.Time
	playbackOffset int64
	showGrid       bool
	metronomeOn    bool
	isHolding      bool
	holdStartTime  int64
	bpm            float64
	timeDivision   int // 1, 2, 4, 8, 16, 32, 64, 128 (represents 1/1, 1/2, 1/4, etc.)
	laneSpeed      float64 // Editor lane speed (separate from user settings during editing)
}

func NewEditorState(args *EditorArgs) *EditorState {
	editor := &EditorState{
		song:          args.Song,
		currentTime:   0,
		selectedTrack: types.TrackLeftTop,
		snapToGrid:    true,
		gridSize:      250, // Quarter note at 120 BPM
		showGrid:      true,
		undoStack:     make([]ChartAction, 0),
		redoStack:     make([]ChartAction, 0),
		bpm:           120.0, // Default BPM
		timeDivision:  4,     // Default to quarter notes (1/4)
		laneSpeed:     user.S().LaneSpeed, // Use player's configured lane speed
	}

	// Initialize empty chart or load existing one
	if args.Song != nil {
		editor.initializeChart()
	} else {
		editor.createEmptyChart()
	}

	// Setup keyboard controls
	editor.setupControls()
	
	// Initialize grid size based on BPM and time division
	editor.updateGridSize()
	
	// Stop background music when entering editor
	audio.StopBGM()

	return editor
}

func (e *EditorState) initializeChart() {
	// If song has charts, load the first available one for editing
	if len(e.song.Charts) > 0 {
		var firstChart *types.Chart
		for _, chart := range e.song.Charts {
			firstChart = chart
			break
		}
		
		e.chart = &EditorChart{
			tracks:   make(map[types.TrackName][]*types.Note),
			events:   make([]types.Event, 0),
			modified: false,
		}

		// Copy notes from existing chart
		for _, track := range firstChart.Tracks {
			notes := make([]*types.Note, len(track.AllNotes))
			copy(notes, track.AllNotes)
			e.chart.tracks[track.Name] = notes
		}
	} else {
		e.createEmptyChart()
	}
}

func (e *EditorState) createEmptyChart() {
	e.chart = &EditorChart{
		tracks:   make(map[types.TrackName][]*types.Note),
		events:   make([]types.Event, 0),
		modified: false,
		metadata: EditorChartMetadata{
			Name:       "New Chart",
			Difficulty: 1,
		},
	}

	// Initialize empty note arrays for all tracks
	for _, trackName := range types.TrackNames() {
		e.chart.tracks[trackName] = make([]*types.Note, 0)
	}
}

func (e *EditorState) setupControls() {
	// Only set up essential actions that don't conflict
	e.SetAction(input.ActionBack, e.exitEditor)
	// Note: Using direct key input in Update() for navigation to avoid conflicts
}

func (e *EditorState) moveLeft() {
	step := e.getTimeStep()
	e.currentTime = max(0, e.currentTime-step)
	e.stopPlayback()
}

func (e *EditorState) moveRight() {
	step := e.getTimeStep()
	if e.song != nil {
		maxTime := int64(e.song.Length)
		e.currentTime = min(maxTime, e.currentTime+step)
	} else {
		e.currentTime += step
	}
	e.stopPlayback()
}

func (e *EditorState) selectTrackUp() {
	// More intuitive track selection order: visual left-to-right, top-to-bottom
	visualOrder := []types.TrackName{
		types.TrackLeftTop, types.TrackCenterTop, types.TrackRightTop,
		types.TrackLeftBottom, types.TrackCenterBottom, types.TrackRightBottom,
	}
	
	currentIndex := 0
	for i, track := range visualOrder {
		if track == e.selectedTrack {
			currentIndex = i
			break
		}
	}
	
	// Navigate upwards in the visual grid
	if currentIndex >= 3 {
		newIndex := currentIndex - 3 // Move up one row
		e.selectedTrack = visualOrder[newIndex]
	}
}

func (e *EditorState) selectTrackDown() {
	// More intuitive track selection order: visual left-to-right, top-to-bottom
	visualOrder := []types.TrackName{
		types.TrackLeftTop, types.TrackCenterTop, types.TrackRightTop,
		types.TrackLeftBottom, types.TrackCenterBottom, types.TrackRightBottom,
	}
	
	currentIndex := 0
	for i, track := range visualOrder {
		if track == e.selectedTrack {
			currentIndex = i
			break
		}
	}
	
	// Navigate downwards in the visual grid
	if currentIndex < 3 {
		newIndex := currentIndex + 3 // Move down one row
		e.selectedTrack = visualOrder[newIndex]
	}
}

func (e *EditorState) selectTrackLeft() {
	// Move left in the visual grid (3x2 layout)
	visualOrder := []types.TrackName{
		types.TrackLeftTop, types.TrackCenterTop, types.TrackRightTop,
		types.TrackLeftBottom, types.TrackCenterBottom, types.TrackRightBottom,
	}
	
	currentIndex := 0
	for i, track := range visualOrder {
		if track == e.selectedTrack {
			currentIndex = i
			break
		}
	}
	
	// Move left within the same row
	row := currentIndex / 3 // 0 = top row, 1 = bottom row
	col := currentIndex % 3 // 0 = left, 1 = center, 2 = right
	
	if col > 0 {
		newIndex := row*3 + (col - 1)
		e.selectedTrack = visualOrder[newIndex]
	}
}

func (e *EditorState) selectTrackRight() {
	// Move right in the visual grid (3x2 layout)
	visualOrder := []types.TrackName{
		types.TrackLeftTop, types.TrackCenterTop, types.TrackRightTop,
		types.TrackLeftBottom, types.TrackCenterBottom, types.TrackRightBottom,
	}
	
	currentIndex := 0
	for i, track := range visualOrder {
		if track == e.selectedTrack {
			currentIndex = i
			break
		}
	}
	
	// Move right within the same row
	row := currentIndex / 3 // 0 = top row, 1 = bottom row
	col := currentIndex % 3 // 0 = left, 1 = center, 2 = right
	
	if col < 2 {
		newIndex := row*3 + (col + 1)
		e.selectedTrack = visualOrder[newIndex]
	}
}

func (e *EditorState) exitEditor() {
	if e.chart.modified {
		// TODO: Show save confirmation dialog
	}
	audio.StopAll()
	
	// Resume background music when exiting to title
	audio.PlayBGM(audio.BGMTitle)
	
	e.SetNextState(types.GameStateTitle, nil)
}

func (e *EditorState) getTimeStep() int64 {
	if input.K.Is(ebiten.KeyShift, input.Held) {
		// Measure (4 beats)
		return e.calculateTimeDivision(1) * 4
	} else if input.K.Is(ebiten.KeyControl, input.Held) {
		// Fine adjustment
		return 10
	}
	// Current time division
	return e.calculateTimeDivision(e.timeDivision)
}

// Calculate milliseconds for a time division (e.g., 4 = quarter note, 8 = eighth note)
func (e *EditorState) calculateTimeDivision(division int) int64 {
	// 60000 ms per minute / BPM = ms per quarter note
	quarterNoteMs := 60000.0 / e.bpm
	// For quarter notes (division=4), return quarterNoteMs
	// For eighth notes (division=8), return quarterNoteMs/2
	// For half notes (division=2), return quarterNoteMs*2
	return int64(quarterNoteMs * 4.0 / float64(division))
}

func (e *EditorState) adjustBPM(delta float64) {
	e.bpm = math.Max(30.0, math.Min(300.0, e.bpm+delta)) // Clamp BPM between 30-300
	e.updateGridSize()
}

func (e *EditorState) adjustTimeDivision(up bool) {
	divisions := []int{1, 2, 4, 8, 16, 32, 64, 128}
	currentIndex := 0
	
	for i, div := range divisions {
		if div == e.timeDivision {
			currentIndex = i
			break
		}
	}
	
	oldDivision := e.timeDivision
	
	if up && currentIndex < len(divisions)-1 {
		e.timeDivision = divisions[currentIndex+1]
	} else if !up && currentIndex > 0 {
		e.timeDivision = divisions[currentIndex-1]
	}
	
	// If division changed, snap current time to nearest valid position
	if oldDivision != e.timeDivision {
		e.updateGridSize()
		e.currentTime = e.snapTime(e.currentTime)
	}
}

func (e *EditorState) updateGridSize() {
	e.gridSize = e.calculateTimeDivision(e.timeDivision)
}

func (e *EditorState) GetBPM() float64 {
	return e.bpm
}

func (e *EditorState) GetTimeDivision() int {
	return e.timeDivision
}

func (e *EditorState) GetLaneSpeed() float64 {
	return e.laneSpeed
}

func (e *EditorState) adjustLaneSpeed(delta float64) {
	e.laneSpeed = math.Max(0.5, math.Min(10.0, e.laneSpeed+delta)) // Clamp between 0.5-10.0
}

// Get current time position as a time division fraction (e.g., "3/4", "5/8")
func (e *EditorState) GetCurrentTimePosition() string {
	// Calculate how many divisions have passed
	divisionMs := e.calculateTimeDivision(e.timeDivision)
	if divisionMs == 0 {
		return "0/1"
	}
	
	totalDivisions := e.currentTime / divisionMs
	
	// Calculate measures and beats within measure
	// Assuming 4/4 time signature for now
	divisionsPerMeasure := int64(e.timeDivision)
	measure := totalDivisions / divisionsPerMeasure
	beatInMeasure := totalDivisions % divisionsPerMeasure
	
	if measure == 0 {
		return fmt.Sprintf("%d/%d", beatInMeasure, e.timeDivision)
	} else {
		return fmt.Sprintf("M%d:%d/%d", measure+1, beatInMeasure, e.timeDivision)
	}
}

func (e *EditorState) snapTime(time int64) int64 {
	if !e.snapToGrid {
		return time
	}
	return (time + e.gridSize/2) / e.gridSize * e.gridSize
}

func (e *EditorState) placeNote() {
	time := e.snapTime(e.currentTime)
	
	// Check if note already exists at this position
	existingNote := e.getNoteAt(e.selectedTrack, time)
	if existingNote != nil {
		// Remove existing note
		e.removeNote(existingNote)
		return
	}

	// Create new note
	var note *types.Note
	if e.isHolding {
		// End hold note
		if e.holdStartTime > 0 {
			duration := time - e.holdStartTime
			if duration > 0 {
				note = types.NewNote(e.selectedTrack, e.holdStartTime, e.holdStartTime+duration)
				e.addNote(note)
			}
		}
		e.isHolding = false
		e.holdStartTime = 0
	} else {
		if input.K.Is(ebiten.KeyShift, input.Held) {
			// Start hold note
			e.isHolding = true
			e.holdStartTime = time
		} else {
			// Place tap note
			note = types.NewNote(e.selectedTrack, time, 0)
			e.addNote(note)
		}
	}
}

func (e *EditorState) addNote(note *types.Note) {
	action := ChartAction{
		Type:      "add",
		TrackName: note.TrackName,
		Time:      note.Target,
		Note:      note,
	}
	
	e.chart.tracks[note.TrackName] = append(e.chart.tracks[note.TrackName], note)
	e.chart.totalNotes++
	if note.IsHoldNote() {
		e.chart.totalHolds++
	}
	e.chart.modified = true
	
	e.undoStack = append(e.undoStack, action)
	e.redoStack = e.redoStack[:0] // Clear redo stack
}

func (e *EditorState) removeNote(note *types.Note) {
	action := ChartAction{
		Type:      "remove",
		TrackName: note.TrackName,
		Time:      note.Target,
		PrevNote:  note,
	}

	// Remove from track
	notes := e.chart.tracks[note.TrackName]
	for i, n := range notes {
		if n == note {
			e.chart.tracks[note.TrackName] = append(notes[:i], notes[i+1:]...)
			break
		}
	}
	
	e.chart.totalNotes--
	if note.IsHoldNote() {
		e.chart.totalHolds--
	}
	e.chart.modified = true
	
	e.undoStack = append(e.undoStack, action)
	e.redoStack = e.redoStack[:0]
}

func (e *EditorState) getNoteAt(trackName types.TrackName, time int64) *types.Note {
	for _, note := range e.chart.tracks[trackName] {
		if abs(note.Target-time) < e.gridSize/4 {
			return note
		}
	}
	return nil
}

func (e *EditorState) togglePlayback() {
	if e.playing {
		e.stopPlayback()
	} else {
		e.startPlayback()
	}
}

func (e *EditorState) startPlayback() {
	e.playing = true
	e.startTime = time.Now()
	e.playbackOffset = e.currentTime
	
	// Start audio from current position if song is available
	if e.song != nil {
		audio.InitSong(e.song)
		audio.SetSongPositionMS(int(e.currentTime))
		audio.PlaySong()
	}
	// If no song, playback will just update the timeline position for editing
}

func (e *EditorState) stopPlayback() {
	e.playing = false
	if e.song != nil {
		audio.StopSong()
	}
	// Snap to nearest time division when stopping playback
	e.currentTime = e.snapTime(e.currentTime)
}

func (e *EditorState) undo() {
	if len(e.undoStack) == 0 {
		return
	}
	
	action := e.undoStack[len(e.undoStack)-1]
	e.undoStack = e.undoStack[:len(e.undoStack)-1]
	
	switch action.Type {
	case "add":
		e.removeNoteInternal(action.Note)
		e.redoStack = append(e.redoStack, action)
	case "remove":
		e.addNoteInternal(action.PrevNote)
		e.redoStack = append(e.redoStack, action)
	}
}

func (e *EditorState) redo() {
	if len(e.redoStack) == 0 {
		return
	}
	
	action := e.redoStack[len(e.redoStack)-1]
	e.redoStack = e.redoStack[:len(e.redoStack)-1]
	
	switch action.Type {
	case "add":
		e.addNoteInternal(action.Note)
		e.undoStack = append(e.undoStack, action)
	case "remove":
		e.removeNoteInternal(action.PrevNote)
		e.undoStack = append(e.undoStack, action)
	}
}

func (e *EditorState) addNoteInternal(note *types.Note) {
	e.chart.tracks[note.TrackName] = append(e.chart.tracks[note.TrackName], note)
	e.chart.totalNotes++
	if note.IsHoldNote() {
		e.chart.totalHolds++
	}
	e.chart.modified = true
}

func (e *EditorState) removeNoteInternal(note *types.Note) {
	notes := e.chart.tracks[note.TrackName]
	for i, n := range notes {
		if n == note {
			e.chart.tracks[note.TrackName] = append(notes[:i], notes[i+1:]...)
			break
		}
	}
	e.chart.totalNotes--
	if note.IsHoldNote() {
		e.chart.totalHolds--
	}
	e.chart.modified = true
}

func (e *EditorState) Update() error {
	e.BaseGameState.Update()
	
	// Update playback time
	if e.playing {
		elapsed := time.Since(e.startTime).Milliseconds()
		e.currentTime = e.playbackOffset + elapsed
		
		// Stop playback if we've reached the end of the song
		if e.song != nil && e.currentTime >= int64(e.song.Length) {
			e.stopPlayback()
		}
		// For songs without audio, allow unlimited playback (useful for editing)
	}
	
	// Handle BPM and time division controls first (to avoid conflicts with navigation modifiers)
	bpmDivisionHandled := false
	
	// BPM controls (Shift + -/+)
	if input.K.Is(ebiten.KeyShift, input.Held) {
		if input.K.Is(ebiten.KeyMinus, input.JustPressed) {
			e.adjustBPM(-1.0)
			bpmDivisionHandled = true
		}
		if input.K.Is(ebiten.KeyEqual, input.JustPressed) { // Plus key
			e.adjustBPM(1.0)
			bpmDivisionHandled = true
		}
	}
	
	// Time division controls (Ctrl + -/+)
	if input.K.Is(ebiten.KeyControl, input.Held) {
		if input.K.Is(ebiten.KeyMinus, input.JustPressed) {
			e.adjustTimeDivision(false) // Decrease division (longer notes)
			bpmDivisionHandled = true
		}
		if input.K.Is(ebiten.KeyEqual, input.JustPressed) { // Plus key
			e.adjustTimeDivision(true) // Increase division (shorter notes)
			bpmDivisionHandled = true
		}
	}
	
	// Lane speed controls (Alt + -/+)
	if input.K.Is(ebiten.KeyAlt, input.Held) {
		if input.K.Is(ebiten.KeyMinus, input.JustPressed) {
			e.adjustLaneSpeed(-0.5)
			bpmDivisionHandled = true
		}
		if input.K.Is(ebiten.KeyEqual, input.JustPressed) { // Plus key
			e.adjustLaneSpeed(0.5)
			bpmDivisionHandled = true
		}
	}
	
	// Handle time navigation (arrow keys) - only if BPM/division wasn't handled
	if !bpmDivisionHandled {
		if input.K.Is(ebiten.KeyArrowLeft, input.JustPressed) {
			e.moveLeft()
		}
		if input.K.Is(ebiten.KeyArrowRight, input.JustPressed) {
			e.moveRight()
		}
	}
	
	// Handle track selection (WASD)
	if input.K.Is(ebiten.KeyW, input.JustPressed) {
		e.selectTrackUp()
	}
	if input.K.Is(ebiten.KeyS, input.JustPressed) {
		e.selectTrackDown()
	}
	if input.K.Is(ebiten.KeyA, input.JustPressed) {
		e.selectTrackLeft()
	}
	if input.K.Is(ebiten.KeyD, input.JustPressed) {
		e.selectTrackRight()
	}
	
	// Handle keyboard input
	if input.K.Is(ebiten.KeySpace, input.JustPressed) {
		e.placeNote()
	}
	
	if input.K.Is(ebiten.KeyP, input.JustPressed) {
		e.togglePlayback()
	}
	
	if input.K.Is(ebiten.KeyR, input.JustPressed) {
		e.currentTime = 0
		e.stopPlayback()
	}
	
	if input.K.Is(ebiten.KeyHome, input.JustPressed) {
		e.currentTime = 0
		e.stopPlayback()
	}
	
	if input.K.Is(ebiten.KeyEnd, input.JustPressed) {
		if e.song != nil {
			e.currentTime = int64(e.song.Length)
		}
		e.stopPlayback()
	}
	
	if input.K.Is(ebiten.KeyPageUp, input.JustPressed) {
		// Jump back by 4 measures
		measureTime := e.calculateTimeDivision(1) * 4 
		e.currentTime = max(0, e.currentTime-measureTime*4)
		e.stopPlayback()
	}
	
	if input.K.Is(ebiten.KeyPageDown, input.JustPressed) {
		// Jump forward by 4 measures
		measureTime := e.calculateTimeDivision(1) * 4
		e.currentTime += measureTime * 4
		e.stopPlayback()
	}
	
	if input.K.Is(ebiten.KeyDelete, input.JustPressed) {
		note := e.getNoteAt(e.selectedTrack, e.snapTime(e.currentTime))
		if note != nil {
			e.removeNote(note)
		}
	}
	
	if input.K.Is(ebiten.KeyG, input.JustPressed) {
		e.showGrid = !e.showGrid
	}
	
	if input.K.Is(ebiten.KeyM, input.JustPressed) {
		e.metronomeOn = !e.metronomeOn
	}
	
	// Ctrl key combinations (only undo/redo here, BPM/division handled earlier)  
	if input.K.Is(ebiten.KeyControl, input.Held) && !bpmDivisionHandled {
		// Undo/Redo
		if input.K.Is(ebiten.KeyZ, input.JustPressed) {
			e.undo()
		}
		if input.K.Is(ebiten.KeyY, input.JustPressed) {
			e.redo()
		}
		if input.K.Is(ebiten.KeyS, input.JustPressed) {
			e.saveChart()
		}
		if input.K.Is(ebiten.KeyO, input.JustPressed) {
			if input.K.Is(ebiten.KeyShift, input.Held) {
				e.openAudioFile()
			} else {
				e.loadChart()
			}
		}
		if input.K.Is(ebiten.KeyN, input.JustPressed) {
			// New chart - reset to empty
			e.song = nil
			e.createEmptyChart()
		}
	}
	
	return nil
}

// Getter methods for renderer
func (e *EditorState) GetNotesForTrack(trackName types.TrackName) []*types.Note {
	return e.chart.tracks[trackName]
}

func (e *EditorState) CurrentTime() int64 {
	return e.currentTime
}

func (e *EditorState) SelectedTrack() types.TrackName {
	return e.selectedTrack
}

func (e *EditorState) ShowGrid() bool {
	return e.showGrid
}

func (e *EditorState) GridSize() int64 {
	return e.gridSize
}

func (e *EditorState) Playing() bool {
	return e.playing
}

func (e *EditorState) IsHolding() bool {
	return e.isHolding
}

func (e *EditorState) HoldStartTime() int64 {
	return e.holdStartTime
}

// File operations
func (e *EditorState) openAudioFile() {
	filePath, err := system.OpenAudioFileDialog()
	if err != nil {
		logger.Warn("Failed to open audio file dialog: %v", err)
		return
	}
	
	if filePath == "" {
		return // User cancelled
	}
	
	e.loadExternalSong(filePath)
}

func (e *EditorState) loadExternalSong(audioPath string) {
	// Create a new song from external audio file
	baseName := filepath.Base(audioPath)
	songName := baseName[:len(baseName)-len(filepath.Ext(baseName))]
	
	// TODO: Get actual duration from audio file
	duration := int64(180000) // 3 minutes default
	
	e.song = &types.Song{
		Title:     songName,
		Artist:    "Unknown Artist",
		AudioPath: audioPath,
		BPM:       120, // Default BPM
		Length:    int(duration),
		Charts:    make(map[types.Difficulty]*types.Chart),
	}
	
	e.createEmptyChart()
	logger.Debug("Loaded external song: %s", songName)
}

func (e *EditorState) saveChart() {
	if e.song == nil {
		logger.Warn("No song loaded to save")
		return
	}
	
	defaultName := fmt.Sprintf("%s_%s.json", e.song.Title, e.chart.metadata.Name)
	filePath, err := system.SaveJSONFileDialog(defaultName)
	if err != nil {
		logger.Warn("Failed to show save dialog: %v", err)
		return
	}
	
	if filePath == "" {
		return // User cancelled
	}
	
	// Export to JSON format
	chartData := e.exportToJSON()
	data, err := json.MarshalIndent(chartData, "", "  ")
	if err != nil {
		logger.Error("Failed to marshal chart data: %v", err)
		return
	}
	
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		logger.Error("Failed to save chart file: %v", err)
		return
	}
	
	e.chart.modified = false
	logger.Debug("Chart saved to: %s", filePath)
}

func (e *EditorState) loadChart() {
	filePath, err := system.OpenJSONFileDialog()
	if err != nil {
		logger.Warn("Failed to open chart dialog: %v", err)
		return
	}
	
	if filePath == "" {
		return // User cancelled
	}
	
	data, err := os.ReadFile(filePath)
	if err != nil {
		logger.Error("Failed to read chart file: %v", err)
		return
	}
	
	var songData schema.SongDataV2
	err = json.Unmarshal(data, &songData)
	if err != nil {
		logger.Error("Failed to parse chart file: %v", err)
		return
	}
	
	// Load the first chart for editing
	if len(songData.Charts) == 0 {
		logger.Warn("No charts found in file")
		return
	}
	
	var firstChartData schema.ChartDataV2
	for _, chartData := range songData.Charts {
		firstChartData = chartData
		break
	}
	
	e.importFromJSON(&songData, &firstChartData)
	logger.Debug("Chart loaded from: %s", filePath)
}

func (e *EditorState) exportToJSON() *schema.SongDataV2 {
	// Create chart data
	chartData := schema.ChartDataV2{
		Name:       e.chart.metadata.Name,
		Difficulty: e.chart.metadata.Difficulty,
		NoteCount:  e.chart.totalNotes,
		HoldCount:  e.chart.totalHolds,
		MaxCombo:   e.calculateMaxCombo(),
		Tracks:     make(map[string][]schema.NoteData),
		Events:     []schema.EventData{},
	}
	
	// Convert notes to JSON format
	for trackName, notes := range e.chart.tracks {
		if len(notes) == 0 {
			continue
		}
		
		trackStr := trackNameToString(trackName)
		jsonNotes := make([]schema.NoteData, len(notes))
		
		// Sort notes by time
		sortedNotes := make([]*types.Note, len(notes))
		copy(sortedNotes, notes)
		sort.Slice(sortedNotes, func(i, j int) bool {
			return sortedNotes[i].Target < sortedNotes[j].Target
		})
		
		for i, note := range sortedNotes {
			jsonNotes[i] = schema.NoteData{
				Time: note.Target,
				Type: schema.NoteTypeTap,
			}
			
			if note.IsHoldNote() {
				jsonNotes[i].Type = schema.NoteTypeHold
				jsonNotes[i].Duration = note.TargetRelease - note.Target
			}
		}
		
		chartData.Tracks[trackStr] = jsonNotes
	}
	
	// Create song data
	songData := &schema.SongDataV2{
		Schema:  "https://slaptrax.dev/schema/song/v2.json",
		Version: 2,
		Metadata: schema.SongMetadata{
			Title:           e.song.Title,
			Artist:          e.song.Artist,
			BPM:             e.song.BPM,
			PreviewStart:    30000, // Default preview start
			Duration:      int64(e.song.Length),
			ChartedBy:       "Chart Editor",
			Version:         "1.0.0",
			DifficultyRange: [2]int{e.chart.metadata.Difficulty, e.chart.metadata.Difficulty},
		},
		Audio: schema.AudioInfo{
			File: filepath.Base(e.song.AudioPath),
		},
		Visual: schema.VisualInfo{
			Theme: "default",
		},
		Charts: map[string]schema.ChartDataV2{
			fmt.Sprintf("%d", e.chart.metadata.Difficulty): chartData,
		},
	}
	
	return songData
}

func (e *EditorState) importFromJSON(songData *schema.SongDataV2, chartData *schema.ChartDataV2) {
	// Create song from metadata
	e.song = &types.Song{
		Title:     songData.Metadata.Title,
		Artist:    songData.Metadata.Artist,
		BPM:       songData.Metadata.BPM,
		Length:  int(songData.Metadata.Duration),
		AudioPath: songData.Audio.File, // This might need full path resolution
		Charts:    make(map[types.Difficulty]*types.Chart),
	}
	
	// Create chart
	e.chart = &EditorChart{
		tracks:   make(map[types.TrackName][]*types.Note),
		events:   make([]types.Event, 0),
		modified: false,
		metadata: EditorChartMetadata{
			Name:       chartData.Name,
			Difficulty: chartData.Difficulty,
		},
		totalNotes: chartData.NoteCount,
		totalHolds: chartData.HoldCount,
	}
	
	// Initialize empty tracks
	for _, trackName := range types.TrackNames() {
		e.chart.tracks[trackName] = make([]*types.Note, 0)
	}
	
	// Import notes
	for trackStr, notes := range chartData.Tracks {
		trackName := stringToTrackName(trackStr)
		if trackName == types.TrackUnknown {
			continue
		}
		
		for _, noteData := range notes {
			var note *types.Note
			if noteData.Type == schema.NoteTypeHold && noteData.Duration > 0 {
				note = types.NewNote(trackName, noteData.Time, noteData.Time+noteData.Duration)
			} else {
				note = types.NewNote(trackName, noteData.Time, 0)
			}
			e.chart.tracks[trackName] = append(e.chart.tracks[trackName], note)
		}
	}
}

func (e *EditorState) calculateMaxCombo() int {
	combo := 0
	for _, notes := range e.chart.tracks {
		combo += len(notes)
	}
	return combo
}

func trackNameToString(trackName types.TrackName) string {
	switch trackName {
	case types.TrackLeftTop:
		return "left_top"
	case types.TrackLeftBottom:
		return "left_bottom"
	case types.TrackCenterTop:
		return "center_top"
	case types.TrackCenterBottom:
		return "center_bottom"
	case types.TrackRightTop:
		return "right_top"
	case types.TrackRightBottom:
		return "right_bottom"
	default:
		return "unknown"
	}
}

func stringToTrackName(trackStr string) types.TrackName {
	switch trackStr {
	case "left_top":
		return types.TrackLeftTop
	case "left_bottom":
		return types.TrackLeftBottom
	case "center_top":
		return types.TrackCenterTop
	case "center_bottom":
		return types.TrackCenterBottom
	case "right_top":
		return types.TrackRightTop
	case "right_bottom":
		return types.TrackRightBottom
	default:
		return types.TrackUnknown
	}
}

func (e *EditorState) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {}

// Helper functions
func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}