package state

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/audio"
	"github.com/liqmix/slaptrax/internal/input"
	"github.com/liqmix/slaptrax/internal/logger"
	"github.com/liqmix/slaptrax/internal/system"
	"github.com/liqmix/slaptrax/internal/types"
	"github.com/liqmix/slaptrax/internal/types/schema"
	"github.com/liqmix/slaptrax/internal/ui"
)

type EditorArgs struct {
	Song      *types.Song // For editing existing bundled songs
	ChartPath string      // Path to existing chart file to load
	AudioPath string      // Path to audio file for new charts
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
	chartPath      string // Store original chart path for re-entry
	audioPath      string // Store original audio path for re-entry
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
	audioOffset    int64   // Audio starting position offset in milliseconds
	eventManager   *types.EventManager // Event system for playback execution
	audioInitialized bool // Track if audio has been initialized to prevent duplicate initialization
	
	// Rapid offset adjustment tracking
	offsetIncreaseHeldStart time.Time // When Alt+Plus started being held
	offsetDecreaseHeldStart time.Time // When Alt+Minus started being held
	offsetRapidAdjustment   bool      // Whether rapid adjustment is active
	
}

func NewEditorState(args *EditorArgs) *EditorState {
	editor := &EditorState{
		song:          args.Song,
		chartPath:     args.ChartPath,
		audioPath:     args.AudioPath,
		currentTime:   0,
		selectedTrack: types.TrackLeftTop,
		snapToGrid:    true,
		gridSize:      250, // Quarter note at 120 BPM
		showGrid:      true,
		undoStack:     make([]ChartAction, 0),
		redoStack:     make([]ChartAction, 0),
		bpm:           120.0, // Default BPM
		timeDivision:  4,     // Default to quarter notes (1/4)
		laneSpeed:     1.0, // Default editor speed
		audioOffset:   0,   // No audio offset by default
		eventManager:  types.NewEventManager(), // Initialize event system
		audioInitialized: false, // Audio not initialized at start
	}

	// Initialize chart based on provided arguments
	if args.Song != nil {
		// Editing an existing bundled song
		editor.initializeChart()
	} else if args.ChartPath != "" {
		// Loading an existing chart folder
		editor.loadChartFromFolder(args.ChartPath)
	} else if args.AudioPath != "" {
		// Creating new chart with provided audio - for now just create empty chart
		// TODO: Implement audio file loading
		editor.createEmptyChart()
	} else {
		// No specific arguments, create empty chart
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
		
		// Initialize beat positions for existing content
		e.initializeBeatPositions()
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
	
	// Auto-generate event markers for new songs
	e.generateDefaultEvents()
}

func (e *EditorState) generateDefaultEvents() {
	// Use exact beat-based positioning (measure = 4 beats)
	// Add countdown beat events on each quarter beat in the second measure (after one measure of silence)
	for beat := 0; beat < 4; beat++ {
		// Use exact beat fractions: 4.0, 5.0, 6.0, 7.0
		beatPosition := 4.0 + float64(beat) // Second measure starts at beat 4
		countdownEvent := types.NewCountdownBeatEventFromBeat(beatPosition, e.bpm)
		e.chart.events = append(e.chart.events, countdownEvent)
	}
	
	// Add audio start marker at the beginning of the third measure (exact beat 8.0)
	audioStartEvent := types.NewAudioStartEventFromBeat(8.0, e.bpm)
	e.chart.events = append(e.chart.events, audioStartEvent)
	
	// Keep events sorted by time
	sort.Slice(e.chart.events, func(i, j int) bool {
		return e.chart.events[i].GetTime() < e.chart.events[j].GetTime()
	})
	
	// Sync events to event manager
	e.syncEventsToManager()
}

func (e *EditorState) syncEventsToManager() {
	e.eventManager.Clear()
	for _, event := range e.chart.events {
		e.eventManager.AddEvent(event)
	}
	fmt.Printf("📝 Synced %d events to EventManager\n", len(e.chart.events))
}

func (e *EditorState) setupControls() {
	// Only set up essential actions that don't conflict
	e.SetAction(input.ActionBack, e.handleBackAction)
	// Note: Using direct key input in Update() for navigation to avoid conflicts
}

func (e *EditorState) handleBackAction() {
	e.exitEditor()
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

func (e *EditorState) goToBeginning() {
	e.currentTime = 0
	e.stopPlayback()
}

func (e *EditorState) goToLastNote() {
	var lastTime int64 = 0
	
	// Find the latest note across all tracks
	for _, notes := range e.chart.tracks {
		for _, note := range notes {
			var noteEndTime int64
			if note.IsHoldNote() {
				noteEndTime = note.TargetRelease
			} else {
				noteEndTime = note.Target
			}
			if noteEndTime > lastTime {
				lastTime = noteEndTime
			}
		}
	}
	
	// If no notes found, stay at current position
	if lastTime > 0 {
		e.currentTime = lastTime
		e.stopPlayback()
	}
}

func (e *EditorState) exitEditor() {
	// Create modal state for exit confirmation
	message := "Are you sure you want to exit the editor?"
	if e.chart.modified {
		message = "You have unsaved changes. Exit without saving?"
	}
	
	// Create the modal once
	var modal *ui.ConfirmModal
	
	modalArgs := &ModalStateArgs{
		Update: func(setNext func(types.GameState, interface{})) {
			// Create modal if not already created
			if modal == nil {
				modal = ui.NewConfirmModal(ui.ConfirmModalOptions{
					Title:   "Exit Editor",
					Message: message,
					YesText: "Exit",
					NoText:  "Stay",
					OnYes: func() {
						e.confirmExit(setNext)
					},
					OnNo: func() {
						e.cancelExit(setNext)
					},
					OnCancel: func() {
						e.cancelExit(setNext)
					},
				})
			}
			modal.Update()
		},
		Draw: func(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
			if modal != nil {
				modal.Draw(screen, opts)
			}
		},
	}
	
	// Transition to modal state
	e.SetNextState(types.GameStateModal, modalArgs)
}

func (e *EditorState) confirmExit(setNext func(types.GameState, interface{})) {
	audio.StopAll()
	
	// Resume background music when exiting to title
	audio.PlayBGM(audio.BGMTitle)
	
	setNext(types.GameStateTitle, nil)
}

func (e *EditorState) cancelExit(setNext func(types.GameState, interface{})) {
	// Return to the previous state (editor) without recreating it
	setNext(types.GameStateBack, nil)
}

func (e *EditorState) getTimeStep() int64 {
	if e.isShiftHeld() {
		// Measure (4 beats)
		return e.calculateTimeDivision(1) * 4
	} else if e.isControlHeld() {
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
	
	// Keep song BPM synchronized with editor BPM
	if e.song != nil {
		e.song.BPM = int(e.bpm)
	}
	
	e.updateGridSize()
	e.recalculateTimestamps() // Update all notes and events to new BPM
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

func (e *EditorState) GetWholeNoteMs() int64 {
	return e.calculateTimeDivision(1) // Always return whole note duration (1/1)
}

func (e *EditorState) GetQuarterNoteMs() int64 {
	return e.calculateTimeDivision(4) // Always return quarter note duration (1/4)
}

func (e *EditorState) GetEighthNoteMs() int64 {
	return e.calculateTimeDivision(8) // Always return eighth note duration (1/8)
}

func (e *EditorState) adjustLaneSpeed(delta float64) {
	e.laneSpeed = math.Max(0.5, math.Min(10.0, e.laneSpeed+delta)) // Clamp between 0.5-10.0
}

func (e *EditorState) GetAudioOffset() int64 {
	return e.audioOffset
}

// Public beat conversion methods for renderer access
func (e *EditorState) MsIntToBeats(ms int64) float64 {
	return e.msIntToBeats(ms)
}

func (e *EditorState) BeatsToMsInt(beats float64) int64 {
	return e.beatsToMsInt(beats)
}

// Beat conversion utilities (high precision)
func (e *EditorState) beatsToMs(beats float64) float64 {
	quarterNoteMs := 60000.0 / e.bpm
	return beats * quarterNoteMs
}

func (e *EditorState) msToBeats(ms float64) float64 {
	quarterNoteMs := 60000.0 / e.bpm
	return ms / quarterNoteMs
}

// Legacy integer conversion for compatibility
func (e *EditorState) beatsToMsInt(beats float64) int64 {
	return int64(math.Round(e.beatsToMs(beats)))
}

func (e *EditorState) msIntToBeats(ms int64) float64 {
	return e.msToBeats(float64(ms))
}

// initializeBeatPositions calculates beat positions for existing content using the song's original BPM
func (e *EditorState) initializeBeatPositions() {
	// Use the song's BPM if available, otherwise use current editor BPM
	originalBPM := e.bpm
	if e.song != nil && e.song.BPM > 0 {
		originalBPM = float64(e.song.BPM)
	}
	
	quarterNoteMs := 60000.0 / originalBPM
	
	// Initialize beat positions for notes that don't have them
	for _, notes := range e.chart.tracks {
		for _, note := range notes {
			if note.TargetBeat == 0 && note.Target > 0 {
				note.TargetBeat = float64(note.Target) / quarterNoteMs
			}
			if note.TargetReleaseBeat == 0 && note.TargetRelease > 0 {
				note.TargetReleaseBeat = float64(note.TargetRelease) / quarterNoteMs
			}
		}
	}
	
	// Initialize beat positions for events that don't have them
	for _, event := range e.chart.events {
		switch ev := event.(type) {
		case *types.AudioStartEvent:
			if ev.BaseEvent.TimeBeat == 0 && ev.BaseEvent.Time > 0 {
				ev.BaseEvent.TimeBeat = float64(ev.BaseEvent.Time) / quarterNoteMs
			}
		case *types.CountdownBeatEvent:
			if ev.BaseEvent.TimeBeat == 0 && ev.BaseEvent.Time > 0 {
				ev.BaseEvent.TimeBeat = float64(ev.BaseEvent.Time) / quarterNoteMs
			}
		}
	}
}

// recalculateTimestamps updates all note and event timestamps when BPM changes
func (e *EditorState) recalculateTimestamps() {
	// First, ensure all content has beat positions initialized
	e.initializeBeatPositions()
	
	// Now recalculate all millisecond positions from beat positions using current BPM
	for _, notes := range e.chart.tracks {
		for _, note := range notes {
			if note.TargetBeat > 0 {
				note.Target = e.beatsToMsInt(note.TargetBeat)
			}
			if note.TargetReleaseBeat > 0 {
				note.TargetRelease = e.beatsToMsInt(note.TargetReleaseBeat)
			}
		}
	}
	
	// Update event timestamps from their beat positions  
	for _, event := range e.chart.events {
		switch ev := event.(type) {
		case *types.AudioStartEvent:
			ev.BaseEvent.RecalculateTime(e.bpm)
		case *types.CountdownBeatEvent:
			ev.BaseEvent.RecalculateTime(e.bpm)
		}
	}
	
	// Re-sync events to event manager with updated timestamps
	e.syncEventsToManager()
}

// createNoteAtBeat creates a note at a specific beat position
func (e *EditorState) createNoteAtBeat(trackName types.TrackName, beatPosition, beatRelease float64) *types.Note {
	note := types.NewNoteFromBeats(trackName, beatPosition, beatRelease, e.bpm)
	return note
}

// createNoteAtTime creates a note at a millisecond position and calculates beat position
func (e *EditorState) createNoteAtTime(trackName types.TrackName, timeMs, releaseMs int64) *types.Note {
	beatPosition := e.msIntToBeats(timeMs)
	beatRelease := float64(0)
	if releaseMs > 0 {
		beatRelease = e.msIntToBeats(releaseMs)
	}
	
	// Snap beat positions to ensure they're exact musical fractions
	beatPosition = e.snapToBeat(beatPosition)
	if beatRelease > 0 {
		beatRelease = e.snapToBeat(beatRelease)
	}
	
	note := types.NewNoteFromBeats(trackName, beatPosition, beatRelease, e.bpm)
	return note
}

func (e *EditorState) adjustAudioOffset(delta int64) {
	// Allow audio offset range from -5 seconds to +5 seconds
	e.audioOffset = max(-5000, min(5000, e.audioOffset+delta))
}


// Helper functions to check for modifier keys (including left/right variants)
func (e *EditorState) isAltHeld() bool {
	return input.K.Is(ebiten.KeyAlt, input.Held) ||
		   input.K.Is(ebiten.KeyAltLeft, input.Held) ||
		   input.K.Is(ebiten.KeyAltRight, input.Held)
}

func (e *EditorState) isShiftHeld() bool {
	return input.K.Is(ebiten.KeyShift, input.Held) ||
		   input.K.Is(ebiten.KeyShiftLeft, input.Held) ||
		   input.K.Is(ebiten.KeyShiftRight, input.Held)
}

func (e *EditorState) isControlHeld() bool {
	return input.K.Is(ebiten.KeyControl, input.Held) ||
		   input.K.Is(ebiten.KeyControlLeft, input.Held) ||
		   input.K.Is(ebiten.KeyControlRight, input.Held)
}

func (e *EditorState) handleRapidOffsetAdjustment() {
	const throttleDelayMs = 1000 // 1 second delay before rapid adjustment starts
	const rapidAdjustmentMs = 50 // Adjust every 50ms during rapid mode
	
	now := time.Now()
	
	// Handle increase (Plus/Equal key)
	if input.K.Is(ebiten.KeyEqual, input.Held) {
		if input.K.Is(ebiten.KeyEqual, input.JustPressed) {
			// Just pressed - do immediate adjustment and start tracking
			e.adjustAudioOffset(10)
			e.offsetIncreaseHeldStart = now
			e.offsetRapidAdjustment = false
		} else {
			// Being held - check if we should start rapid adjustment
			heldDuration := now.Sub(e.offsetIncreaseHeldStart)
			if heldDuration >= time.Duration(throttleDelayMs)*time.Millisecond && !e.offsetRapidAdjustment {
				e.offsetRapidAdjustment = true
			}
			
			// If in rapid mode, adjust every rapidAdjustmentMs
			if e.offsetRapidAdjustment && heldDuration.Milliseconds()%rapidAdjustmentMs < 16 { // ~60fps tolerance
				e.adjustAudioOffset(10)
			}
		}
	} else {
		// Reset increase tracking when key is released
		e.offsetIncreaseHeldStart = time.Time{}
		if !input.K.Is(ebiten.KeyMinus, input.Held) {
			e.offsetRapidAdjustment = false
		}
	}
	
	// Handle decrease (Minus key)
	if input.K.Is(ebiten.KeyMinus, input.Held) {
		if input.K.Is(ebiten.KeyMinus, input.JustPressed) {
			// Just pressed - do immediate adjustment and start tracking
			e.adjustAudioOffset(-10)
			e.offsetDecreaseHeldStart = now
			e.offsetRapidAdjustment = false
		} else {
			// Being held - check if we should start rapid adjustment
			heldDuration := now.Sub(e.offsetDecreaseHeldStart)
			if heldDuration >= time.Duration(throttleDelayMs)*time.Millisecond && !e.offsetRapidAdjustment {
				e.offsetRapidAdjustment = true
			}
			
			// If in rapid mode, adjust every rapidAdjustmentMs
			if e.offsetRapidAdjustment && heldDuration.Milliseconds()%rapidAdjustmentMs < 16 { // ~60fps tolerance
				e.adjustAudioOffset(-10)
			}
		}
	} else {
		// Reset decrease tracking when key is released
		e.offsetDecreaseHeldStart = time.Time{}
		if !input.K.Is(ebiten.KeyEqual, input.Held) {
			e.offsetRapidAdjustment = false
		}
	}
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

// snapToBeat snaps a beat position to the nearest grid division
func (e *EditorState) snapToBeat(beats float64) float64 {
	if !e.snapToGrid {
		return beats
	}
	
	// Calculate divisions per beat (e.g., 4 = quarter notes, 16 = sixteenth notes)
	divisionsPerBeat := float64(e.timeDivision) / 4.0
	
	// Snap to nearest division
	snappedBeats := math.Round(beats * divisionsPerBeat) / divisionsPerBeat
	return snappedBeats
}

// snapTime snaps milliseconds to beat grid (legacy function for compatibility)
func (e *EditorState) snapTime(time int64) int64 {
	if !e.snapToGrid {
		return time
	}
	
	// Convert to beats, snap, then convert back
	beats := e.msIntToBeats(time)
	snappedBeats := e.snapToBeat(beats)
	return e.beatsToMsInt(snappedBeats)
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
				note = e.createNoteAtTime(e.selectedTrack, e.holdStartTime, e.holdStartTime+duration)
				e.addNote(note)
			}
		}
		e.isHolding = false
		e.holdStartTime = 0
	} else {
		if e.isShiftHeld() {
			// Start hold note
			e.isHolding = true
			e.holdStartTime = time
		} else {
			// Place tap note
			note = e.createNoteAtTime(e.selectedTrack, time, 0)
			e.addNote(note)
		}
	}
}

func (e *EditorState) toggleNoteOnTrack(trackName types.TrackName) {
	time := e.snapTime(e.currentTime)
	
	// Check if note already exists at this position
	existingNote := e.getNoteAt(trackName, time)
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
				note = e.createNoteAtTime(trackName, e.holdStartTime, e.holdStartTime+duration)
				e.addNote(note)
			}
		}
		e.isHolding = false
		e.holdStartTime = 0
	} else {
		if e.isShiftHeld() {
			// Start hold note
			e.isHolding = true
			e.holdStartTime = time
		} else {
			// Place tap note
			note = e.createNoteAtTime(trackName, time, 0)
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
	
	// Reset event manager for fresh playback
	e.eventManager.Reset()
	
	// Initialize audio if song is available, but don't start playing yet
	// The EventManager will handle playing the audio when it hits the audio start event
	e.safeInitSong(e.song)
	// If no song, playback will just update the timeline position for editing
}

func (e *EditorState) stopPlayback() {
	e.playing = false
	
	// Stop all audio to prevent overlapping
	audio.StopAll()
	e.audioInitialized = false
	
	// Snap to nearest time division when stopping playback
	e.currentTime = e.snapTime(e.currentTime)
}

// safeInitSong prevents multiple audio initializations
func (e *EditorState) safeInitSong(song *types.Song) {
	if !e.audioInitialized && song != nil {
		audio.InitSong(song)
		e.audioInitialized = true
		logger.Debug("Audio initialized for song: %s", song.Title)
	}
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

// Event management methods
func (e *EditorState) addEvent(event types.Event) {
	e.chart.events = append(e.chart.events, event)
	
	// Keep events sorted by time
	sort.Slice(e.chart.events, func(i, j int) bool {
		return e.chart.events[i].GetTime() < e.chart.events[j].GetTime()
	})
	
	e.chart.modified = true
	e.syncEventsToManager()
}

func (e *EditorState) removeEvent(time int64, eventType string) bool {
	for i, event := range e.chart.events {
		// Check if event is at the same time and of the same type
		if abs(event.GetTime()-time) < e.gridSize/4 && event.GetType() == eventType {
			// Remove the event
			e.chart.events = append(e.chart.events[:i], e.chart.events[i+1:]...)
			e.chart.modified = true
			e.syncEventsToManager()
			return true
		}
	}
	return false
}

func (e *EditorState) getEventsAt(time int64) []types.Event {
	var events []types.Event
	for _, event := range e.chart.events {
		if abs(event.GetTime()-time) < e.gridSize/4 {
			events = append(events, event)
		}
	}
	return events
}

func (e *EditorState) getAllEvents() []types.Event {
	// Return a copy to avoid external modification
	events := make([]types.Event, len(e.chart.events))
	copy(events, e.chart.events)
	return events
}

// GetAllEvents returns all events (public method for renderer)
func (e *EditorState) GetAllEvents() []types.Event {
	return e.getAllEvents()
}

// Event marker toggle methods
func (e *EditorState) toggleAudioStartMarker() {
	time := e.snapTime(e.currentTime)
	
	// Check if an audio start event already exists at this position
	if e.removeEvent(time, "audio_start") {
		return // Event was removed
	}
	
	// Remove any existing audio start markers (only one allowed)
	e.removeAllAudioStartMarkers()
	
	// Create new audio start event using beat position
	beatPosition := e.msIntToBeats(time)
	beatPosition = e.snapToBeat(beatPosition) // Ensure exact beat positioning
	event := types.NewAudioStartEventFromBeat(beatPosition, e.bpm)
	e.addEvent(event)
}

func (e *EditorState) removeAllAudioStartMarkers() {
	// Remove all existing audio start events
	var filteredEvents []types.Event
	for _, event := range e.chart.events {
		if event.GetType() != "audio_start" {
			filteredEvents = append(filteredEvents, event)
		}
	}
	
	if len(filteredEvents) != len(e.chart.events) {
		e.chart.events = filteredEvents
		e.chart.modified = true
		e.syncEventsToManager()
	}
}

func (e *EditorState) toggleCountdownBeatMarker() {
	time := e.snapTime(e.currentTime)
	
	// Check if a countdown beat event already exists at this position
	if e.removeEvent(time, "countdown_beat") {
		return // Event was removed
	}
	
	// Create new countdown beat event using beat position
	beatPosition := e.msIntToBeats(time)
	beatPosition = e.snapToBeat(beatPosition) // Ensure exact beat positioning
	event := types.NewCountdownBeatEventFromBeat(beatPosition, e.bpm)
	e.addEvent(event)
}

func (e *EditorState) Update() error {
	e.BaseGameState.Update()
	
	// Update playback time
	if e.playing {
		elapsed := time.Since(e.startTime).Milliseconds()
		e.currentTime = e.playbackOffset + elapsed
		
		// Execute events during playback
		ctx := &types.EventContext{
			CurrentTime:          e.currentTime,
			Song:                 e.song,
			Chart:                nil, // TODO: Convert EditorChart to Chart if needed
			AudioOffset:          e.audioOffset,
			AudioInitSong:        e.safeInitSong,
			AudioPlaySong:        audio.PlaySong,
			AudioSetSongPosition: audio.SetSongPositionMS,
			AudioPlaySFX:         func(sfxCode string) {
				switch sfxCode {
				case "hat":
					audio.PlaySFX(audio.SFXHat)
				case "select":
					audio.PlaySFX(audio.SFXSelect)
				default:
					// Default to select sound for unknown codes
					audio.PlaySFX(audio.SFXSelect)
				}
			},
		}
		if err := e.eventManager.Update(e.currentTime, ctx); err != nil {
			logger.Warn("Event execution error: %v", err)
		}
		
		// Stop playback if we've reached the end of the song
		if e.song != nil && e.currentTime >= int64(e.song.Length) {
			e.stopPlayback()
		}
		// For songs without audio, allow unlimited playback (useful for editing)
	}
	
	// Handle BPM and time division controls first (to avoid conflicts with navigation modifiers)
	bpmDivisionHandled := false
	
	// BPM controls (Shift + -/+)
	if e.isShiftHeld() {
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
	if e.isControlHeld() {
		if input.K.Is(ebiten.KeyMinus, input.JustPressed) {
			e.adjustTimeDivision(false) // Decrease division (longer notes)
			bpmDivisionHandled = true
		}
		if input.K.Is(ebiten.KeyEqual, input.JustPressed) { // Plus key
			e.adjustTimeDivision(true) // Increase division (shorter notes)
			bpmDivisionHandled = true
		}
	}
	
	// Handle arrow keys - context depends on modifiers
	if !bpmDivisionHandled {
		// Lane speed controls (Alt + Up/Down arrows) and audio offset (Alt + Plus/Minus)
		if e.isAltHeld() {
			if input.K.Is(ebiten.KeyArrowUp, input.JustPressed) {
				e.adjustLaneSpeed(0.5) // Increase lane speed
				bpmDivisionHandled = true
			}
			if input.K.Is(ebiten.KeyArrowDown, input.JustPressed) {
				e.adjustLaneSpeed(-0.5) // Decrease lane speed
				bpmDivisionHandled = true
			}
			
			// Handle audio offset adjustment with throttled rapid adjustment
			e.handleRapidOffsetAdjustment()
			bpmDivisionHandled = true
		} else {
			// Time navigation (Left/Right arrows) and time division (Up/Down arrows)
			if input.K.Is(ebiten.KeyArrowLeft, input.JustPressed) {
				e.moveLeft()
			}
			if input.K.Is(ebiten.KeyArrowRight, input.JustPressed) {
				e.moveRight()
			}
			if input.K.Is(ebiten.KeyArrowUp, input.JustPressed) {
				e.adjustTimeDivision(true) // Increase division (shorter notes)
			}
			if input.K.Is(ebiten.KeyArrowDown, input.JustPressed) {
				e.adjustTimeDivision(false) // Decrease division (longer notes)
			}
		}
	}
	
	// Handle direct track note toggling (QWE for top tracks, ASD for bottom tracks)
	if input.K.Is(ebiten.KeyQ, input.JustPressed) {
		e.toggleNoteOnTrack(types.TrackLeftTop)
	}
	if input.K.Is(ebiten.KeyW, input.JustPressed) {
		e.toggleNoteOnTrack(types.TrackCenterTop)
	}
	if input.K.Is(ebiten.KeyE, input.JustPressed) {
		e.toggleNoteOnTrack(types.TrackRightTop)
	}
	if input.K.Is(ebiten.KeyA, input.JustPressed) {
		e.toggleNoteOnTrack(types.TrackLeftBottom)
	}
	if input.K.Is(ebiten.KeyS, input.JustPressed) {
		e.toggleNoteOnTrack(types.TrackCenterBottom)
	}
	if input.K.Is(ebiten.KeyD, input.JustPressed) {
		e.toggleNoteOnTrack(types.TrackRightBottom)
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
	
	if input.K.Is(ebiten.KeyB, input.JustPressed) {
		if e.isShiftHeld() {
			e.goToLastNote()
		} else {
			e.goToBeginning()
		}
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
	
	// Event markers
	if input.K.Is(ebiten.KeyK, input.JustPressed) {
		e.toggleAudioStartMarker()
	}
	
	if input.K.Is(ebiten.KeyI, input.JustPressed) {
		e.toggleCountdownBeatMarker()
	}
	
	// Ctrl key combinations (only undo/redo here, BPM/division handled earlier)  
	if e.isControlHeld() && !bpmDivisionHandled {
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
			if e.isShiftHeld() {
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

func (e *EditorState) IsPlaying() bool {
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
	// Stop any existing playback first
	e.stopPlayback()
	
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
	
	// Reset audio initialization flag for new song
	e.audioInitialized = false
	
	e.createEmptyChart()
	logger.Debug("Loaded external song: %s", songName)
}

func (e *EditorState) saveChart() {
	if e.song == nil {
		logger.Warn("No song loaded to save")
		return
	}
	
	// Select folder to save the chart in
	folderPath, err := system.SelectFolderDialog("Select folder to save chart")
	if err != nil {
		logger.Warn("Failed to show folder dialog: %v", err)
		return
	}
	
	if folderPath == "" {
		return // User cancelled
	}
	
	// Create song folder name based on title
	songFolderName := strings.ReplaceAll(e.song.Title, " ", "-")
	songFolderName = strings.ToLower(songFolderName)
	songFolderPath := filepath.Join(folderPath, songFolderName)
	
	// Create the song folder
	err = os.MkdirAll(songFolderPath, 0755)
	if err != nil {
		logger.Error("Failed to create song folder: %v", err)
		return
	}
	
	// Export to JSON format and save as song.json
	chartData := e.exportToJSON()
	data, err := json.MarshalIndent(chartData, "", "  ")
	if err != nil {
		logger.Error("Failed to marshal chart data: %v", err)
		return
	}
	
	songJSONPath := filepath.Join(songFolderPath, "song.json")
	err = os.WriteFile(songJSONPath, data, 0644)
	if err != nil {
		logger.Error("Failed to save chart file: %v", err)
		return
	}
	
	// Copy audio file if available
	if e.song.AudioPath != "" {
		srcAudioPath := e.song.AudioPath
		destAudioPath := filepath.Join(songFolderPath, "audio.ogg")
		
		// Check if source exists
		if _, err := os.Stat(srcAudioPath); err == nil {
			// Copy the audio file
			err = e.copyFile(srcAudioPath, destAudioPath)
			if err != nil {
				logger.Warn("Failed to copy audio file: %v", err)
			} else {
				logger.Debug("Audio file copied to: %s", destAudioPath)
			}
		} else {
			logger.Warn("Source audio file not found: %s", srcAudioPath)
		}
	}
	
	e.chart.modified = false
	logger.Debug("Chart saved to folder: %s", songFolderPath)
}

// copyFile copies a file from src to dst
func (e *EditorState) copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()
	
	_, err = io.Copy(destFile, sourceFile)
	return err
}

// loadChartFromFolder loads a chart from a folder path (used during initialization)
func (e *EditorState) loadChartFromFolder(folderPath string) {
	// Validate that the folder contains required files
	songJSONPath := filepath.Join(folderPath, "song.json")
	audioPath := filepath.Join(folderPath, "audio.ogg")
	
	// Check if song.json exists
	if _, err := os.Stat(songJSONPath); os.IsNotExist(err) {
		logger.Error("Chart folder missing song.json: %s", folderPath)
		e.createEmptyChart()
		return
	}
	
	// Check if audio.ogg exists
	if _, err := os.Stat(audioPath); os.IsNotExist(err) {
		logger.Warn("Chart folder missing audio.ogg: %s", folderPath)
		audioPath = "" // Clear audio path if not found
	}
	
	// Read the song.json file
	data, err := os.ReadFile(songJSONPath)
	if err != nil {
		logger.Error("Failed to read chart file: %v", err)
		e.createEmptyChart()
		return
	}
	
	var songData schema.SongDataV2
	err = json.Unmarshal(data, &songData)
	if err != nil {
		logger.Error("Failed to parse chart file: %v", err)
		e.createEmptyChart()
		return
	}
	
	// Load the first chart for editing
	if len(songData.Charts) == 0 {
		logger.Warn("No charts found in file")
		e.createEmptyChart()
		return
	}
	
	var firstChartData schema.ChartDataV2
	for _, chartData := range songData.Charts {
		firstChartData = chartData
		break
	}
	
	// Set the audio path in the song data if found
	if audioPath != "" {
		songData.Audio.File = audioPath
	}
	
	e.importFromJSON(&songData, &firstChartData)
	logger.Debug("Chart loaded from folder during initialization: %s", folderPath)
}

func (e *EditorState) loadChart() {
	// Select folder containing the chart
	folderPath, err := system.SelectFolderDialog("Select chart folder")
	if err != nil {
		logger.Warn("Failed to open folder dialog: %v", err)
		return
	}
	
	if folderPath == "" {
		return // User cancelled
	}
	
	// Validate that the folder contains required files
	songJSONPath := filepath.Join(folderPath, "song.json")
	audioPath := filepath.Join(folderPath, "audio.ogg")
	
	// Check if song.json exists
	if _, err := os.Stat(songJSONPath); os.IsNotExist(err) {
		logger.Error("Chart folder missing song.json: %s", folderPath)
		return
	}
	
	// Check if audio.ogg exists
	if _, err := os.Stat(audioPath); os.IsNotExist(err) {
		logger.Warn("Chart folder missing audio.ogg: %s", folderPath)
		audioPath = "" // Clear audio path if not found
	}
	
	// Read the song.json file
	data, err := os.ReadFile(songJSONPath)
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
	
	// Set the audio path in the song data if found
	if audioPath != "" {
		songData.Audio.File = audioPath
	}
	
	e.importFromJSON(&songData, &firstChartData)
	logger.Debug("Chart loaded from folder: %s", folderPath)
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
				Time:     note.Target,
				Type:     schema.NoteTypeTap,
				TimeBeat: note.TargetBeat,
			}
			
			if note.IsHoldNote() {
				jsonNotes[i].Type = schema.NoteTypeHold
				jsonNotes[i].Duration = note.TargetRelease - note.Target
				jsonNotes[i].DurationBeat = note.TargetReleaseBeat - note.TargetBeat
			}
		}
		
		chartData.Tracks[trackStr] = jsonNotes
	}
	
	// Convert events to JSON format
	for _, event := range e.chart.events {
		eventData := schema.EventData{
			Time: event.GetTime(),
			Type: event.GetType(),
		}
		
		// Add beat position if available
		switch ev := event.(type) {
		case *types.AudioStartEvent:
			eventData.TimeBeat = ev.BaseEvent.TimeBeat
		case *types.CountdownBeatEvent:
			eventData.TimeBeat = ev.BaseEvent.TimeBeat
		// Add other event types as needed
		}
		
		chartData.Events = append(chartData.Events, eventData)
	}
	
	// Create song data
	songData := &schema.SongDataV2{
		Schema:  "https://slaptrax.dev/schema/song/v2.json",
		Version: 2,
		Metadata: schema.SongMetadata{
			Title:           e.song.Title,
			Artist:          e.song.Artist,
			BPM:             int(e.bpm), // Use current editor BPM
			PreviewStart:    30000, // Default preview start
			Duration:      int64(e.song.Length),
			ChartedBy:       "Chart Editor",
			Version:         "1.0.0",
			DifficultyRange: [2]int{e.chart.metadata.Difficulty, e.chart.metadata.Difficulty},
		},
		Audio: schema.AudioInfo{
			File:   filepath.Base(e.song.AudioPath),
			Offset: e.audioOffset, // Save current audio offset
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
	// Stop any existing playback and cleanup audio
	e.stopPlayback()
	
	// Create song from metadata
	e.song = &types.Song{
		Title:     songData.Metadata.Title,
		Artist:    songData.Metadata.Artist,
		BPM:       songData.Metadata.BPM,
		Length:  int(songData.Metadata.Duration),
		AudioPath: songData.Audio.File, // This might need full path resolution
		Charts:    make(map[types.Difficulty]*types.Chart),
	}
	
	// Set editor BPM and audio offset from loaded data
	e.bpm = float64(songData.Metadata.BPM)
	e.audioOffset = songData.Audio.Offset
	
	// Reset audio initialization flag for new song
	e.audioInitialized = false
	
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
			
			// Prefer beat-based data if available, fall back to millisecond data
			if noteData.TimeBeat > 0 {
				// Use beat-based positioning
				releaseBeat := float64(0)
				if noteData.Type == schema.NoteTypeHold && noteData.DurationBeat > 0 {
					releaseBeat = noteData.TimeBeat + noteData.DurationBeat
				}
				note = types.NewNoteFromBeats(trackName, noteData.TimeBeat, releaseBeat, e.bpm)
			} else {
				// Fall back to millisecond-based positioning and calculate beats
				if noteData.Type == schema.NoteTypeHold && noteData.Duration > 0 {
					note = e.createNoteAtTime(trackName, noteData.Time, noteData.Time+noteData.Duration)
				} else {
					note = e.createNoteAtTime(trackName, noteData.Time, 0)
				}
			}
			
			e.chart.tracks[trackName] = append(e.chart.tracks[trackName], note)
		}
	}
	
	// Import events
	for _, eventData := range chartData.Events {
		event, err := types.CreateEventFromDataWithBPM(eventData, e.bpm)
		if err != nil {
			logger.Warn("Failed to create event: %v", err)
			continue
		}
		e.chart.events = append(e.chart.events, event)
	}
	
	// Keep events sorted by time
	sort.Slice(e.chart.events, func(i, j int) bool {
		return e.chart.events[i].GetTime() < e.chart.events[j].GetTime()
	})
	
	// Sync events to event manager
	e.syncEventsToManager()
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

func (e *EditorState) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
}

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