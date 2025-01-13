package audio

import (
	"path"

	"github.com/liqmix/ebiten-holiday-2024/internal/input"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
)

const (
	bgmDir = "bgm"
	sfxDir = "sfx"
)

type BGMCode string

const (
	BGMTitle BGMCode = "title"
)

func (b BGMCode) Path() string {
	return path.Join(bgmDir, string(b)+".ogg")
}

type SFXCode string

const (
	SFXNone   SFXCode = "none"
	SFXOffset SFXCode = "offset"
	SFXHat    SFXCode = "hat"

	// Menu sounds
	SFXPrev   SFXCode = "prev"
	SFXNext   SFXCode = "next"
	SFXSelect SFXCode = "select"
	SFXBack   SFXCode = "back"

	// Track hit sounds
	SFXTopLeftHit      SFXCode = "hitsoundTopLeft"
	SFXTopRightHit     SFXCode = "hitsoundTopRight"
	SFXBottomLeftHit   SFXCode = "hitsoundBottomLeft"
	SFXBottomRightHit  SFXCode = "hitsoundBottomRight"
	SFXBottomCenterHit SFXCode = "hitsoundBottomCenter"
	SFXTopCenterHit    SFXCode = "hitsoundTopCenter"
)

func (s SFXCode) Path() string {
	return path.Join(sfxDir, string(s)+".ogg")
}

func AllSFXCodes() []SFXCode {
	return []SFXCode{
		SFXPrev,
		SFXNext,
		SFXSelect,
		SFXBack,
		SFXTopLeftHit,
		SFXTopRightHit,
		SFXBottomLeftHit,
		SFXBottomRightHit,
		SFXBottomCenterHit,
		SFXTopCenterHit,
		SFXOffset,
		SFXHat,
	}
}

func TrackSFX(name types.TrackName) SFXCode {
	switch name {
	case types.TrackLeftTop:
		return SFXTopLeftHit
	case types.TrackRightTop:
		return SFXTopRightHit
	case types.TrackLeftBottom:
		return SFXBottomLeftHit
	case types.TrackRightBottom:
		return SFXBottomRightHit
	case types.TrackCenterBottom:
		return SFXBottomCenterHit
	case types.TrackCenterTop:
		return SFXTopCenterHit
	}
	return SFXNone
}

func SFXTrack(s SFXCode) types.TrackName {
	switch s {
	case SFXTopLeftHit:
		return types.TrackLeftTop
	case SFXTopRightHit:
		return types.TrackRightTop
	case SFXBottomLeftHit:
		return types.TrackLeftBottom
	case SFXBottomRightHit:
		return types.TrackRightBottom
	case SFXBottomCenterHit:
		return types.TrackCenterBottom
	case SFXTopCenterHit:
		return types.TrackCenterTop
	}
	return types.TrackUnknown
}
func ActionSFX(a input.Action) SFXCode {
	switch a {
	case input.ActionBack:
		return SFXBack
	case input.ActionSelect:
		return SFXSelect
	case input.ActionUp:
		return SFXPrev
	case input.ActionDown:
		return SFXNext
	case input.ActionLeft:
		return SFXPrev
	case input.ActionRight:
		return SFXNext
	}
	return SFXNone
}
