package def

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/state/play"
)

type PlayRenderer interface {
	Init(s *play.State)
	Draw(screen *ebiten.Image)
	DrawBackground(screen *ebiten.Image)
	DrawProfile(screen *ebiten.Image)
	DrawSongInfo(screen *ebiten.Image)
	DrawScore(screen *ebiten.Image)
	DrawTracks(screen *ebiten.Image)
	DrawEffects(screen *ebiten.Image)
}
