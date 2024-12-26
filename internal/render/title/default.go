package title

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/config"
	"github.com/liqmix/ebiten-holiday-2024/internal/state/title"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
)

type Default struct{}

func (r *Default) Draw(screen *ebiten.Image, state *title.State) {
	// Draw everything to canvas at base resolution
	centerX := config.CANVAS_WIDTH / 2
	centerY := config.CANVAS_HEIGHT / 2

	ui.DrawTextAt(screen, state.Text, centerX, centerY, config.FONT_SCALE)
	ui.DrawImageAt(screen, state.Flag, centerX, centerY-100, 5.0)

}
