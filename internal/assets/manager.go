package assets

type AssetInit struct {
	Locale string
}

func Init(a AssetInit) {
	InitAudio()
	InitLocales(a.Locale)
	InitSongs()
}
