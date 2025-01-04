package play

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/state"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
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
	state            *state.Play
	vectorCache      *VectorCache
	vectorCollection *ui.VectorCollection
	// animations map[string]*Animation
}

func NewPlayRender(s state.State) *Play {
	p := &Play{state: s.(*state.Play)}
	p.vectorCache = NewVectorCache()
	p.vectorCollection = ui.NewVectorCollection()

	p.BaseRenderer.Init(p.static)
	return p
}

func (r *Play) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	r.BaseRenderer.Draw(screen, opts)
	r.drawMeasureMarkers(screen)

	for _, track := range r.state.Tracks {
		width := judgementWidth
		if track.IsPressed() {
			width *= 2
		}
		cachedPath := r.vectorCache.GetJudgementLinePath(track.Name, track.IsPressed())
		r.vectorCollection.Add(cachedPath.vertices, cachedPath.indices)
	}

	for _, track := range r.state.Tracks {
		if len(track.ActiveNotes) == 0 {
			continue
		}
		for _, note := range track.ActiveNotes {
			// drawNote(screen, note, pts, color)
			path := r.vectorCache.GetNotePath(track.Name, note)
			if path != nil {
				r.vectorCollection.Add(path.vertices, path.indices)
			}
		}
	}
	r.vectorCollection.Draw(screen)

	r.drawEffects(screen)
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

	r.vectorCollection.Clear()
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
