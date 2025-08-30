package shaders

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/types"
	"github.com/liqmix/slaptrax/internal/ui"
)

// ShaderRenderer handles shader-based note rendering
type ShaderRenderer struct {
	// Base geometry for notes (simple quad that will be transformed by shader)
	baseVertices []ebiten.Vertex
	baseIndices  []uint16
	
	// Base texture for shader rendering
	baseTexture *ebiten.Image
}

// NewShaderRenderer creates a new shader-based renderer
func NewShaderRenderer() *ShaderRenderer {
	renderer := &ShaderRenderer{
		baseTexture: ui.BaseTriImg,
	}
	
	renderer.createBaseGeometry()
	return renderer
}

// createBaseGeometry creates a simple quad that shaders will transform
func (sr *ShaderRenderer) createBaseGeometry() {
	// Create a basic quad - we'll calculate proper bounds per note
	sr.baseIndices = []uint16{0, 1, 2, 0, 2, 3}
}

// createBoundedGeometry creates geometry bounded to the note's actual area
func (sr *ShaderRenderer) createBoundedGeometry(trackPoints []*ui.Point, centerPoint *ui.Point, progress float32) []ebiten.Vertex {
	// Calculate the bounding box for this note
	minX, minY, maxX, maxY := sr.calculateNoteBounds(trackPoints, centerPoint, progress)
	
	// Add padding to ensure we don't clip edges
	padding := float32(20)
	minX -= padding
	minY -= padding
	maxX += padding
	maxY += padding
	
	// Create vertices for the bounded quad
	return []ebiten.Vertex{
		{DstX: minX, DstY: minY, SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: maxX, DstY: minY, SrcX: 1, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: maxX, DstY: maxY, SrcX: 1, SrcY: 1, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: minX, DstY: maxY, SrcX: 0, SrcY: 1, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
	}
}

// calculateNoteBounds determines the screen-space bounding box for a note
func (sr *ShaderRenderer) calculateNoteBounds(trackPoints []*ui.Point, centerPoint *ui.Point, progress float32) (float32, float32, float32, float32) {
	// Calculate interpolated positions
	centerX, centerY := centerPoint.ToRender32()
	
	minX := centerX
	minY := centerY
	maxX := centerX
	maxY := centerY
	
	// Check all track points at the given progress
	for _, point := range trackPoints {
		if point == nil {
			continue
		}
		
		pointX, pointY := point.ToRender32()
		
		// Interpolate position based on progress
		interpX := centerX + (pointX-centerX)*progress
		interpY := centerY + (pointY-centerY)*progress
		
		// Update bounds
		if interpX < minX {
			minX = interpX
		}
		if interpX > maxX {
			maxX = interpX
		}
		if interpY < minY {
			minY = interpY
		}
		if interpY > maxY {
			maxY = interpY
		}
	}
	
	return minX, minY, maxX, maxY
}

// calculateHoldNoteBounds determines bounds for a hold note spanning from startProgress to endProgress
func (sr *ShaderRenderer) calculateHoldNoteBounds(trackPoints []*ui.Point, centerPoint *ui.Point, startProgress, endProgress float32) (float32, float32, float32, float32) {
	// Calculate bounds at both start and end progress
	startMinX, startMinY, startMaxX, startMaxY := sr.calculateNoteBounds(trackPoints, centerPoint, startProgress)
	endMinX, endMinY, endMaxX, endMaxY := sr.calculateNoteBounds(trackPoints, centerPoint, endProgress)
	
	// Combine bounds to encompass the entire hold area
	minX := min(startMinX, endMinX)
	minY := min(startMinY, endMinY)
	maxX := max(startMaxX, endMaxX)
	maxY := max(startMaxY, endMaxY)
	
	return minX, minY, maxX, maxY
}

// createHoldNoteBoundedGeometry creates geometry for a hold note
func (sr *ShaderRenderer) createHoldNoteBoundedGeometry(trackPoints []*ui.Point, centerPoint *ui.Point, startProgress, endProgress float32) []ebiten.Vertex {
	// Calculate the bounding box for this hold note
	minX, minY, maxX, maxY := sr.calculateHoldNoteBounds(trackPoints, centerPoint, startProgress, endProgress)
	
	// Add padding to ensure we don't clip edges
	padding := float32(20)
	minX -= padding
	minY -= padding
	maxX += padding
	maxY += padding
	
	// Create vertices for the bounded quad
	return []ebiten.Vertex{
		{DstX: minX, DstY: minY, SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: maxX, DstY: minY, SrcX: 1, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: maxX, DstY: maxY, SrcX: 1, SrcY: 1, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: minX, DstY: maxY, SrcX: 0, SrcY: 1, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
	}
}

// RenderNote renders a single note using shaders
func (sr *ShaderRenderer) RenderNote(img *ebiten.Image, track types.TrackName, note *types.Note, trackPoints []*ui.Point, centerPoint *ui.Point) {
	if Manager == nil || Manager.noteShader == nil {
		return
	}
	
	uniforms := CreateNoteUniforms(track, note, trackPoints, centerPoint)
	if uniforms == nil {
		return
	}
	
	// Create bounded geometry for this specific note
	vertices := sr.createBoundedGeometry(trackPoints, centerPoint, uniforms.Progress)
	
	options := &ebiten.DrawTrianglesShaderOptions{}
	options.Uniforms = map[string]interface{}{
		"Progress":   uniforms.Progress,
		"Point1X":    uniforms.Point1X,
		"Point1Y":    uniforms.Point1Y,
		"Point2X":    uniforms.Point2X,
		"Point2Y":    uniforms.Point2Y,
		"Point3X":    uniforms.Point3X,
		"Point3Y":    uniforms.Point3Y,
		"CenterX":    uniforms.CenterX,
		"CenterY":    uniforms.CenterY,
		"Width":      uniforms.Width,
		"WidthScale": uniforms.WidthScale,
		"ColorR":     uniforms.ColorR,
		"ColorG":     uniforms.ColorG,
		"ColorB":     uniforms.ColorB,
		"ColorA":     uniforms.ColorA,
		"Glow":       uniforms.Glow,
		"Solo":       uniforms.Solo,
		"TimeMs":     uniforms.TimeMs,
		"FadeInThreshold":  uniforms.FadeInThreshold,
		"FadeOutThreshold": uniforms.FadeOutThreshold,
	}
	
	img.DrawTrianglesShader(vertices, sr.baseIndices, Manager.noteShader, options)
}

// RenderHoldNote renders a hold note using shaders
func (sr *ShaderRenderer) RenderHoldNote(img *ebiten.Image, track types.TrackName, note *types.Note, trackPoints []*ui.Point, centerPoint *ui.Point) {
	if Manager == nil || Manager.holdNoteShader == nil {
		return
	}
	
	uniforms := CreateHoldNoteUniforms(track, note, trackPoints, centerPoint)
	if uniforms == nil {
		return
	}
	
	// Create bounded geometry for this hold note spanning from start to end progress
	vertices := sr.createHoldNoteBoundedGeometry(trackPoints, centerPoint, uniforms.HoldStartProgress, uniforms.HoldEndProgress)
	
	options := &ebiten.DrawTrianglesShaderOptions{}
	options.Uniforms = map[string]interface{}{
		"Progress":           uniforms.Progress,
		"Point1X":            uniforms.Point1X,
		"Point1Y":            uniforms.Point1Y,
		"Point2X":            uniforms.Point2X,
		"Point2Y":            uniforms.Point2Y,
		"Point3X":            uniforms.Point3X,
		"Point3Y":            uniforms.Point3Y,
		"CenterX":            uniforms.CenterX,
		"CenterY":            uniforms.CenterY,
		"Width":              uniforms.Width,
		"WidthScale":         uniforms.WidthScale,
		"ColorR":             uniforms.ColorR,
		"ColorG":             uniforms.ColorG,
		"ColorB":             uniforms.ColorB,
		"ColorA":             uniforms.ColorA,
		"Glow":               uniforms.Glow,
		"Solo":               uniforms.Solo,
		"TimeMs":             uniforms.TimeMs,
		"HoldStartProgress":  uniforms.HoldStartProgress,
		"HoldEndProgress":    uniforms.HoldEndProgress,
		"WasHit":             uniforms.WasHit,
		"WasReleased":        uniforms.WasReleased,
		"FadeInThreshold":    uniforms.FadeInThreshold,
		"FadeOutThreshold":   uniforms.FadeOutThreshold,
	}
	
	img.DrawTrianglesShader(vertices, sr.baseIndices, Manager.holdNoteShader, options)
}

// Global shader renderer instance
var Renderer *ShaderRenderer

// InitRenderer initializes the shader renderer
func InitRenderer() {
	Renderer = NewShaderRenderer()
}