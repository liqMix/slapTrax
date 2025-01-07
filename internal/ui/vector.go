package ui

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/liqmix/ebiten-holiday-2024/internal/cache"
	"github.com/liqmix/ebiten-holiday-2024/internal/display"
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

func (p *Point) Copy() *Point {
	return &Point{X: p.X, Y: p.Y}
}

func (p *Point) Translate(dx, dy float64) {
	p.X += dx
	p.Y += dy
}
func (p *Point) TranslateX(dx float64) {
	p.X += dx
}
func (p *Point) TranslateY(dy float64) {
	p.Y += dy
}

func (p *Point) V() (float64, float64) {
	return p.X, p.Y
}

func (p *Point) ToRender() (float64, float64) {
	x, y := display.Window.RenderSize()
	return p.X * float64(x), p.Y * float64(y)
}
func (p *Point) ToRenderInt() (int, int) {
	x, y := p.ToRender()
	return int(x), int(y)
}

func (p *Point) ToRender32() (float32, float32) {
	x, y := p.ToRender()
	return float32(x), float32(y)
}

func PointFromRender(x, y float64) *Point {
	rx, ry := display.Window.RenderSize()
	return &Point{x / float64(rx), y / float64(ry)}
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
		if points[i] == nil {
			continue
		}
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

// Vector collection for rendering multiple paths in one draw call// Vector collection for rendering multiple paths in one draw call
type VectorCollection struct {
	vertices []ebiten.Vertex
	indices  []uint16
	vertIdx  int
	idxIdx   int
}

const (
	initSize   = 1024
	growFactor = 2
)

func NewVectorCollection() *VectorCollection {
	return &VectorCollection{
		vertices: make([]ebiten.Vertex, initSize),
		indices:  make([]uint16, initSize),
	}
}

func (vc *VectorCollection) growVerts(needed int) {
	newCap := cap(vc.vertices)
	for newCap < vc.vertIdx+needed {
		newCap *= growFactor
	}
	newSlice := make([]ebiten.Vertex, newCap)
	copy(newSlice, vc.vertices[:vc.vertIdx])
	vc.vertices = newSlice
}

func (vc *VectorCollection) growIndices(needed int) {
	newCap := cap(vc.indices)
	for newCap < vc.idxIdx+needed {
		newCap *= growFactor
	}
	newSlice := make([]uint16, newCap)
	copy(newSlice, vc.indices[:vc.idxIdx])
	vc.indices = newSlice
}

func (vc *VectorCollection) Add(verts []ebiten.Vertex, inds []uint16) {
	vlen, ilen := len(verts), len(inds)
	if vlen == 0 || ilen == 0 {
		return
	}

	// Check capacity and grow if needed
	if vc.vertIdx+vlen > cap(vc.vertices) {
		vc.growVerts(vlen)
	}
	if vc.idxIdx+ilen > cap(vc.indices) {
		vc.growIndices(ilen)
	}

	// Direct slice copying
	offset := uint16(vc.vertIdx)
	copy(vc.vertices[vc.vertIdx:], verts)

	// Manual loop is faster than range for small slices
	dst := vc.indices[vc.idxIdx:]
	for i := 0; i < ilen; i++ {
		dst[i] = inds[i] + offset
	}

	vc.vertIdx += vlen
	vc.idxIdx += ilen
}

func (vc *VectorCollection) AddPath(path *cache.CachedPath) {
	if path == nil {
		return
	}
	vc.Add(path.Vertices, path.Indices)
}

func (vc *VectorCollection) Draw(img *ebiten.Image) {
	if vc.vertIdx == 0 {
		return
	}
	img.DrawTriangles(vc.vertices[:vc.vertIdx], vc.indices[:vc.idxIdx], BaseTriImg, nil)
}

func (vc *VectorCollection) Clear() {
	vc.vertIdx = 0
	vc.idxIdx = 0
}
