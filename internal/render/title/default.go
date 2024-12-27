package title

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/config"
	"github.com/liqmix/ebiten-holiday-2024/internal/state/title"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
	"github.com/liqmix/ebiten-holiday-2024/internal/user"
)

type Default struct {
	state *title.State
}

func (r Default) New(s *title.State) TitleRenderer {
	return &Default{
		state: s,
	}
}
func (r *Default) Draw(screen *ebiten.Image) {
	s := user.Settings()
	centerX := s.RenderWidth / 2
	centerY := s.RenderHeight / 2

	ui.DrawTextAt(screen, r.state.Text, centerX, centerY, config.FONT_SCALE)
	ui.DrawImageAt(screen, r.state.Flag, centerX, centerY-100, 5.0)

}
