package standard

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/liqmix/ebiten-holiday-2024/internal/song"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
)

// Get the rendering of the judgement line
func (r *Standard) getJudgementLine(config *LaneConfig) vector.Path {
	left := config.Left.RenderPoint()
	right := config.Right.RenderPoint()
	center := config.Center.RenderPoint()

	if left.X == right.X || left.Y == right.Y {
		return ui.GetVectorPath([]ui.Point{
			left,
			right,
		}, config.CurveAmount)
	}

	return ui.GetVectorPath([]ui.Point{
		left,
		center,
		right,
	}, config.CurveAmount)
}

func (r *Standard) drawJudgementLine(screen *ebiten.Image, trackName song.TrackName) {
	config, ok := laneConfigs[trackName]
	if !ok {
		return
	}
	img := r.getJudgementLine(config)
	var width float32 = 2.0
	if r.state.IsTrackPressed(trackName) {
		width = 8.0
	}
	vs, is := img.AppendVerticesAndIndicesForStroke(nil, nil, &vector.StrokeOptions{
		Width: width,
	})
	ui.ColorVertices(vs, getNoteColor(trackName))
	screen.DrawTriangles(vs, is, ui.BaseTriImg, nil)
}
