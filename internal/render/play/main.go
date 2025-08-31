package play

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/cache"
	"github.com/liqmix/slaptrax/internal/display"
	"github.com/liqmix/slaptrax/internal/logger"
	"github.com/liqmix/slaptrax/internal/render/shaders"
	"github.com/liqmix/slaptrax/internal/state"
	"github.com/liqmix/slaptrax/internal/ui"
	"github.com/liqmix/slaptrax/internal/user"
)


// The default renderer for the play state.
type Play struct {
	display.BaseRenderer
	state            *state.Play
	vectorCollection *ui.VectorCollection

	hitRecordIdx int
	lastCacheCheck bool
	isFrozen     bool // Flag to indicate if we're rendering for frozen state
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
	
	// Initialize play area layouts
	ReinitLayouts()
	p.lastCacheCheck = true // Assume cache exists initially
	
	return p
}

func (r *Play) shouldReinitLayouts() bool {
	// Simple heuristic: check if cache is empty (indicating it was cleared)
	// This is a simplified approach that assumes settings changes clear the cache
	_, cacheExists := cache.Image.Get("canvas")
	if r.lastCacheCheck && !cacheExists {
		r.lastCacheCheck = false
		return true
	}
	r.lastCacheCheck = cacheExists
	return false
}

// SetFrozen sets the frozen state to skip animated effects when cached
func (r *Play) SetFrozen(frozen bool) {
	r.isFrozen = frozen
}

func (r *Play) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	// Check if layouts need to be reinitialized
	if cache.LayoutReinitRequested {
		ReinitLayouts()
		cache.LayoutReinitRequested = false
		logger.Info("Layouts reinitialized after settings change")
	} else if r.shouldReinitLayouts() {
		ReinitLayouts()
		logger.Info("Layouts reinitialized after cache cleared")
	}
	
	r.BaseRenderer.Draw(screen, opts)
	
	// Skip animated shader effects when frozen (for pause caching)
	if !r.isFrozen {
		// Render tunnel background first (underneath everything)
		if !user.S().DisableLaneEffects && shaders.LaneRendererInstance != nil {
			shaders.LaneRendererInstance.RenderTunnelBackground(screen, &playCenterPoint, float32(playLeft), float32(playRight), float32(playTop), float32(playBottom))
		}
		
		// Shader-based lane rendering
		for _, track := range r.state.Tracks {
			// Render lane background using shader
			if !user.S().DisableLaneEffects && shaders.LaneRendererInstance != nil {
				trackPoints := notePoints[track.Name]
				isActive := track.IsPressed()
				shaders.LaneRendererInstance.RenderLaneBackground(screen, track.Name, trackPoints, &playCenterPoint, isActive)
			}
		}
		
		// Render shader-based measure markers
		r.drawMeasureMarkersShader(screen)
	}

	// Track vectors
	for _, track := range r.state.Tracks {
		r.addNotePathShader(track, screen)
		r.addJudgementPath(track)
		if !user.S().DisableLaneEffects {
			// Only keep track effects (particles, etc.), not the track path itself
			r.addTrackEffects(track)
		}
	}

	if !user.S().DisableHitEffects {
		r.addHitEffects()
	}

	r.drawStats(screen, opts)
	r.vectorCollection.Draw(screen)
	r.vectorCollection.Clear()

}

// These are static items we only need to render once
func (r *Play) static(img *ebiten.Image, opts *ebiten.DrawImageOptions) {
	r.renderBackground(img, opts)
	// Header rendering moved to unified system in game.go
}

func (r *Play) renderBackground(img *ebiten.Image, _ *ebiten.DrawImageOptions) {
	// Background is already rendered by game.go, no need to fill
	// The tunnel background will provide proper occlusion where needed
}
