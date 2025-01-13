package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/ebiten-holiday-2024/internal/external"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/liqmix/ebiten-holiday-2024/internal/user"
	"github.com/tinne26/etxt"
)

type UserProfile struct {
	username string
	rank     float64
	title    string
	position Point

	connected bool
}

func NewUserProfile() *UserProfile {
	topOffset := 0.08
	leftOffset := topOffset / 4
	return &UserProfile{
		position:  Point{X: leftOffset, Y: topOffset},
		connected: external.HasConnection(),
	}
}

func (u *UserProfile) Update() {
	if !u.connected {
		return
	}

	if u.username != user.Current().Username || u.rank != user.Current().Rank {
		u.username = user.Current().Username
		u.rank = user.Current().Rank
		u.title = types.RankTitleFromRank(u.rank).String()
	}
}

func (u *UserProfile) Draw(image *ebiten.Image, opts *ebiten.DrawImageOptions) {
	if !u.connected || external.GetLoginState() == external.StateOffline {
		return
	}

	textOpts := GetDefaultTextOptions()
	textOpts.Align = etxt.Left
	textOpts.Scale = 2

	// Draw the username and rank in top left corner
	if u.username != "" {
		textOpts.Color = CornerTrackColor()
		DrawTextAt(image, u.username, &u.position, textOpts, opts)
		textOpts.Color = CenterTrackColor()
		textOpts.Scale = 1.3
		center := &Point{X: u.position.X, Y: u.position.Y + 0.05}
		DrawTextAt(image, fmt.Sprintf("SF: %.2f", u.rank), center, textOpts, opts)
		textOpts.Color = types.RankTitleFromRank(u.rank).Color()
		DrawTextAt(image, u.title, &Point{X: u.position.X, Y: u.position.Y + 0.1}, textOpts, opts)
	}
}
