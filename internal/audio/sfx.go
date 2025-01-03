package audio

import (
	"path"

	"github.com/liqmix/ebiten-holiday-2024/internal/config"
)

type SFXCode string

const (
	SFXOffset           SFXCode = "offset"
	SFXNoteHit          SFXCode = "hit"
	SFXMoveSelectorHigh SFXCode = "moveSelectorHigh"
	SFXMoveSelectorLow  SFXCode = "moveSelectorLow"
	SFXSelect           SFXCode = "select"
)

func (s SFXCode) Path() string {
	return path.Join(config.SFX_DIR, string(s)+".mp3")
}

func AllSFX() []SFXCode {
	return []SFXCode{
		SFXOffset,
		SFXNoteHit,
		SFXMoveSelectorHigh,
		SFXMoveSelectorLow,
		SFXSelect,
	}
}
