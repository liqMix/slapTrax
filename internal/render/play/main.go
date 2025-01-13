package play

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/cache"
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
	vectorCollection *ui.VectorCollection

	hitRecordIdx int
}

func NewPlayRender(s state.State) *Play {
	p := &Play{
		state:            s.(*state.Play),
		vectorCollection: ui.NewVectorCollection(),
	}
	p.BaseRenderer.Init(p.static)

	// awful but idc
	cache.Path.RemoveCbs()
	cb := func() {
		go func() {
			if cache.Path.IsBuilding() {
				return
			}
			cache.Path.SetIsBuilding(true)
			RebuildVectorCache()
			cache.Path.SetIsBuilding(false)
		}()
	}
	cache.Path.AddCb(&cb)
	cache.Path.Clear()
	return p
}

func (r *Play) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	if cache.Path.IsBuilding() {
		return
	}

	r.BaseRenderer.Draw(screen, opts)
	r.drawMeasureMarkers(screen)

	// Track vemctors
	for _, track := range r.state.Tracks {
		r.addNotePath(track)
		r.addJudgementPath(track)
		if !user.S().DisableLaneEffects {
			r.addTrackPath(track)
			r.addTrackEffects(track)
		}
	}

	if !user.S().DisableHitEffects {
		r.addHitEffects()
	}

	r.drawHeader(screen, opts)
	r.drawStats(screen, opts)
	r.vectorCollection.Draw(screen)
	r.vectorCollection.Clear()

}

// These are static items we only need to render once
func (r *Play) static(img *ebiten.Image, opts *ebiten.DrawImageOptions) {
	r.renderBackground(img, opts)
	r.drawStaticHeader(img, opts)
}

func (r *Play) renderBackground(img *ebiten.Image, _ *ebiten.DrawImageOptions) {
	// TODO: actually make some sort of background?
	img.Fill(types.Black.C())
}
