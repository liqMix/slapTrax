package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/liqmix/slaptrax/internal/assets"
	"github.com/liqmix/slaptrax/internal/external"
	"github.com/liqmix/slaptrax/internal/types"
	"github.com/liqmix/slaptrax/internal/user"
	"github.com/tinne26/etxt"
)

type UserProfile struct {
	username string
	rank     float64
	title    string
	position Point
	avatar   *ebiten.Image

	connected bool
	
	// Pre-created UI elements to avoid per-frame allocation
	avatarElement *Element
	avatarGroup   *UIGroup
}

func NewUserProfile() *UserProfile {
	// Header constants with equal margins from all edges
	headerCenterY := 0.1   // headerTop + headerHeight/2 = 0.0 + 0.2/2
	equalMargin := 0.025   // Equal distance from top, sides, and border
	headerHeight := 0.2    // 1/5th of screen height
	avatarSize := headerHeight * 0.4 // Use 40% of header height (same as album art)
	
	// Pre-create UI elements once
	art := NewElement()
	art.SetDisabled(true)
	art.SetSize(Point{
		X: avatarSize,
		Y: avatarSize,
	})
	art.SetImage(assets.GetImage("default_art.png"))
	
	avatarGroup := NewUIGroup()
	avatarGroup.SetDisabled(true)
	avatarGroup.SetSize(Point{X: avatarSize, Y: avatarSize})
	avatarGroup.Add(art)
	
	return &UserProfile{
		position:      Point{X: equalMargin, Y: headerCenterY},
		avatar:        assets.GetImage("default_art.png"), // Use default album art as avatar
		connected:     external.HasConnection(),
		avatarElement: art,
		avatarGroup:   avatarGroup,
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
		// Update avatar element when data changes (if needed)
		u.avatarElement.SetImage(u.avatar)
	}
}

func (u *UserProfile) Draw(image *ebiten.Image, opts *ebiten.DrawImageOptions) {
	// Header constants with equal margins from all edges (EXACT same as song details)
	headerHeight := 0.2   // 1/5th of screen height
	headerCenterY := 0.1  // headerTop + headerHeight/2
	equalMargin := 0.025  // Equal distance from top, sides, and border
	
	// Avatar - positioned on left side of header (mirror of album art)
	avatarSize := headerHeight * 0.4 // Use 40% of header height (same as album art)
	
	// Update position of pre-created elements (no allocation)
	u.avatarElement.SetCenter(Point{
		X: equalMargin + (avatarSize / 2), // Left side (mirror of right side)
		Y: headerCenterY,                  // Same Y as album art
	})
	
	u.avatarGroup.SetCenter(Point{X: equalMargin + (avatarSize / 2), Y: headerCenterY})
	u.avatarGroup.Draw(image, opts)

	// Text - positioned to the right of the avatar (mirror of song details)
	textMargin := 0.02
	textOpts := GetDefaultTextOptions()
	textOpts.Align = etxt.Left // Keep left-aligned for user details
	textOpts.Scale = 1.3       // Match song title scale

	// Username (mirror of song title)
	textCenter := &Point{
		X: equalMargin + avatarSize + textMargin,
		Y: headerCenterY - 0.03, // Adjusted for smaller text
	}

	if !u.connected || external.GetLoginState() == external.StateOffline {
		// Show "slapGuest" when not logged in
		textOpts.Color = CornerTrackColor()
		DrawTextAt(image, "slapGuest", textCenter, textOpts, opts)
		return
	}

	// Draw the username, rank, and title (mirroring song details structure)
	if u.username != "" {
		// Username (bold, like song title)
		textOpts.Color = CornerTrackColor()
		DrawTextAt(image, u.username, textCenter, textOpts, opts)
		textCenter.X += 0.001
		DrawTextAt(image, u.username, textCenter, textOpts, opts) // Bold effect
		textCenter.X -= 0.001

		textOpts.Scale = 1.0 // Match artist scale
		textCenter.Y += 0.04

		// Rank (like artist)
		textOpts.Color = CenterTrackColor()
		DrawTextAt(image, fmt.Sprintf("SF: %.2f", u.rank), textCenter, textOpts, opts)
		textCenter.Y += 0.03

		// Title (like album)
		textOpts.Scale = 0.9 // Match album scale
		textOpts.Color = types.RankTitleFromRank(u.rank).Color()
		DrawTextAt(image, u.title, textCenter, textOpts, opts)
	}
}
