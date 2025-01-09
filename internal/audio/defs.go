package audio

import (
	"path"
)

const (
	bgmDir = "bgm"
	sfxDir = "sfx"
)

type BGMCode string

const (
	BGMIntro   BGMCode = "intro"
	BGMMenu    BGMCode = "menu"
	BGMResults BGMCode = "results"
)

func (b BGMCode) Path() string {
	return path.Join(bgmDir, string(b)+".mp3")
}

type SFXCode string

const (
	SFXNone       SFXCode = "none"
	SFXOffset     SFXCode = "offset"
	SFXHat        SFXCode = "hat"
	SFXSelectUp   SFXCode = "selectup"
	SFXSelectDown SFXCode = "selectdown"
	// SFXNoteHit          SFXCode = "hit"
	// SFXSelect     SFXCode = "select"
)

func (s SFXCode) Path() string {
	return path.Join(sfxDir, string(s)+".mp3")
}

func AllSFX() []SFXCode {
	return []SFXCode{
		SFXOffset,
		SFXHat,
		SFXSelectUp,
		SFXSelectDown,
	}
}
