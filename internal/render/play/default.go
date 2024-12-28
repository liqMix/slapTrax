package play

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/liqmix/ebiten-holiday-2024/internal/cache"
	"github.com/liqmix/ebiten-holiday-2024/internal/song"
	play "github.com/liqmix/ebiten-holiday-2024/internal/state/play"
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

// The default renderer for the play state.
type Default struct {
	state *play.State
}

func (r Default) New(s *play.State) PlayRenderer {
	baseImg.Fill(color.White)
	return &Default{
		state: s,
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
func (r *Default) drawScore(screen *ebiten.Image)    {}

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
}

var (
	offset  = 0.05
	spacing = 0.02

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
	centerTrackWidth = 0.3

	// Edge
	edgeCenter = NormPoint{
		X: edgeLeft + (edgeWidth / 2),
		Y: mainCenter.Y,
	}
	edgeLeft   = mainRight + offset*2
	edgeRight  = 1 - offset
	edgeHeight = 0.8 * mainHeight
	edgeTop    = 1 - (offset + edgeHeight)
	edgeBottom = 1 - offset
	edgeWidth  = edgeRight - edgeLeft
)

func (r *Default) drawMainTracks(screen *ebiten.Image) {
	topSpacing := spacing / 2
	bottomSpacing := (centerTrackWidth / 2) + spacing

	// Draw the rounded tracks first
	curve := 0.7

	var img *ebiten.Image
	var ok bool

	//// Left side
	// Bottom
	img, ok = cache.GetImage(string(song.LeftBottom))
	if !ok {
		img = r.createLaneBackground(LaneConfig{
			CurveAmount:    curve,
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
				X: mainCenter.X - bottomSpacing,
				Y: mainBottom,
			},
		})
		cache.SetImage(string(song.LeftBottom), img)
	}
	screen.DrawImage(img, nil)

	// Top
	img, ok = cache.GetImage(string(song.LeftTop))
	if !ok {
		img = r.createLaneBackground(LaneConfig{
			CurveAmount:    curve,
			VanishingPoint: mainCenter,
			Left: NormPoint{
				X: mainCenter.X - topSpacing,
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
		})
		cache.SetImage(string(song.LeftTop), img)
	}
	screen.DrawImage(img, nil)

	//// Right Side
	// Bottom
	img, ok = cache.GetImage(string(song.RightBottom))
	if !ok {
		img = r.createLaneBackground(LaneConfig{
			CurveAmount:    curve,
			VanishingPoint: mainCenter,
			Left: NormPoint{
				X: mainCenter.X + bottomSpacing,
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
		})
		cache.SetImage(string(song.RightBottom), img)
	}
	screen.DrawImage(img, nil)

	// Top
	img, ok = cache.GetImage(string(song.RightTop))
	if !ok {
		img = r.createLaneBackground(LaneConfig{
			CurveAmount:    curve,
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
				X: mainCenter.X + topSpacing,
				Y: mainTop,
			},
		})
		cache.SetImage(string(song.RightTop), img)
	}
	screen.DrawImage(img, nil)

	//// Then draw the straight track
	img, ok = cache.GetImage(string(song.Center))
	if !ok {
		img = r.createLaneBackground(LaneConfig{
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
		})
		cache.SetImage(string(song.Center), img)
	}
	screen.DrawImage(img, nil)
}

func (r *Default) drawEdgeTracks(screen *ebiten.Image) {
	// Draw the top track
	img, ok := cache.GetImage(string(song.EdgeTop))
	if !ok {
		s := user.Settings()
		img = ebiten.NewImage(s.RenderWidth, s.RenderHeight)
		leftHalf := r.createLaneBackground(LaneConfig{
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
		rightHalf := r.createLaneBackground(LaneConfig{
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

	// Draw the bottom track
	img, ok = cache.GetImage(string(song.EdgeTap1))
	if !ok {
		img = r.createLaneBackground(LaneConfig{
			CurveAmount:    0,
			VanishingPoint: edgeCenter,
			Left: NormPoint{
				X: edgeLeft,
				Y: edgeBottom,
			},
			Center: NormPoint{
				X: edgeCenter.X,
				Y: edgeBottom,
			},
			Right: NormPoint{
				X: edgeRight,
				Y: edgeBottom,
			},
		})
		cache.SetImage(string(song.EdgeTap1), img)
	}
	screen.DrawImage(img, nil)
}

type LaneConfig struct {
	Left           NormPoint // Bottom left point (0-1 space)
	Center         NormPoint // Bottom center point (0-1 space)
	Right          NormPoint // Bottom right point (0-1 space)
	VanishingPoint NormPoint // Convergence point (0-1 space)
	CurveAmount    float64   // 0 = right angle at center, 1 = straight line
}

func (r *Default) createLaneBackground(config LaneConfig) *ebiten.Image {
	s := user.Settings()
	img := ebiten.NewImage(s.RenderWidth, s.RenderHeight)

	left := config.Left.RenderPoint()
	right := config.Right.RenderPoint()
	center := config.Center.RenderPoint()
	// curve := config.CurveAmount
	vanishingPoint := config.VanishingPoint.RenderPoint()

	// Draw the target line
	targetLine := vector.Path{}
	targetLine.MoveTo(left.X, left.Y)
	targetLine.LineTo(center.X, center.Y)
	targetLine.LineTo(right.X, right.Y)

	// Draw the guide lines to vanishing point
	guideLineCenter := vector.Path{}
	guideLineCenter.MoveTo(center.X, center.Y)
	guideLineCenter.LineTo(vanishingPoint.X, vanishingPoint.Y)

	guideLineLeft := vector.Path{}
	guideLineLeft.MoveTo(left.X, left.Y)
	guideLineLeft.LineTo(vanishingPoint.X, vanishingPoint.Y)

	guideLineRight := vector.Path{}
	guideLineRight.MoveTo(right.X, right.Y)
	guideLineRight.LineTo(vanishingPoint.X, vanishingPoint.Y)

	opts := vector.StrokeOptions{
		Width:   4,
		LineCap: vector.LineCapRound,
	}

	vs, is := targetLine.AppendVerticesAndIndicesForStroke(nil, nil, &opts)
	img.DrawTriangles(vs, is, baseImg, nil)

	opts = vector.StrokeOptions{
		Width: 1,
	}

	var guideAlpha float32 = 0.5
	vs, is = guideLineCenter.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{
		Width:    1,
		LineJoin: vector.LineJoinRound,
	})
	for i := range vs {
		vs[i].ColorA = guideAlpha
	}
	img.DrawTriangles(vs, is, baseImg, nil)

	vs, is = guideLineLeft.AppendVerticesAndIndicesForStroke(nil, nil, &opts)
	for i := range vs {
		vs[i].ColorA = guideAlpha
	}
	img.DrawTriangles(vs, is, baseImg, nil)

	vs, is = guideLineRight.AppendVerticesAndIndicesForStroke(nil, nil, &opts)
	for i := range vs {
		vs[i].ColorA = guideAlpha
	}
	img.DrawTriangles(vs, is, baseImg, nil)
	return img
}
