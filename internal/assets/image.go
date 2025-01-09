package assets

import (
	"embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//go:embed images/*
var imageFS embed.FS

func GetImage(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFileSystem(imageFS, path)
	if err != nil {
		return nil
	}
	return img
}
