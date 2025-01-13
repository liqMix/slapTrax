package ui

import (
	"image/color"

	"github.com/liqmix/slaptrax/internal/types"
	"github.com/liqmix/slaptrax/internal/user"
)

func BrightenColor(c color.RGBA, scale float32) color.RGBA {
	return color.RGBA{
		R: uint8(min(float32(c.R)+(255-float32(c.R))*scale, 255)),
		G: uint8(min(float32(c.G)+(255-float32(c.G))*scale, 255)),
		B: uint8(min(float32(c.B)+(255-float32(c.B))*scale, 255)),
		A: c.A,
	}
}

func ApplyAlphaScale(c color.RGBA, alpha float32) color.RGBA {
	return color.RGBA{
		R: c.R,
		G: c.G,
		B: c.B,
		A: c.A - uint8(float32(c.A)*alpha),
	}
}

func DarkenColor(c color.RGBA, scale float32) color.RGBA {
	return color.RGBA{
		R: uint8(max(float32(c.R)*scale, 0)),
		G: uint8(max(float32(c.G)*scale, 0)),
		B: uint8(max(float32(c.B)*scale, 0)),
		A: c.A,
	}
}

func BorderColor() color.RGBA {
	track := types.TrackLeftTop
	color := track.NoteColor()
	return DarkenColor(color, 0.85)
}

func CornerTrackColor() color.RGBA {
	theme := user.S().NoteColorTheme
	return types.NoteColorTheme(theme).CornerColor()
}

func CenterTrackColor() color.RGBA {
	theme := user.S().NoteColorTheme
	return types.NoteColorTheme(theme).CenterColor()
}
