package standard

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/cache"
	"github.com/liqmix/ebiten-holiday-2024/internal/locale"
	play "github.com/liqmix/ebiten-holiday-2024/internal/state/play"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/liqmix/ebiten-holiday-2024/internal/ui"
	"github.com/liqmix/ebiten-holiday-2024/internal/user"
)

// type Animation struct {
// 	startTime int64
// }

// func (a *Animation) Draw() float32 {
// 	return 0
// }

// The default renderer for the play state.
type Standard struct {
	state            *play.State
	renderEdgeTracks bool
	// animations map[string]*Animation
}

func (r *Standard) Init(s *play.State) {
	settings := user.Settings()
	displayEdgeTracks := !settings.NoEdgeTracks && s.Chart.HasEdgeTracks()
	SetLayout(displayEdgeTracks)

	r.state = s
	r.renderEdgeTracks = displayEdgeTracks
}

func (r *Standard) Draw(screen *ebiten.Image) {
	r.DrawBackground(screen)
	r.DrawProfile(screen)
	r.DrawSongInfo(screen)
	r.DrawScore(screen)
	r.DrawTracks(screen)
	r.DrawEffects(screen)
}

// TODO: later after tracks and notes
func (r *Standard) DrawProfile(screen *ebiten.Image)  {}
func (r *Standard) DrawSongInfo(screen *ebiten.Image) {}
func (r *Standard) DrawEffects(screen *ebiten.Image)  {}

func (r *Standard) DrawScore(screen *ebiten.Image) {
	// Draw the score at the top of the screen
	score := r.state.Score

	// Draw the score at the top of the screen
	s := user.Settings()
	x := 0.95 * float64(s.RenderWidth)
	y := 0.05 * float64(s.RenderHeight)

	perfectText := fmt.Sprintf(locale.String(types.L_HIT_PERFECT)+": %d", score.Perfect)
	ui.DrawTextRightAt(screen, perfectText, int(x), int(y), 1)

	y += 20
	goodText := fmt.Sprintf(locale.String(types.L_HIT_GOOD)+": %d", score.Good)
	ui.DrawTextRightAt(screen, goodText, int(x), int(y), 1)

	y += 20
	badText := fmt.Sprintf(locale.String(types.L_HIT_BAD)+": %d", score.Bad)
	ui.DrawTextRightAt(screen, badText, int(x), int(y), 1)

	y += 20
	missText := fmt.Sprintf(locale.String(types.L_HIT_MISS)+": %d", score.Miss)
	ui.DrawTextRightAt(screen, missText, int(x), int(y), 1)
}

func (r *Standard) DrawBackground(screen *ebiten.Image) {
	// If we've already created the background, or the render size hasn't changed
	s := user.Settings()

	bg, ok := cache.GetImage("play.background")
	if !ok {
		// Create the background image
		bg = ebiten.NewImage(s.RenderWidth, s.RenderHeight)
		// TODO: actually make some sort of background
		bg.Fill(color.Gray16{0x0000})
		cache.SetImage("play.background", bg)
	}
	screen.DrawImage(bg, nil)
}

func (r *Standard) DrawTracks(screen *ebiten.Image) {
	r.drawMainTracks(screen)
	r.drawEdgeTracks(screen)
	// r.drawMeasureMarkers(screen)
	r.drawNotes(screen)
}
