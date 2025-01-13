package assets

import (
	"embed"
	"path"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//go:embed images/*.png
var imageFS embed.FS

const imageDir = "images"

var loadedImageCache = map[string]*ebiten.Image{}

func GetImage(filename string) *ebiten.Image {
	if img, ok := loadedImageCache[filename]; ok {
		return ebiten.NewImageFromImage(img)
	}
	img, _, err := ebitenutil.NewImageFromFileSystem(imageFS, path.Join(imageDir, filename))
	if err != nil {
		return nil
	}
	loadedImageCache[filename] = img
	return img
}
