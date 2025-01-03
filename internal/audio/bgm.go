package audio

import (
	"path"

	"github.com/liqmix/ebiten-holiday-2024/internal/config"
)

type BGMCode string

const (
	BGMIntro   BGMCode = "intro"
	BGMMenu    BGMCode = "menu"
	BGMResults BGMCode = "results"
)

func (b BGMCode) Path() string {
	return path.Join(config.BGM_DIR, string(b)+".mp3")
}
