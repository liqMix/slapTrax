package play

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/types"
	"github.com/liqmix/slaptrax/internal/ui"
	"github.com/tinne26/etxt"
)

func (r *Play) drawHeader(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	r.drawScore(screen, opts)

}
func (r *Play) drawStaticHeader(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	r.drawSongDetails(screen, opts)
}

// Right side of header
func (r *Play) drawSongDetails(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	center := headerCenter
	size := ui.Point{
		X: headerWidth,
		Y: headerHeight,
	}
	group := ui.NewUIGroup()
	group.SetDisabled(true)
	group.SetCenter(center)
	group.SetSize(size)

	// Art
	artWidth := headerHeight
	art := ui.NewElement()
	art.SetDisabled(true)
	art.SetSize(ui.Point{
		X: artWidth,
		Y: artWidth,
	})
	art.SetCenter(ui.Point{
		X: headerRight - (artWidth / 4),
		Y: headerCenter.Y + (artWidth / 4),
	})
	art.SetImage(r.state.Song.Art)
	group.Add(art)

	group.Draw(screen, opts)

	// Text
	offset := 0.04
	textOpts := ui.GetDefaultTextOptions()
	textOpts.Align = etxt.Right
	textOpts.Scale = 1.1

	// Song title
	textCenter := &ui.Point{
		X: headerRight - artWidth,
		Y: headerCenter.Y - 0.012,
	}

	// Bold it
	ui.DrawTextAt(screen, r.state.Song.Title, textCenter, textOpts, opts)
	textCenter.X += 0.0005
	ui.DrawTextAt(screen, r.state.Song.Title, textCenter, textOpts, opts)
	textCenter.X -= 0.0005
	textCenter.Y += offset

	textOpts.Scale = 1.0

	// Artist
	textCenter.X += 0.002
	ui.DrawTextAt(screen, r.state.Song.Artist, textCenter, textOpts, opts)
	textCenter.Y += offset

	// Album
	ui.DrawTextAt(screen, r.state.Song.Album, textCenter, textOpts, opts)
	textCenter.Y += offset
}

// Top center of screen
func (r *Play) drawScore(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	score := r.state.Score
	center := &ui.Point{
		X: playCenterX,
		Y: headerCenter.Y + (headerHeight / 4),
	}
	textOpts := &ui.TextOptions{
		Align: etxt.Center,
		Scale: 3.0,
		Color: types.White.C(),
	}

	totalScore := fmt.Sprintf("%d", score.TotalScore)

	ui.DrawTextAt(screen, totalScore, center, textOpts, opts)
}
