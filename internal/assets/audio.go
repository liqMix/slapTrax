package assets

import (
	"embed"
	"strings"
)

//go:embed sfx/*.mp3 bgm/*.mp3 songs/**/*.mp3
var audioFS embed.FS
var loadedAudio = make(map[string][]byte)

func InitAudio() {}

func GetAudio(path string) ([]byte, error) {
	var err error
	var data, ok = loadedAudio[path]
	if !ok {
		data, err = audioFS.ReadFile(path)
		if err != nil {
			return nil, err
		}
		loadedAudio[path] = data
	}
	return data, nil
}

func AudioExtFromPath(path string) AudioExt {
	for _, ext := range []AudioExt{Wav, Ogg, Mp3} {
		if ext.Is(path) {
			return ext
		}
	}
	panic("Unsupported audio format:" + path)
}

type AudioExt string

const (
	Wav AudioExt = "wav"
	Ogg AudioExt = "ogg"
	Mp3 AudioExt = "mp3"
)

func (a AudioExt) String() string { return string(a) }
func (a AudioExt) Ext() string    { return "." + a.String() }

func (a AudioExt) Is(path string) bool {
	return strings.HasSuffix(path, a.Ext())
}
