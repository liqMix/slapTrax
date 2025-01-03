package user

import "github.com/liqmix/ebiten-holiday-2024/internal/assets"

type User struct {
	Name     string
	Settings *UserSettings
}

var u = &User{
	Name:     "funyarinpa",
	Settings: NewUserSettings(),
}

func Init() {
	if u == nil {
		u = &User{
			Name:     "funyarinpa",
			Settings: NewUserSettings(),
		}
		u.Settings.Apply()
	}

}

func S() *UserSettings {
	return u.Settings
}

func Volume() *assets.Volume {
	return &assets.Volume{
		Bgm:  u.Settings.Audio.BGMVolume,
		Sfx:  u.Settings.Audio.SFXVolume,
		Song: u.Settings.Audio.SongVolume,
	}
}
