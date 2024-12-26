package types

import "github.com/liqmix/ebiten-holiday-2024/internal/l"

type Theme string

const (
	ThemeDefault    Theme = "theme.default"
	ThemeLeftBehind Theme = "theme.leftbehind"
)

func (t Theme) String() string {
	return l.String(string(t))
}
