package shaders

import (
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
	
	// Set fade thresholds - start fade much earlier to reduce center clutter
	// These values work with smoothProgress - use small values due to perspective compression
	uniforms.FadeInThreshold = 0.02  // Start fading in later to reduce center clutter
	uniforms.FadeOutThreshold = 0.06  // Reach full visibility quickly after fade starts
	
	return uniforms
}

// CreateHoldNoteUniforms creates uniforms for a hold note
func CreateHoldNoteUniforms(track types.TrackName, note *types.Note, trackPoints []*ui.Point, centerPoint *ui.Point) *HoldNoteUniforms {
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
	
	// Adjust alpha based on hold state
	if note.WasHit() {
		if note.WasReleased() {
			holdUniforms.ColorA = 50.0 / 255.0
		} else {
			holdUniforms.ColorA = 200.0 / 255.0
		}
	} else {
		holdUniforms.ColorA = 100.0 / 255.0
	}
	
	return holdUniforms
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
	}
	return append(base, hold...)
}

// Helper functions duplicated from play package to avoid circular imports
const (
	minT        = 0.01 // Small value for vanishing point calculation (matches play package)
	noteWidth   = 18.0
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