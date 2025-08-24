package assets

import "github.com/liqmix/slaptrax/internal/l"

type AssetInit struct {
	Locale string
}

func Init(a AssetInit) {
	InitAudio()
	InitLocales(a.Locale)
	InitSongs()
	
	// Connect the localization function to the l package
	l.GetLocaleString = GetLocaleString
}
