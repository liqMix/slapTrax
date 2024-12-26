package play

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/song"
)

type LeftBehind struct {
	Renderer
}

func (r *LeftBehind) Draw(screen *ebiten.Image) {
	for _, t := range r.state.Tracks {
		r.drawTrack(screen, &t)
	}
}

func (r *LeftBehind) drawTrack(screen *ebiten.Image, t *song.Track) {
	for _, n := range t.Notes {
		r.drawNote(screen, &n)
	}
}

func (r *LeftBehind) drawNote(screen *ebiten.Image, n *song.Note) {
}
