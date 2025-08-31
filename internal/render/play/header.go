package play

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/types"
	"github.com/liqmix/slaptrax/internal/ui"
	"github.com/tinne26/etxt"
)

func (r *Play) drawHeader(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	r.drawScore(screen, opts)

}
func (r *Play) drawStaticHeader(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	// Draw header background panel first (before content)
	r.drawHeaderBackground(screen, opts)
	r.drawSongDetails(screen, opts)
}

// Draw header background panel with bottom border only
func (r *Play) drawHeaderBackground(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	// Header panel dimensions (full width, 1/5 height)
	center := headerCenter
	size := ui.Point{
		X: headerWidth,
		Y: headerHeight,
	}
	
	// Draw main background panel
	ui.DrawFilledRect(screen, &center, &size, color.RGBA{R: 0, G: 0, B: 0, A: 120})
	
	// Draw bottom border only
	borderHeight := 0.003 // Thin border
	borderCenter := ui.Point{
		X: center.X,
		Y: headerBottom - (borderHeight / 2), // Position at bottom edge
	}
	borderSize := ui.Point{
		X: headerWidth,
		Y: borderHeight,
	}
	
	// Use corner track color for the border
	ui.DrawFilledRect(screen, &borderCenter, &borderSize, ui.CornerTrackColor())
}

// Right side of header
func (r *Play) drawSongDetails(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	// Art - positioned in right side of header
	artSize := headerHeight * 0.4 // Use 40% of header height
	rightMargin := 0.05 // Margin from right edge
	
	art := ui.NewElement()
	art.SetDisabled(true)
	art.SetSize(ui.Point{
		X: artSize,
		Y: artSize,
	})
	art.SetCenter(ui.Point{
		X: headerRight - rightMargin - (artSize / 2),
		Y: headerCenter.Y,
	})
	art.SetImage(r.state.Song.Art)
	
	// Draw art directly
	artGroup := ui.NewUIGroup()
	artGroup.SetDisabled(true)
	artGroup.SetCenter(ui.Point{X: headerRight - rightMargin - (artSize / 2), Y: headerCenter.Y})
	artGroup.SetSize(ui.Point{X: artSize, Y: artSize})
	artGroup.Add(art)
	artGroup.Draw(screen, opts)

	// Text - positioned to the left of the art
	textMargin := 0.02
	textOpts := ui.GetDefaultTextOptions()
	textOpts.Align = etxt.Right
	textOpts.Scale = 1.3 // Reduced from 1.8

	// Song title
	textCenter := &ui.Point{
		X: headerRight - rightMargin - artSize - textMargin,
		Y: headerCenter.Y - 0.03, // Adjusted for smaller text
	}

	// Bold title
	ui.DrawTextAt(screen, r.state.Song.Title, textCenter, textOpts, opts)
	textCenter.X += 0.001
	ui.DrawTextAt(screen, r.state.Song.Title, textCenter, textOpts, opts)
	textCenter.X -= 0.001

	textOpts.Scale = 1.0 // Reduced from 1.4
	textCenter.Y += 0.04

	// Artist
	ui.DrawTextAt(screen, r.state.Song.Artist, textCenter, textOpts, opts)
	textCenter.Y += 0.03

	// Album
	textOpts.Scale = 0.9 // Reduced from 1.2
	ui.DrawTextAt(screen, r.state.Song.Album, textCenter, textOpts, opts)
}

// Center of header - score display
func (r *Play) drawScore(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	score := r.state.Score
	center := &ui.Point{
		X: 0.5, // Center of screen
		Y: headerCenter.Y,
	}
	textOpts := &ui.TextOptions{
		Align: etxt.Center,
		Scale: 4.0, // Larger scale for the bigger header
		Color: types.White.C(),
	}

	totalScore := fmt.Sprintf("%d", score.TotalScore)

	ui.DrawTextAt(screen, totalScore, center, textOpts, opts)
}
