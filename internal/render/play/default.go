package play

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/song"
	play "github.com/liqmix/ebiten-holiday-2024/internal/state/play"
)

type Default struct{}

func (r *Default) Draw(screen *ebiten.Image, state *play.State) {
	for _, t := range state.Tracks {
		r.drawTrack(screen, &t)
	}
}

func (r *Default) drawTrack(screen *ebiten.Image, t *song.Track) {
	for _, n := range t.Notes {
		r.drawNote(screen, &n)
	}
}

func (r *Default) drawNote(screen *ebiten.Image, n *song.Note) {
}
