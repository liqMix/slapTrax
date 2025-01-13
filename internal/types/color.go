package types

import (
	"fmt"
	"image/color"

	"github.com/liqmix/slaptrax/internal/l"
	"github.com/liqmix/slaptrax/internal/logger"
	"github.com/liqmix/slaptrax/internal/user"
)

func ColorFromHex(hex string) color.RGBA {
	defaultC := White.C()

	if hex == "" {
		return defaultC
	}

	var c color.RGBA = color.RGBA{}
	c.A = 0xff
	if hex[0] != '#' {
		logger.Error("Invalid format: %s", hex)
		return defaultC
	}

	hex = hex[1:]
	switch len(hex) {
	case 6:
		_, err := fmt.Sscanf(hex, "%02x%02x%02x", &c.R, &c.G, &c.B)
		if err != nil {
			logger.Error("Invalid hex: %s", hex)
			return defaultC
		}
		return c
	case 8:
		_, err := fmt.Sscanf(hex, "%02x%02x%02x%02x", &c.R, &c.G, &c.B, &c.A)
		if err != nil {
			logger.Error("Invalid hex: %s", hex)
			return defaultC
		}
		return c
	default:
		logger.Error("Invalid hex: %s", hex)
		return defaultC
	}
}

func GameColorFromHex(hex string) GameColor {
	c := ColorFromHex(hex)
	return GameColor(c)
}

func HexFromColor(c color.RGBA) string {
	return fmt.Sprintf("#%02x%02x%02x%02x", c.R, c.G, c.B, c.A)
}

type GameColor color.RGBA

var (
	Black     GameColor = GameColor(color.RGBA{0, 0, 0, 255})
	Gray      GameColor = GameColor(color.RGBA{100, 100, 100, 255})
	White     GameColor = GameColor(color.RGBA{200, 200, 200, 255})
	Orange    GameColor = GameColor(color.RGBA{230, 130, 0, 255})
	Blue      GameColor = GameColor(color.RGBA{50, 50, 255, 255})
	LightBlue GameColor = GameColor(color.RGBA{173, 216, 230, 255})
	Yellow    GameColor = GameColor(color.RGBA{230, 230, 0, 255})
	Red       GameColor = GameColor(color.RGBA{200, 50, 50, 255})
	Green     GameColor = GameColor(color.RGBA{50, 205, 50, 255})
	Purple    GameColor = GameColor(color.RGBA{150, 30, 150, 255})
	Pink      GameColor = GameColor(color.RGBA{200, 145, 145, 255})
)

var AllGameColors = []GameColor{
	// Black,
	Gray,
	White,
	Orange,
	Blue,
	LightBlue,
	Yellow,
	Red,
	Green,
	Purple,
	Pink,
}

func (c GameColor) C() color.RGBA {
	return color.RGBA(c)
}

func (c GameColor) String() string {
	s := "color."
	switch c {
	case Black:
		s += "black"
	case Gray:
		s += "gray"
	case White:
		s += "white"
	case Orange:
		s += "orange"
	case Blue:
		s += "blue"
	case LightBlue:
		s += "lightBlue"
	case Yellow:
		s += "yellow"
	case Red:
		s += "red"
	case Green:
		s += "green"
	case Purple:
		s += "purple"
	case Pink:
		s += "pink"
	default:
		s += "custom"
	}
	return s
}

type NoteColorTheme string

const (
	NoteColorThemeDefault   NoteColorTheme = l.NOTE_COLOR_DEFAULT
	NoteColorThemeMono                     = l.NOTE_COLOR_MONO
	NoteColorThemeDusk                     = l.NOTE_COLOR_DUSK
	NoteColorThemeDawn                     = l.NOTE_COLOR_DAWN
	NoteColorThemeAurora                   = l.NOTE_COLOR_AURORA
	NoteColorThemeArorua                   = l.NOTE_COLOR_ARORUA
	NoteColorThemeHamburger                = l.NOTE_COLOR_HAMBURGER
	NoteColorThemeClassic                  = l.NOTE_COLOR_CLASSIC
	NoteColorThemeCustom                   = l.NOTE_COLOR_CUSTOM
)

func AllNoteColorThemes() []NoteColorTheme {
	return []NoteColorTheme{
		NoteColorThemeDefault,
		NoteColorThemeClassic,
		NoteColorThemeDusk,
		NoteColorThemeDawn,
		NoteColorThemeAurora,
		NoteColorThemeArorua,
		NoteColorThemeMono,
		NoteColorThemeHamburger,
		NoteColorThemeCustom,
	}
}

var themeToColors = map[NoteColorTheme]map[TrackType]color.RGBA{
	NoteColorThemeDefault: {
		TrackTypeCenter: Red.C(),
		TrackTypeCorner: Orange.C(),
	},
	NoteColorThemeClassic: {
		TrackTypeCenter: Yellow.C(),
		TrackTypeCorner: Orange.C(),
	},
	NoteColorThemeDusk: {
		TrackTypeCenter: Blue.C(),
		TrackTypeCorner: Purple.C(),
	},
	NoteColorThemeDawn: {
		TrackTypeCenter: Pink.C(),
		TrackTypeCorner: Yellow.C(),
	},
	NoteColorThemeHamburger: {
		TrackTypeCenter: Red.C(),
		TrackTypeCorner: Yellow.C(),
	},
	NoteColorThemeAurora: {
		TrackTypeCenter: LightBlue.C(),
		TrackTypeCorner: Pink.C(),
	},
	NoteColorThemeArorua: {
		TrackTypeCenter: Pink.C(),
		TrackTypeCorner: LightBlue.C(),
	},
	NoteColorThemeMono: {
		TrackTypeCenter: White.C(),
		TrackTypeCorner: Gray.C(),
	},
}

func (t NoteColorTheme) CenterColor() color.RGBA {
	if t == NoteColorThemeCustom {
		return ColorFromHex(user.S().CenterNoteColor)
	}
	c := themeToColors[t][TrackTypeCenter]
	return c
}
func (t NoteColorTheme) CornerColor() color.RGBA {
	if t == NoteColorThemeCustom {
		return ColorFromHex(user.S().CornerNoteColor)
	}
	c := themeToColors[t][TrackTypeCorner]
	return c
}
