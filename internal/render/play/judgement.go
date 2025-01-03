package play

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
)

func (r *Play) drawJudgementLines(screen *ebiten.Image) {
	for _, track := range r.state.Tracks {
		width := judgementWidth
		if track.IsPressed() {
			width *= 2
		}
		color := track.Name.NoteColor()
		points := GetJudgementPoints(track.Name)
		path := ui.GetVectorPath(points)
		path.Draw(screen, width, color)
	}
}
