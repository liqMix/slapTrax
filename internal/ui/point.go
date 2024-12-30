package ui

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/liqmix/ebiten-holiday-2024/internal/user"
)

// Create a 1x1 white image as the base texture
func createBaseTriImg() *ebiten.Image {
	img := ebiten.NewImage(1, 1)
	img.Fill(color.White)
	return img
}

var BaseTriImg = createBaseTriImg()

// Actual screen coordinates
type Point struct {
	X, Y float32
}

// Ratio of render coordinates (0-1)
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

func GetVectorPath(points []Point, curveAmount float64) vector.Path {
	path := vector.Path{}
	if len(points) < 1 {
		return path
	}

	start := points[0]
	path.MoveTo(start.X, start.Y)

	for i := 1; i < len(points); i++ {
		next := points[i]
		path.LineTo(next.X, next.Y)
	}
	return path
}

func ColorVertices(vs []ebiten.Vertex, color color.RGBA) {
	for i := range vs {
		vs[i].ColorR = float32(color.R) / 255
		vs[i].ColorG = float32(color.G) / 255
		vs[i].ColorB = float32(color.B) / 255
		vs[i].ColorA = float32(color.A) / 255
	}
}

func DrawDashedLine(screen *ebiten.Image, start Point, end Point, dashLength float32, gapLength float32, color color.RGBA) {
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
		path := GetVectorPath([]Point{dashStart, dashEnd}, 0)

		vs, is := path.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{
			Width: 1,
		})

		ColorVertices(vs, color)
		screen.DrawTriangles(vs, is, BaseTriImg, nil)
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

		path := GetVectorPath([]Point{finalStart, finalEnd}, 0)

		vs, is := path.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{
			Width: 1,
		})

		ColorVertices(vs, color)
		screen.DrawTriangles(vs, is, BaseTriImg, nil)
	}
}
