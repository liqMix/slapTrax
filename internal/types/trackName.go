package types

import (
	"image/color"

	"github.com/liqmix/ebiten-holiday-2024/internal/input"
	"github.com/liqmix/ebiten-holiday-2024/internal/user"
)

type TrackName int

// Order is critical here
const (
	TrackLeftBottom TrackName = iota
	TrackLeftTop
	TrackCenterBottom
	TrackCenterTop
	TrackRightBottom
	TrackRightTop
	TrackUnknown
)

func TrackNames() []TrackName {
	return []TrackName{
		TrackLeftBottom,
		TrackLeftTop,
		TrackRightBottom,
		TrackRightTop,
		TrackCenterBottom,
		TrackCenterTop,
	}
}

func (t TrackName) String() string {
	switch t {
	case TrackLeftBottom:
		return "LeftBottom"
	case TrackLeftTop:
		return "LeftTop"
	case TrackRightBottom:
		return "RightBottom"
	case TrackRightTop:
		return "RightTop"
	case TrackCenterBottom:
		return "CenterBottom"
	case TrackCenterTop:
		return "CenterTop"
	}
	return "Unknown"
}

func (t TrackName) Action() input.Action {
	switch t {
	case TrackLeftBottom:
		return input.ActionLeftBottom
	case TrackLeftTop:
		return input.ActionLeftTop
	case TrackRightBottom:
		return input.ActionRightBottom
	case TrackRightTop:
		return input.ActionRightTop
	case TrackCenterBottom:
		return input.ActionCenterBottom
	case TrackCenterTop:
		return input.ActionCenterTop
	}
	return input.ActionUnknown
}

func (t TrackName) NoteColor() color.RGBA {
	return TrackTypeFromName(t).Color()
}

type TrackType int

const (
	TrackTypeCenter TrackType = iota
	TrackTypeCorner
)

func TrackTypeFromName(n TrackName) TrackType {
	switch n {
	case TrackCenterBottom, TrackCenterTop:
		return TrackTypeCenter
	}
	return TrackTypeCorner
}

func (t TrackType) Color() color.RGBA {
	theme := NoteColorTheme(user.S().NoteColorTheme)
	switch t {
	case TrackTypeCenter:
		return theme.CenterColor()
	case TrackTypeCorner:
		return theme.CornerColor()
	}
	return White.C()
}
