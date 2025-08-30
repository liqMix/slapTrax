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
	// Create a simple quad covering the entire screen
	// The shader will handle positioning and scaling
	width := float32(1920) // Max expected screen width
	height := float32(1080) // Max expected screen height
	
	sr.baseVertices = []ebiten.Vertex{
		{DstX: 0, DstY: 0, SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: width, DstY: 0, SrcX: 1, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: width, DstY: height, SrcX: 1, SrcY: 1, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: 0, DstY: height, SrcX: 0, SrcY: 1, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
	}
	
	sr.baseIndices = []uint16{0, 1, 2, 0, 2, 3}
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
	}
	
	img.DrawTrianglesShader(sr.baseVertices, sr.baseIndices, Manager.noteShader, options)
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
	}
	
	img.DrawTrianglesShader(sr.baseVertices, sr.baseIndices, Manager.holdNoteShader, options)
}

// Global shader renderer instance
var Renderer *ShaderRenderer

// InitRenderer initializes the shader renderer
func InitRenderer() {
	Renderer = NewShaderRenderer()
}