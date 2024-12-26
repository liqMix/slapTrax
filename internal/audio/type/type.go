package audiotype

type AudioType string

const (
	SFX         AudioType = "audio.sfx"
	BGM         AudioType = "audio.music"
	Song        AudioType = "audio.song"
	SongPreview AudioType = "audio.songPreview"
)
