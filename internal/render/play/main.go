package play

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/display"
	"github.com/liqmix/ebiten-holiday-2024/internal/state"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
	"github.com/liqmix/ebiten-holiday-2024/internal/user"
)

// The default renderer for the play state.
type Play struct {
	display.BaseRenderer
	state            *state.Play
	vectorCache      *VectorCache
	vectorCollection *ui.VectorCollection

	hitRecordIdx int
}

func NewPlayRender(s state.State) *Play {
	display.ResetCaches()

	p := &Play{
		state: s.(*state.Play),
	}
	p.vectorCache = NewVectorCache()
	display.AttachCache(p.vectorCache)
	p.vectorCollection = ui.NewVectorCollection()
	p.BaseRenderer.Init(p.static)

	return p
}

func (r *Play) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	r.BaseRenderer.Draw(screen, opts)
	r.drawMeasureMarkers(screen)
	// if user.S.FullScreenLane {
	// 	OscillateWindowOffset(r.state.CurrentTime())
	// }

	// Track vemctors
	for _, track := range r.state.Tracks {
		r.addNotePath(track)
		r.addJudgementPath(track)
		if !user.S.DisableLaneEffects {
			r.addTrackPath(track)
			r.addTrackEffects(track)
		}
	}

	if !user.S.DisableHitEffects {
		r.addHitEffects()
	}
	r.vectorCollection.Draw(screen)

	// Effects and score
	r.drawScore(screen, opts)

	r.vectorCollection.Clear()
}

// These are static items we only need to render once
func (r *Play) static(img *ebiten.Image, opts *ebiten.DrawImageOptions) {
	r.renderBackground(img, opts)
	r.renderProfile(img, opts)
	r.renderSongInfo(img, opts)
}

func (r *Play) renderBackground(img *ebiten.Image, _ *ebiten.DrawImageOptions) {
	// TODO: actually make some sort of background?
	img.Fill(types.Black.C())
}

// TODO: later after tracks and notes
func (r *Play) renderProfile(img *ebiten.Image, opts *ebiten.DrawImageOptions)  {}
func (r *Play) renderSongInfo(img *ebiten.Image, opts *ebiten.DrawImageOptions) {}
