package play

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/types"
	"github.com/liqmix/slaptrax/internal/ui"
	"github.com/tinne26/etxt"
)

type HitMessage struct {
	HitRecord   *types.HitRecord
	DisplayTime int64 // When this message was first displayed
}

var (
	fadeOutHitMs  = int64(400) // 400ms for quick fade out
	activeHits    []HitMessage // Stack of currently active hit messages
	lastHitRecord *types.HitRecord // Track the last processed hit record
)

func (r *Play) drawStats(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	r.drawCombo(screen, opts)
	r.drawHitDetails(screen, opts)
}

func (r *Play) drawCombo(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	// Draw the combo text in the center of the combo box ? bit cluttered...
	// Draw the combo above play area
	combo := r.state.Score.Combo

	if combo > 0 {
		comboText := fmt.Sprintf("%d", combo)
		ui.DrawTextAt(
			screen,
			comboText,
			// &headerCenterPoint,
			&comboCenter,
			&ui.TextOptions{
				Align: etxt.Center,
				Scale: 2.0,
				Color: types.White.C(),
			},
			opts,
		)
	}
}
func (r *Play) drawHitDetails(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	currentMs := r.state.CurrentTime()
	lastHit := r.state.Score.GetLastHitRecord()
	
	// Add new hit messages to the stack
	if lastHit != nil && lastHit != lastHitRecord {
		activeHits = append(activeHits, HitMessage{
			HitRecord:   lastHit,
			DisplayTime: currentMs,
		})
		lastHitRecord = lastHit
	}
	
	// Remove expired hit messages and draw active ones
	validHits := make([]HitMessage, 0)
	
	for _, hitMsg := range activeHits {
		timeSinceHit := currentMs - hitMsg.DisplayTime
		
		// Remove expired messages
		if timeSinceHit >= fadeOutHitMs {
			continue // Don't add to validHits (removes it)
		}
		
		// Keep this message
		validHits = append(validHits, hitMsg)
		
		// Calculate fade-out and falling animation
		opacity := 1.0 - float64(timeSinceHit)/float64(fadeOutHitMs)
		fallDistance := float64(timeSinceHit) / float64(fadeOutHitMs) * 0.03
		
		// All messages spawn at same position (no stacking offset)
		
		hitType := hitMsg.HitRecord.HitRating
		c := hitType.Color().C()
		// Don't modify c.A - use ColorScale instead
		
		// Make SLIP and SLOP text smaller than SLAP
		scale := 1.0
		if hitType == types.Slip || hitType == types.Slop {
			scale = 0.8
		}
		
		// Create custom draw options with fade-out opacity
		fadeOpts := &ebiten.DrawImageOptions{}
		if opts != nil {
			// Copy existing options
			*fadeOpts = *opts
		}
		// Apply fade-out opacity through ColorScale
		fadeOpts.ColorScale.ScaleAlpha(float32(opacity))
		
		ui.DrawTextAt(
			screen,
			hitType.String(),
			&ui.Point{
				X: comboCenter.X,
				Y: comboCenter.Y + 0.05 + fallDistance, // All messages start at same position
			},
			&ui.TextOptions{
				Align: etxt.Center,
				Scale: scale,
				Color: c, // Use original color without alpha modification
			},
			fadeOpts, // Use custom options with fade opacity
		)
	}
	
	// Update active hits list
	activeHits = validHits
}
