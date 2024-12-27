package play

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/song"
	play "github.com/liqmix/ebiten-holiday-2024/internal/state/play"
	"github.com/liqmix/ebiten-holiday-2024/internal/user"
)

// The default renderer for the play state.
type Default struct {
	state            *play.State
	background       *ebiten.Image
	trackBackgrounds map[song.TrackName]*ebiten.Image
}

func (r Default) New(s *play.State) PlayRenderer {
	return &Default{
		state:            s,
		trackBackgrounds: make(map[song.TrackName]*ebiten.Image),
	}
}

func (r *Default) Draw(screen *ebiten.Image) {
	r.drawBackground(screen)
	r.drawTracks(screen)
	r.drawProfile(screen)
	r.drawSongInfo(screen)
	r.drawScore(screen)
}

// TODO: later after tracks and notes
func (r *Default) drawProfile(screen *ebiten.Image)  {}
func (r *Default) drawSongInfo(screen *ebiten.Image) {}
func (r *Default) drawScore(screen *ebiten.Image)    {}

func (r *Default) drawBackground(screen *ebiten.Image) {
	// If we've already created the background, or the render size hasn't changed
	s := user.Settings()

	if r.background != nil {
		x, y := r.background.Bounds().Dx(), r.background.Bounds().Dy()
		if x == s.RenderWidth && y == s.RenderWidth {
			screen.DrawImage(r.background, nil)
			return
		}
	}

	// Create the background image
	r.background = ebiten.NewImage(s.RenderWidth, s.RenderHeight)

	// TODO: create the actual background

	// Create the background for the main tracks
	//// These tracks will be curved lanes around a center point.
	//// The top tracks with nearly touch on the top of the screen,
	//// however the bottom tracks will have a gap for the center track to fit in.

	// Create the background for the center track
	//// This track will be a straight lane between the bottom left and bottom right tracks

	// Create the background for the edge tracks

	// draw the background
	screen.DrawImage(r.background, nil)
}

func (r *Default) drawTracks(screen *ebiten.Image) {
	for _, t := range r.state.Tracks {
		if r.trackBackgrounds[t.Name] != nil {

		}
		switch t.Name {
		case song.Center:
			r.drawCenterTrack(screen, &t)
			return
		case song.EdgeTop:
		case song.EdgeTap1:
		case song.EdgeTap2:
		case song.EdgeTap3:
			r.drawEdgeTrack(screen, &t)
			return
		default:
			r.drawMainTrack(screen, &t)
		}
	}
}

// draws a main track, one of the four rounded tracks
func (r *Default) drawMainTrack(screen *ebiten.Image, t *song.Track) {
}

// draws the center track, a straight lane
func (r *Default) drawCenterTrack(screen *ebiten.Image, t *song.Track) {
}

func (r *Default) drawEdgeTrack(screen *ebiten.Image, t *song.Track) {
}
