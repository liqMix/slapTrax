package play

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/cache"
	"github.com/liqmix/slaptrax/internal/display"
	"github.com/liqmix/slaptrax/internal/logger"
	"github.com/liqmix/slaptrax/internal/render/shaders"
	"github.com/liqmix/slaptrax/internal/state"
	"github.com/liqmix/slaptrax/internal/types"
	"github.com/liqmix/slaptrax/internal/ui"
	"github.com/liqmix/slaptrax/internal/user"
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

	// Initialize shader system
	if err := shaders.InitManager(); err != nil {
		logger.Error("Failed to initialize shader manager: %v", err)
		logger.Info("Falling back to vertex-based rendering")
		ShaderRenderingEnabled = false
	} else {
		shaders.InitRenderer()
		// Enable shader rendering by default for testing
		ShaderRenderingEnabled = true
		logger.Info("Shader-based note rendering initialized successfully")
	}

	// Only initialize vector cache system if shaders failed
	if !ShaderRenderingEnabled {
		logger.Info("Initializing vector cache system")
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
	} else {
		logger.Info("Skipping vector cache initialization - using shaders")
	}
	return p
}

func (r *Play) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	if cache.Path.IsBuilding() {
		return
	}

	r.BaseRenderer.Draw(screen, opts)
	r.drawMeasureMarkers(screen)

	// Track vectors
	for _, track := range r.state.Tracks {
		if ShaderRenderingEnabled {
			r.addNotePathShader(track, screen)
		} else {
			r.addNotePath(track)
		}
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
