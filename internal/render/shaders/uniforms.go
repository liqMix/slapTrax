package shaders

import (
	"time"
	
	"github.com/liqmix/slaptrax/internal/display"
	"github.com/liqmix/slaptrax/internal/types"
	"github.com/liqmix/slaptrax/internal/ui"
	"github.com/liqmix/slaptrax/internal/user"
)

// NoteUniforms contains all parameters needed for note shader rendering
type NoteUniforms struct {
	// Progress from center (0.0) to judgment line (1.0)
	Progress float32
	
	// Track geometry - 3 points that define the note shape
	Point1X, Point1Y float32 // First point of track shape
	Point2X, Point2Y float32 // Second point (usually corner/center)
	Point3X, Point3Y float32 // Third point of track shape
	
	// Center point (vanishing point)
	CenterX, CenterY float32
	
	// Note properties
	Width      float32 // Base width before scaling
	WidthScale float32 // Additional width multiplier
	
	// Color and transparency
	ColorR, ColorG, ColorB, ColorA float32
	
	// Effects
	Glow     float32 // Glow intensity
	Solo     float32 // 1.0 if solo, 0.0 otherwise
	TimeMs   float32 // Current time in milliseconds for animations
	
	// Fade thresholds
	FadeInThreshold  float32 // Progress below which notes are invisible
	FadeOutThreshold float32 // Progress above which notes are fully visible
}

// HoldNoteUniforms contains parameters for hold note rendering
type HoldNoteUniforms struct {
	NoteUniforms
	
	// Hold note specific
	HoldStartProgress float32 // Where the hold starts (0.0-1.0)
	HoldEndProgress   float32 // Where the hold ends (0.0-1.0)
	WasHit            float32 // 1.0 if hit, 0.0 otherwise
	WasReleased       float32 // 1.0 if released, 0.0 otherwise
	IsActive          float32 // 1.0 if currently active (being held), 0.0 otherwise
	BPM               float32 // Song BPM for oscillation sync
}

// HitEffectUniforms contains parameters for hit effect rendering (shadow that travels backwards)
type HitEffectUniforms struct {
	NoteUniforms
	
	// Hit effect specific
	HitProgress    float32 // Progress from hit point (1.0) to vanishing point (0.0) - reverse travel
	EffectOpacity  float32 // Fade-out opacity (1.0 at start, 0.0 at end)
}

// CreateNoteUniforms creates uniforms for a regular note
func CreateNoteUniforms(track types.TrackName, note *types.Note, trackPoints []*ui.Point, centerPoint *ui.Point) *NoteUniforms {
	if len(trackPoints) < 3 || centerPoint == nil {
		return nil
	}
	
	uniforms := &NoteUniforms{}
	
	// Set progress
	uniforms.Progress = smoothProgress(float64(note.Progress))
	
	// Set track geometry - all 3 points that define the track shape
	uniforms.Point1X, uniforms.Point1Y = trackPoints[0].ToRender32()
	uniforms.Point2X, uniforms.Point2Y = trackPoints[1].ToRender32()
	uniforms.Point3X, uniforms.Point3Y = trackPoints[2].ToRender32()
	
	// Set center point (vanishing point)
	uniforms.CenterX, uniforms.CenterY = centerPoint.ToRender32()
	
	// Set note properties
	uniforms.Width = getNoteWidth()
	uniforms.WidthScale = 1.0
	if !note.Solo {
		uniforms.WidthScale = 1.5 // noteComboRatio equivalent
	}
	
	// Set color
	color := track.NoteColor()
	if !note.Solo {
		color = types.White.C()
	}
	uniforms.ColorR = float32(color.R) / 255.0
	uniforms.ColorG = float32(color.G) / 255.0
	uniforms.ColorB = float32(color.B) / 255.0
	uniforms.ColorA = 1.0 // Let shader handle fading based on progress and thresholds
	
	// Set effects
	uniforms.Solo = 0.0
	if note.Solo {
		uniforms.Solo = 1.0
	}
	uniforms.Glow = 0.0
	
	// Set fade thresholds - faster fade-in with earlier start for better preparation time
	// These values work with smoothProgress - use small values due to perspective compression
	uniforms.FadeInThreshold = 0.005  // Start fading in much earlier (was 0.01)
	uniforms.FadeOutThreshold = 0.025 // Reach full visibility faster (was 0.04)
	
	// Set current time for animations (use modulo to keep values manageable for sine calculations)
	uniforms.TimeMs = float32(time.Now().UnixMilli() % 100000)
	
	return uniforms
}

// CreateHoldNoteUniforms creates uniforms for a hold note
func CreateHoldNoteUniforms(track types.TrackName, note *types.Note, trackPoints []*ui.Point, centerPoint *ui.Point, bpm float32) *HoldNoteUniforms {
	baseUniforms := CreateNoteUniforms(track, note, trackPoints, centerPoint)
	if baseUniforms == nil {
		return nil
	}
	
	holdUniforms := &HoldNoteUniforms{
		NoteUniforms: *baseUniforms,
	}
	
	// Set hold-specific properties
	holdUniforms.HoldStartProgress = smoothProgress(float64(note.Progress))
	holdUniforms.HoldEndProgress = smoothProgress(float64(note.ReleaseProgress))
	
	holdUniforms.WasHit = 0.0
	if note.WasHit() {
		holdUniforms.WasHit = 1.0
	}
	
	holdUniforms.WasReleased = 0.0
	if note.WasReleased() {
		holdUniforms.WasReleased = 1.0
	}
	
	holdUniforms.IsActive = 0.0
	if note.IsActive {
		holdUniforms.IsActive = 1.0
	}
	
	// Set BPM for oscillation effects
	holdUniforms.BPM = bpm
	
	// Set base alpha for hold notes
	holdUniforms.ColorA = 1.0
	
	// Apply opacity based on active state (this handles dimming for inactive notes)
	holdUniforms.ColorA = note.GetHoldOpacity()
	
	return holdUniforms
}

// CreateHitEffectUniforms creates uniforms for a hit effect
func CreateHitEffectUniforms(track types.TrackName, note *types.Note, trackPoints []*ui.Point, centerPoint *ui.Point, hitProgress, effectOpacity float32) *HitEffectUniforms {
	baseUniforms := CreateNoteUniforms(track, note, trackPoints, centerPoint)
	if baseUniforms == nil {
		return nil
	}
	
	hitUniforms := &HitEffectUniforms{
		NoteUniforms: *baseUniforms,
	}
	
	// Set hit effect specific properties
	hitUniforms.HitProgress = hitProgress
	hitUniforms.EffectOpacity = effectOpacity
	
	// Override Progress with linear (non-smoothed) value for constant speed
	// The base CreateNoteUniforms applies smoothProgress which creates easing
	// For hit effects, we want constant speed, so use raw linear progress
	hitUniforms.Progress = hitProgress
	
	// Make the effect much more transparent than regular notes
	// Set base alpha to 0.1 (instead of 1.0) for very high transparency
	hitUniforms.ColorA = 0.1
	
	// Remove glow effects for clean shadow
	hitUniforms.Glow = 0.0
	
	return hitUniforms
}

// ToSlice converts NoteUniforms to float32 slice for shader
func (u *NoteUniforms) ToSlice() []float32 {
	return []float32{
		u.Progress,
		u.Point1X, u.Point1Y,
		u.Point2X, u.Point2Y,
		u.Point3X, u.Point3Y,
		u.CenterX, u.CenterY,
		u.Width, u.WidthScale,
		u.ColorR, u.ColorG, u.ColorB, u.ColorA,
		u.Glow, u.Solo, u.TimeMs,
		u.FadeInThreshold, u.FadeOutThreshold,
	}
}

// ToSlice converts HoldNoteUniforms to float32 slice for shader
func (u *HoldNoteUniforms) ToSlice() []float32 {
	base := u.NoteUniforms.ToSlice()
	hold := []float32{
		u.HoldStartProgress, u.HoldEndProgress,
		u.WasHit, u.WasReleased,
		u.IsActive,
		u.BPM,
	}
	return append(base, hold...)
}

// ToSlice converts HitEffectUniforms to float32 slice for shader
func (u *HitEffectUniforms) ToSlice() []float32 {
	base := u.NoteUniforms.ToSlice()
	effect := []float32{
		u.HitProgress,
		u.EffectOpacity,
	}
	return append(base, effect...)
}

// Helper functions duplicated from play package to avoid circular imports
const (
	minT        = 0.01 // Small value for vanishing point calculation (matches play package)
	noteWidth   = 36.0 // Doubled from 18.0 for thicker notes
	noteMaxAlpha = uint8(255)
)

func smoothProgress(progress float64) float32 {
	if progress >= 1 {
		return 1
	} else if progress <= 0 {
		return 0
	}
	return float32(minT / (minT + (1-minT)*(1-progress)))
}

func getNoteWidth() float32 {
	renderWidth, _ := display.Window.RenderSize()
	// scale note width based on render width
	// default defined for 1280
	return noteWidth * (float32(renderWidth) / 1280) * user.S().NoteWidth
}

func getFadeAlpha(progress float32, max uint8) uint8 {
	if progress <= 0 {
		return 0
	} else if progress >= 1 {
		return max
	}
	return uint8(float32(max) * progress)
}

func getNoteFadeAlpha(progress float32) uint8 {
	return getFadeAlpha(progress, noteMaxAlpha)
}