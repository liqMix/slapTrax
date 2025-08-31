package play

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/types"
	"github.com/liqmix/slaptrax/internal/ui"
	"github.com/tinne26/etxt"
)

var (
	fadeOutHitMs       = int64(500)
	lastDisplayedHitMs int64
	prevHit            *types.HitRecord
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
	// Draw the last hit text
	lastHit := r.state.Score.GetLastHitRecord()
	if lastHit != nil {
		// currentMs := r.state.CurrentTime()
		// if lastHit != prevHit {
		// 	lastDisplayedHitMs = currentMs
		// }
		// opacity := 1.0 - float64(currentMs-lastDisplayedHitMs)/float64(fadeOutHitMs)
		// c.A = uint8(opacity * 255)
		// prevHit = lastHit

		hitType := lastHit.HitRating
		c := hitType.Color().C()

		// Make SLIP and SLOP text smaller than SLAP
		scale := 1.0
		if hitType == types.Slip || hitType == types.Slop {
			scale = 0.8
		}

		ui.DrawTextAt(
			screen,
			hitType.String(),
			&ui.Point{
				X: comboCenter.X,
				Y: comboCenter.Y + 0.05,
			},
			&ui.TextOptions{
				Align: etxt.Center,
				Scale: scale,
				Color: c,
			},
			opts,
		)
	}
}
