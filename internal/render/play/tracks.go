package play

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
)

var (
	guideColor = color.RGBA{
		R: types.Gray.R,
		G: types.Gray.G,
		B: types.Gray.B,
		A: 150,
	}
)

func (r *Play) renderTracks(screen *ebiten.Image) {
	for _, track := range types.TrackNames() {
		pts := notePoints[track]
		r.renderLaneBackground(screen, pts)
	}

	// Then draw the area for combo / hit display
	// It's a box with rounded corners at the very center of main tracks
	// It will obscure some of the vanishing point lines

	// // Actual box
	// border := 0.05
	// pos.X += border / 2
	// pos.Y += border / 2
	// size.X -= border
	// size.Y -= border
	// ui.DrawFilledRect(screen, pos, size, types.Black)
}

func (r *Play) renderLaneBackground(img *ebiten.Image, points []*ui.Point) {
	// vc := ui.NewVectorCollection()

	// left := points[0]
	// center := points[1]
	// right := points[2]

	// // Draw the guide lines to vanishing point
	// guideLineLeft := getVectorPath([]Point{
	// 	left,
	// 	vanishingPoint,
	// }, 0)

	// guideLineRight := getVectorPath([]Point{
	// 	right,
	// 	vanishingPoint,
	// }, 0)

	// vs, is = guideLineLeft.AppendVerticesAndIndicesForStroke(nil, nil, &opts)
	// colorVertices(vs, guideColor)
	// img.DrawTriangles(vs, is, baseImg, nil)

	// vs, is = guideLineRight.AppendVerticesAndIndicesForStroke(vs, is, &opts)
	// colorVertices(vs, guideColor)
	// img.DrawTriangles(vs, is, baseImg, nil)
	// pathWidth, _ := ui.Point{
	// 	X: 0.001,
	// 	Y: 0,
	// }.ToRender32()
	// guidePaths := ui.GetDashedPaths(
	// 	center,
	// 	&playCenterPoint,
	// )
	// for _, guidePath := range guidePaths {
	// 	guidePath.Draw(img, pathWidth, guideColor)
	// }
}

// func (r *Play) renderEdgeTracks(screen *ebiten.Image) {
// 	// Draw the top track
// 	r.renderLaneBackground(screen, &LaneConfig{
// 		CurveAmount:    0,
// 		VanishingPoint: edgeCenter,
// 		Left: ui.Point{
// 			X: edgeCenter.X,
// 			Y: edgeTop,
// 		},
// 		Center: ui.Point{
// 			X: edgeLeft,
// 			Y: edgeTop,
// 		},
// 		Right: ui.Point{
// 			X: edgeLeft,
// 			Y: edgeCenter.Y,
// 		},
// 	})
// 	r.renderLaneBackground(screen, &LaneConfig{
// 		CurveAmount:    0,
// 		VanishingPoint: edgeCenter,
// 		Left: ui.Point{
// 			X: edgeRight,
// 			Y: edgeCenter.Y,
// 		},
// 		Center: ui.Point{
// 			X: edgeRight,
// 			Y: edgeTop,
// 		},
// 		Right: ui.Point{
// 			X: edgeCenter.X,
// 			Y: edgeTop,
// 		},
// 	})

// 	// Draw the bottom track
// 	r.renderLaneBackground(screen, laneConfigs[types.EdgeTap1])
// 	r.renderLaneBackground(screen, laneConfigs[types.EdgeTap2])
// 	r.renderLaneBackground(screen, laneConfigs[types.EdgeTap3])

// 	r.drawJudgementLine(screen, types.EdgeTap1)
// 	r.drawJudgementLine(screen, types.EdgeTap2)
// 	r.drawJudgementLine(screen, types.EdgeTap3)
// }

// func (r *Play) renderTracks(screen *ebiten.Image) {
// 	r.renderMainTracks(screen)
// 	if !r.hideEdgeTracks {
// 		r.renderEdgeTracks(screen)
// 	}
// }
