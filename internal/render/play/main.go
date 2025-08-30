package play

import (
	"github.com/hajimehoshi/ebiten/v2"
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
		return nil
	}
	
	shaders.InitRenderer()
	ShaderRenderingEnabled = true
	logger.Info("Shader-based note rendering initialized successfully")
	return p
}

func (r *Play) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	r.BaseRenderer.Draw(screen, opts)
	r.drawMeasureMarkers(screen)

	// Track vectors
	for _, track := range r.state.Tracks {
		r.addNotePathShader(track, screen)
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
