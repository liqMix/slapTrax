package play

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/state"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
)

// type Animation struct {
// 	startTime int64
// }

// func (a *Animation) Draw() float32 {
// 	return 0
// }

// The default renderer for the play state.
type Play struct {
	types.BaseRenderer
	state *state.Play
	// animations map[string]*Animation
}

func NewPlayRender(s state.State) *Play {
	p := &Play{state: s.(*state.Play)}
	p.BaseRenderer.Init(p.static)
	return p
}

func (r *Play) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	r.BaseRenderer.Draw(screen, opts)
	r.drawJudgementLines(screen)
	r.drawMeasureMarkers(screen)

	r.drawEffects(screen)
	r.drawNotes(screen)
	// pos := &ui.Point{
	// 	X: playCenterX,
	// 	Y: playCenterY,
	// }
	// size := &ui.Point{
	// 	X: centerComboSize.X,
	// 	Y: centerComboSize.Y,
	// }
	// ui.DrawBorderedFilledRect(screen, pos, size, types.Black, types.White, 0.075)
	r.drawScore(screen)
}

// These are static items we only need to render once
func (r *Play) static(img *ebiten.Image) {
	r.renderBackground(img)
	r.renderProfile(img)
	r.renderSongInfo(img)
	r.renderTracks(img)
}

// TODO: later after tracks and notes
func (r *Play) renderProfile(img *ebiten.Image) *ebiten.Image {
	return img
}
func (r *Play) renderSongInfo(img *ebiten.Image) *ebiten.Image {
	return img
}

func (r *Play) renderBackground(img *ebiten.Image) {
	// TODO: actually make some sort of background?
	img.Fill(types.Black)
}
