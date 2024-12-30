package standard

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/liqmix/ebiten-holiday-2024/internal/cache"
	"github.com/liqmix/ebiten-holiday-2024/internal/song"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
	"github.com/liqmix/ebiten-holiday-2024/internal/user"
)

func (r *Standard) createLaneBackground(config *LaneConfig) *ebiten.Image {
	s := user.Settings()
	img := ebiten.NewImage(s.RenderWidth, s.RenderHeight)

	// left := config.Left.RenderPoint()
	// right := config.Right.RenderPoint()
	center := config.Center.RenderPoint()
	vanishingPoint := config.VanishingPoint.RenderPoint()

	// Draw the target line
	// targetLine := r.getJudgementLine(config)

	// Draw the guide lines to vanishing point
	// guideLineLeft := getVectorPath([]Point{
	// 	left,
	// 	vanishingPoint,
	// }, 0)

	// guideLineRight := getVectorPath([]Point{
	// 	right,
	// 	vanishingPoint,
	// }, 0)

	var guideColor = color.RGBA{
		R: gray.R,
		G: gray.G,
		B: gray.B,
		A: 150,
	}

	ui.DrawDashedLine(img,
		center,
		vanishingPoint,
		10,
		10,
		guideColor,
	)

	// vs, is = guideLineLeft.AppendVerticesAndIndicesForStroke(nil, nil, &opts)
	// colorVertices(vs, guideColor)
	// img.DrawTriangles(vs, is, baseImg, nil)

	// vs, is = guideLineRight.AppendVerticesAndIndicesForStroke(vs, is, &opts)
	// colorVertices(vs, guideColor)
	// img.DrawTriangles(vs, is, baseImg, nil)
	return img
}

func (r *Standard) drawMainTracks(screen *ebiten.Image) {
	var laneConfig *LaneConfig
	var img *ebiten.Image
	var ok bool

	for _, track := range song.MainTracks {
		laneConfig = laneConfigs[track]

		img, ok = cache.GetImage(string(track))
		if !ok {
			img = r.createLaneBackground(laneConfig)
			cache.SetImage(string(track), img)
		}
		screen.DrawImage(img, nil)

		// Judgement line
		// Enlarge if track is active
		r.drawJudgementLine(screen, track)
	}

	// Then draw the area for combo / hit display
	// It's a box with rounded corners at the very center of main tracks
	// It will obscure some of the vanishing point lines
	s := user.Settings()
	img, ok = cache.GetImage("play.maincombo")
	if !ok {
		img = ebiten.NewImage(s.RenderWidth, s.RenderHeight)

		// Draw the box
		x := (mainCenter.X - (mainScoreScale / 2)) * float64(s.RenderWidth)
		y := (mainCenter.Y - (mainScoreScale / 2)) * float64(s.RenderHeight)
		width := mainScoreScale * float64(s.RenderWidth)
		height := mainScoreScale * float64(s.RenderHeight)

		// Fill
		border := 0.005 * float64(s.RenderHeight)

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

	// Draw the combo text in the center of the combo box
	combo := r.state.Score.Combo
	if combo > 0 {
		comboText := fmt.Sprintf("%d", r.state.Score.Combo)
		ui.DrawTextCenterAt(screen, comboText, int(mainCenter.X*float64(s.RenderWidth)), int(mainCenter.Y*float64(s.RenderHeight)), 1)
	}
}

func (r *Standard) drawEdgeTracks(screen *ebiten.Image) {
	if !r.renderEdgeTracks {
		return
	}

	// Draw the top track
	img, ok := cache.GetImage(string(song.EdgeTop))
	if !ok {
		s := user.Settings()
		img = ebiten.NewImage(s.RenderWidth, s.RenderHeight)
		leftHalf := r.createLaneBackground(&LaneConfig{
			CurveAmount:    0,
			VanishingPoint: edgeCenter,
			Left: ui.NormPoint{
				X: edgeCenter.X,
				Y: edgeTop,
			},
			Center: ui.NormPoint{
				X: edgeLeft,
				Y: edgeTop,
			},
			Right: ui.NormPoint{
				X: edgeLeft,
				Y: edgeCenter.Y,
			},
		})
		rightHalf := r.createLaneBackground(&LaneConfig{
			CurveAmount:    0,
			VanishingPoint: edgeCenter,
			Left: ui.NormPoint{
				X: edgeRight,
				Y: edgeCenter.Y,
			},
			Center: ui.NormPoint{
				X: edgeRight,
				Y: edgeTop,
			},
			Right: ui.NormPoint{
				X: edgeCenter.X,
				Y: edgeTop,
			},
		})
		img.DrawImage(leftHalf, nil)
		img.DrawImage(rightHalf, nil)
		cache.SetImage(string(song.EdgeTop), img)
	}
	screen.DrawImage(img, nil)
	r.drawJudgementLine(screen, song.EdgeTop)

	// Draw the bottom track
	img, ok = cache.GetImage(string(song.EdgeTap1))
	if !ok {
		s := user.Settings()
		img = ebiten.NewImage(s.RenderWidth, s.RenderHeight)
		img.DrawImage(r.createLaneBackground(laneConfigs[song.EdgeTap1]), nil)
		img.DrawImage(r.createLaneBackground(laneConfigs[song.EdgeTap2]), nil)
		img.DrawImage(r.createLaneBackground(laneConfigs[song.EdgeTap3]), nil)

		cache.SetImage(string(song.EdgeTap1), img)
	}
	screen.DrawImage(img, nil)
	r.drawJudgementLine(screen, song.EdgeTap1)
	r.drawJudgementLine(screen, song.EdgeTap2)
	r.drawJudgementLine(screen, song.EdgeTap3)
}
