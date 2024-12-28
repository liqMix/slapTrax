package play

import (
	"image/color"
	"math"

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
	baseImg.Fill(white)
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
	mainScoreScale   = 0.1

	// Edge
	edgeCenter = NormPoint{
		X: edgeLeft + (edgeWidth / 2),
		Y: mainCenter.Y,
	}
	edgeLeft   = mainRight + offset
	edgeRight  = 1 - offset
	edgeHeight = 0.8 * mainHeight
	edgeTop    = 1 - (offset + edgeHeight)
	edgeBottom = 1 - offset
	edgeWidth  = edgeRight - edgeLeft

	// Colors
	black = color.RGBA{0, 0, 0, 255}
	gray  = color.RGBA{100, 100, 100, 255}
	white = color.RGBA{200, 200, 200, 255}
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

	// Then draw the area for combo / hit display
	// It's a box with rounded corners at the very center of main tracks
	// It will obscure some of the vanishing point lines
	img, ok = cache.GetImage("play.maincombo")
	if !ok {
		s := user.Settings()
		img = ebiten.NewImage(s.RenderWidth, s.RenderHeight)

		// Draw the box
		x := (mainCenter.X - (mainScoreScale / 2)) * float64(s.RenderWidth)
		y := (mainCenter.Y - (mainScoreScale / 2)) * float64(s.RenderHeight)
		width := mainScoreScale * float64(s.RenderWidth)
		height := mainScoreScale * float64(s.RenderHeight)

		// Fill
		border := 0.01 * float64(s.RenderHeight)

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
			LaneCount:      3,
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
	LaneCount      int       // Number of lanes
}

func (r *Default) createLaneBackground(config LaneConfig) *ebiten.Image {
	s := user.Settings()
	img := ebiten.NewImage(s.RenderWidth, s.RenderHeight)

	left := config.Left.RenderPoint()
	right := config.Right.RenderPoint()
	center := config.Center.RenderPoint()
	// curve := config.CurveAmount
	vanishingPoint := config.VanishingPoint.RenderPoint()
	laneCount := math.Max(1, float64(config.LaneCount))

	// Draw the target line
	targetLine := vector.Path{}
	targetLine.MoveTo(left.X, left.Y)
	targetLine.LineTo(center.X, center.Y)
	targetLine.LineTo(right.X, right.Y)

	// Draw the guide lines to vanishing point
	guideLineLeft := vector.Path{}
	guideLineLeft.MoveTo(left.X, left.Y)
	guideLineLeft.LineTo(vanishingPoint.X, vanishingPoint.Y)

	guideLineRight := vector.Path{}
	guideLineRight.MoveTo(right.X, right.Y)
	guideLineRight.LineTo(vanishingPoint.X, vanishingPoint.Y)

	opts := vector.StrokeOptions{
		Width:   3,
		LineCap: vector.LineCapRound,
	}

	vs, is := targetLine.AppendVerticesAndIndicesForStroke(nil, nil, &opts)
	img.DrawTriangles(vs, is, baseImg, nil)

	opts = vector.StrokeOptions{
		Width: 1,
	}

	var guideColor = color.RGBA{
		R: gray.R,
		G: gray.G,
		B: gray.B,
		A: 150,
	}
	// Draw separate center guides lines for each lane
	// currently onyl supports straight lanes
	if laneCount > 1 {
		// Calculate lane width
		laneWidth := (right.X - left.X) / float32(laneCount)
		for i := 1; i < int(laneCount); i++ {
			from := Point{
				X: left.X + (laneWidth * float32(i)),
				Y: left.Y,
			}
			r.drawDashedLine(img,
				from,
				vanishingPoint,
				10,
				10,
				guideColor,
			)
		}
	} else {
		r.drawDashedLine(img,
			center,
			vanishingPoint,
			10,
			10,
			guideColor,
		)
	}

	vs, is = guideLineLeft.AppendVerticesAndIndicesForStroke(nil, nil, &opts)
	colorVertices(vs, guideColor)
	img.DrawTriangles(vs, is, baseImg, nil)

	vs, is = guideLineRight.AppendVerticesAndIndicesForStroke(vs, is, &opts)
	colorVertices(vs, guideColor)
	img.DrawTriangles(vs, is, baseImg, nil)
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

		path := vector.Path{}
		path.MoveTo(dashStart.X, dashStart.Y)
		path.LineTo(dashEnd.X, dashEnd.Y)

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

		path := vector.Path{}
		path.MoveTo(finalStart.X, finalStart.Y)
		path.LineTo(finalEnd.X, finalEnd.Y)

		vs, is := path.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{
			Width: 1,
		})

		colorVertices(vs, color)
		img.DrawTriangles(vs, is, baseImg, nil)
	}
}
