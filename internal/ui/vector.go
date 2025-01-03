package ui

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
)

// Create a 1x1 white image as the base texture
func createBaseTriImg() *ebiten.Image {
	img := ebiten.NewImage(1, 1)
	img.Fill(color.White)
	return img
}

var BaseTriImg = createBaseTriImg()

// Normalized coordinates (0.0-1.0)
type Point struct {
	X, Y float64
}

func (p Point) V() (float64, float64) {
	return p.X, p.Y
}

func (p Point) ToRender() (float64, float64) {
	x, y := types.Window.RenderSize()
	return p.X * float64(x), p.Y * float64(y)
}

func (p Point) ToRender32() (float32, float32) {
	x, y := p.ToRender()
	return float32(x), float32(y)
}

func PointFromRender(x, y float64) Point {
	rx, ry := types.Window.RenderSize()
	return Point{x / float64(rx), y / float64(ry)}
}

type VectorPath struct {
	vector.Path
}

func (d *VectorPath) Draw(img *ebiten.Image, width float32, color color.RGBA) {
	vs, is := d.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{
		Width: width,
	})
	ColorVertices(vs, color)
	img.DrawTriangles(vs, is, BaseTriImg, nil)
}

func GetVectorPath(points []*Point) *VectorPath {
	if len(points) < 1 {
		return &VectorPath{}
	}

	path := VectorPath{}
	path.MoveTo(points[0].ToRender32())
	for i := 1; i < len(points); i++ {
		path.LineTo(points[i].ToRender32())
	}
	return &path
}

func ColorVertices(vs []ebiten.Vertex, c color.RGBA) {
	r, g, b, a := float32(c.R)/255, float32(c.G)/255, float32(c.B)/255, float32(c.A)/255

	for i := range vs {
		vs[i].ColorR = r
		vs[i].ColorG = g
		vs[i].ColorB = b
		vs[i].ColorA = a
	}
}

func GetDashedPaths(start *Point, end *Point) []*VectorPath {
	dashLength, gapLength := 0.01, 0.01
	dx := end.X - start.X
	dy := end.Y - start.Y
	length := math.Sqrt(float64(dx*dx + dy*dy))

	// Normalize direction vector
	nx := dx / length
	ny := dy / length

	// Calculate number of segments
	totalLength := dashLength + gapLength
	segments := length / totalLength

	// Draw dash segments
	paths := []*VectorPath{}
	for i := 0; i < int(segments); i++ {
		dashStart := &Point{
			X: start.X + nx*float64(i)*totalLength,
			Y: start.Y + ny*float64(i)*totalLength,
		}
		dashEnd := &Point{
			X: dashStart.X + nx*dashLength,
			Y: dashStart.Y + ny*dashLength,
		}
		paths = append(paths, GetVectorPath([]*Point{dashStart, dashEnd}))
	}
	return paths
}
