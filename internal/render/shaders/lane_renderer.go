package shaders

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/display"
	"github.com/liqmix/slaptrax/internal/types"
	"github.com/liqmix/slaptrax/internal/ui"
)

// LaneUniforms contains parameters for lane background rendering
type LaneUniforms struct {
	Point1X, Point1Y float32 // Lane corners
	Point2X, Point2Y float32
	Point3X, Point3Y float32
	Point4X, Point4Y float32
	CenterX, CenterY float32 // Vanishing point
	ColorR, ColorG, ColorB float32 // Lane color
	IsActive float32 // 1.0 if active, 0.0 if inactive
}

// MarkerUniforms contains parameters for marker rendering
type MarkerUniforms struct {
	Progress         float32 // 0.0 to 1.0
	Corner1X, Corner1Y float32 // Marker square corners
	Corner2X, Corner2Y float32
	Corner3X, Corner3Y float32
	Corner4X, Corner4Y float32
	CenterX, CenterY float32 // Vanishing point
	MarkerType       float32 // 0.0 for beat, 1.0 for measure
	TimeMs           float32 // For animations
	FadeInThreshold  float32
	FadeOutThreshold float32
}

// LaneRenderer handles shader-based lane and marker rendering
type LaneRenderer struct {
	baseVertices []ebiten.Vertex
	baseIndices  []uint16
	baseTexture  *ebiten.Image
}

// NewLaneRenderer creates a new lane renderer
func NewLaneRenderer() *LaneRenderer {
	renderer := &LaneRenderer{
		baseTexture: ui.BaseTriImg,
	}
	renderer.createBaseGeometry()
	return renderer
}

// createBaseGeometry creates a simple quad for shader rendering
func (lr *LaneRenderer) createBaseGeometry() {
	lr.baseIndices = []uint16{0, 1, 2, 0, 2, 3}
}

// createLaneBounds calculates the bounding box for a lane
func (lr *LaneRenderer) createLaneBounds(trackPoints []*ui.Point, centerPoint *ui.Point) (float32, float32, float32, float32) {
	centerX, centerY := centerPoint.ToRender32()
	
	minX := centerX
	minY := centerY
	maxX := centerX
	maxY := centerY
	
	// Include all track points
	for _, point := range trackPoints {
		if point == nil {
			continue
		}
		
		pointX, pointY := point.ToRender32()
		
		if pointX < minX {
			minX = pointX
		}
		if pointX > maxX {
			maxX = pointX
		}
		if pointY < minY {
			minY = pointY
		}
		if pointY > maxY {
			maxY = pointY
		}
	}
	
	// Add padding
	padding := float32(50)
	return minX - padding, minY - padding, maxX + padding, maxY + padding
}

// createBoundedGeometry creates vertices for a bounded area
func (lr *LaneRenderer) createBoundedGeometry(minX, minY, maxX, maxY float32) []ebiten.Vertex {
	return []ebiten.Vertex{
		{DstX: minX, DstY: minY, SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: maxX, DstY: minY, SrcX: 1, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: maxX, DstY: maxY, SrcX: 1, SrcY: 1, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: minX, DstY: maxY, SrcX: 0, SrcY: 1, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
	}
}

// CreateLaneUniforms creates uniforms for lane background rendering
func CreateLaneUniforms(trackName types.TrackName, trackPoints []*ui.Point, centerPoint *ui.Point, isActive bool) *LaneUniforms {
	if len(trackPoints) < 3 || centerPoint == nil {
		return nil
	}
	
	uniforms := &LaneUniforms{}
	
	// Set lane corners (use first 4 points if available)
	uniforms.Point1X, uniforms.Point1Y = trackPoints[0].ToRender32()
	uniforms.Point2X, uniforms.Point2Y = trackPoints[1].ToRender32()
	uniforms.Point3X, uniforms.Point3Y = trackPoints[2].ToRender32()
	
	// Use 4th point if available, otherwise use 3rd point
	if len(trackPoints) > 3 && trackPoints[3] != nil {
		uniforms.Point4X, uniforms.Point4Y = trackPoints[3].ToRender32()
	} else {
		uniforms.Point4X, uniforms.Point4Y = uniforms.Point3X, uniforms.Point3Y
	}
	
	// Set center point
	uniforms.CenterX, uniforms.CenterY = centerPoint.ToRender32()
	
	// Set lane color
	color := trackName.NoteColor()
	uniforms.ColorR = float32(color.R) / 255.0
	uniforms.ColorG = float32(color.G) / 255.0
	uniforms.ColorB = float32(color.B) / 255.0
	
	// Set active state
	uniforms.IsActive = 0.0
	if isActive {
		uniforms.IsActive = 1.0
	}
	
	return uniforms
}

// CreateMarkerUniforms creates uniforms for marker rendering
func CreateMarkerUniforms(progress float64, markerCorners []*ui.Point, centerPoint *ui.Point, ismeasure bool) *MarkerUniforms {
	if len(markerCorners) < 4 || centerPoint == nil {
		return nil
	}
	
	uniforms := &MarkerUniforms{}
	
	// Set progress
	uniforms.Progress = float32(progress)
	
	// Set marker corners
	uniforms.Corner1X, uniforms.Corner1Y = markerCorners[0].ToRender32()
	uniforms.Corner2X, uniforms.Corner2Y = markerCorners[1].ToRender32()
	uniforms.Corner3X, uniforms.Corner3Y = markerCorners[2].ToRender32()
	uniforms.Corner4X, uniforms.Corner4Y = markerCorners[3].ToRender32()
	
	// Set center point
	uniforms.CenterX, uniforms.CenterY = centerPoint.ToRender32()
	
	// Set marker type
	uniforms.MarkerType = 0.0 // Beat marker
	if ismeasure {
		uniforms.MarkerType = 1.0 // Measure marker
	}
	
	// Set time for animations
	uniforms.TimeMs = float32(time.Now().UnixMilli() % 100000)
	
	// Set fade thresholds (delayed compared to notes to avoid interference)
	uniforms.FadeInThreshold = 0.02
	uniforms.FadeOutThreshold = 0.05
	
	return uniforms
}

// RenderLaneBackground renders a lane background using the lane shader
func (lr *LaneRenderer) RenderLaneBackground(img *ebiten.Image, trackName types.TrackName, trackPoints []*ui.Point, centerPoint *ui.Point, isActive bool) {
	if Manager == nil {
		return
	}
	
	laneShader := Manager.GetLaneShader()
	if laneShader == nil {
		return
	}
	
	uniforms := CreateLaneUniforms(trackName, trackPoints, centerPoint, isActive)
	if uniforms == nil {
		return
	}
	
	// Create bounded geometry for this lane
	minX, minY, maxX, maxY := lr.createLaneBounds(trackPoints, centerPoint)
	vertices := lr.createBoundedGeometry(minX, minY, maxX, maxY)
	
	options := &ebiten.DrawTrianglesShaderOptions{}
	options.Blend = ebiten.BlendSourceOver
	options.Uniforms = map[string]interface{}{
		"Point1X":   uniforms.Point1X,
		"Point1Y":   uniforms.Point1Y,
		"Point2X":   uniforms.Point2X,
		"Point2Y":   uniforms.Point2Y,
		"Point3X":   uniforms.Point3X,
		"Point3Y":   uniforms.Point3Y,
		"Point4X":   uniforms.Point4X,
		"Point4Y":   uniforms.Point4Y,
		"CenterX":   uniforms.CenterX,
		"CenterY":   uniforms.CenterY,
		"ColorR":    uniforms.ColorR,
		"ColorG":    uniforms.ColorG,
		"ColorB":    uniforms.ColorB,
		"IsActive":  uniforms.IsActive,
	}
	
	img.DrawTrianglesShader(vertices, lr.baseIndices, laneShader, options)
}

// RenderMarker renders a measure/beat marker using the marker shader
func (lr *LaneRenderer) RenderMarker(img *ebiten.Image, progress float64, markerCorners []*ui.Point, centerPoint *ui.Point, ismeasure bool) {
	if Manager == nil {
		return
	}
	
	markerShader := Manager.GetMarkerShader()
	if markerShader == nil {
		return
	}
	
	uniforms := CreateMarkerUniforms(progress, markerCorners, centerPoint, ismeasure)
	if uniforms == nil {
		return
	}
	
	// Create bounded geometry for this marker
	minX, minY, maxX, maxY := lr.createMarkerBounds(markerCorners, centerPoint, float32(progress))
	vertices := lr.createBoundedGeometry(minX, minY, maxX, maxY)
	
	options := &ebiten.DrawTrianglesShaderOptions{}
	options.Blend = ebiten.BlendSourceOver
	options.Uniforms = map[string]interface{}{
		"Progress":         uniforms.Progress,
		"Corner1X":         uniforms.Corner1X,
		"Corner1Y":         uniforms.Corner1Y,
		"Corner2X":         uniforms.Corner2X,
		"Corner2Y":         uniforms.Corner2Y,
		"Corner3X":         uniforms.Corner3X,
		"Corner3Y":         uniforms.Corner3Y,
		"Corner4X":         uniforms.Corner4X,
		"Corner4Y":         uniforms.Corner4Y,
		"CenterX":          uniforms.CenterX,
		"CenterY":          uniforms.CenterY,
		"MarkerType":       uniforms.MarkerType,
		"TimeMs":           uniforms.TimeMs,
		"FadeInThreshold":  uniforms.FadeInThreshold,
		"FadeOutThreshold": uniforms.FadeOutThreshold,
	}
	
	img.DrawTrianglesShader(vertices, lr.baseIndices, markerShader, options)
}

// createMarkerBounds calculates bounds for a marker
func (lr *LaneRenderer) createMarkerBounds(markerCorners []*ui.Point, centerPoint *ui.Point, progress float32) (float32, float32, float32, float32) {
	centerX, centerY := centerPoint.ToRender32()
	
	minX := centerX
	minY := centerY
	maxX := centerX
	maxY := centerY
	
	// Calculate interpolated marker positions
	for _, corner := range markerCorners {
		if corner == nil {
			continue
		}
		
		cornerX, cornerY := corner.ToRender32()
		interpX := centerX + (cornerX-centerX)*progress
		interpY := centerY + (cornerY-centerY)*progress
		
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
	
	// Add padding for marker width
	padding := float32(10)
	return minX - padding, minY - padding, maxX + padding, maxY + padding
}

// Global lane renderer instance
var LaneRendererInstance *LaneRenderer

// RenderTunnelBackground renders the tunnel background constrained to play area using the tunnel shader
func (lr *LaneRenderer) RenderTunnelBackground(img *ebiten.Image, centerPoint *ui.Point, playLeft, playRight, playTop, playBottom float32) {
	if Manager == nil {
		return
	}
	
	tunnelShader := Manager.GetTunnelShader()
	if tunnelShader == nil {
		return
	}
	
	centerX, centerY := centerPoint.ToRender32()
	
	// Convert play area bounds to render coordinates
	renderWidth, renderHeight := display.Window.RenderSize()
	playLeftRender := playLeft * float32(renderWidth)
	playRightRender := playRight * float32(renderWidth)
	playTopRender := playTop * float32(renderHeight)
	playBottomRender := playBottom * float32(renderHeight)
	
	// Create geometry constrained to play area bounds
	vertices := []ebiten.Vertex{
		{DstX: playLeftRender, DstY: playTopRender, SrcX: 0, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: playRightRender, DstY: playTopRender, SrcX: 1, SrcY: 0, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: playRightRender, DstY: playBottomRender, SrcX: 1, SrcY: 1, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
		{DstX: playLeftRender, DstY: playBottomRender, SrcX: 0, SrcY: 1, ColorR: 1, ColorG: 1, ColorB: 1, ColorA: 1},
	}
	
	playAreaWidth := playRightRender - playLeftRender
	playAreaHeight := playBottomRender - playTopRender
	
	options := &ebiten.DrawTrianglesShaderOptions{}
	options.Blend = ebiten.BlendSourceOver
	options.Uniforms = map[string]interface{}{
		"CenterX":        centerX,
		"CenterY":        centerY,
		"PlayAreaWidth":  playAreaWidth,
		"PlayAreaHeight": playAreaHeight,
	}
	
	img.DrawTrianglesShader(vertices, lr.baseIndices, tunnelShader, options)
}

// InitLaneRenderer initializes the lane renderer
func InitLaneRenderer() {
	LaneRendererInstance = NewLaneRenderer()
}