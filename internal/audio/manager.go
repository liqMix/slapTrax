package audio

import (
	"bytes"
	"io"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/liqmix/slaptrax/internal/assets"
	"github.com/liqmix/slaptrax/internal/config"
	"github.com/liqmix/slaptrax/internal/logger"
	"github.com/liqmix/slaptrax/internal/types"
	"github.com/solarlune/resound"
	"github.com/solarlune/resound/effects"
)

type channels struct {
	bgm            *resound.DSPChannel
	song           *resound.DSPChannel
	previewPrev    *resound.DSPChannel
	previewCurrent *resound.DSPChannel
}

type players struct {
	bgm            *resound.Player
	song           *resound.Player
	previewCurrent *resound.Player
	previewPrev    *resound.Player
}

type previewStatus struct {
	start      int64
	end        int64
	restarting bool
}

type audioMan struct {
	audioProperties resound.AudioProperties
	sfxCache        *SFXPlayerCache
	players         *players
	channels        *channels

	previewStatus *previewStatus
	volume        *Volume
}

type Volume struct {
	BGM  float64
	SFX  float64
	Song float64
}

var manager *audioMan

const sampleRate = 48000

func InitAudioManager(v *Volume) {
	logger.Info("Initializing audio manager...")
	audio.NewContext(sampleRate)
	audioProperties := resound.NewAudioProperties()

	sfxCache, err := NewSFXPlayerCache(audioProperties, v.SFX)
	if err != nil {
		panic("Error initializing SFX cache: " + err.Error())
	}

	manager = &audioMan{
		audioProperties: audioProperties,

		sfxCache: sfxCache,
		channels: newChannels(v),
		players:  &players{},
		volume:   v,
	}
}

func newChannels(v *Volume) *channels {
	bgm := resound.NewDSPChannel()
	song := resound.NewDSPChannel()
	previewPrev := resound.NewDSPChannel()
	previewCurrent := resound.NewDSPChannel()

	bgm.AddEffect("volume", effects.NewVolume().SetStrength(v.BGM))
	song.AddEffect("volume", effects.NewVolume().SetStrength(v.Song))
	previewPrev.AddEffect("volume", effects.NewVolume().SetStrength(v.Song))
	previewCurrent.AddEffect("volume", effects.NewVolume().SetStrength(v.Song))

	return &channels{
		bgm:            bgm,
		song:           song,
		previewPrev:    previewPrev,
		previewCurrent: previewCurrent,
	}
}

func getVolumeEffect(c *resound.DSPChannel) *effects.Volume {
	if effect, ok := c.Effects["volume"]; !ok {
		return nil
	} else {
		return effect.(*effects.Volume)
	}
}

func SetBGMVolume(v float64) {
	manager.volume.BGM = v
	if manager.channels.bgm != nil {
		getVolumeEffect(manager.channels.bgm).SetStrength(v)
	}
}

func SetSFXVolume(v float64) {
	manager.volume.SFX = v
	if manager.sfxCache != nil {
		manager.sfxCache.SetVolume(v)
	}
}

func SetSongVolume(v float64) {
	manager.volume.Song = v
	if manager.channels.song != nil {
		getVolumeEffect(manager.channels.song).SetStrength(v)
	}
	if manager.channels.previewCurrent != nil {
		getVolumeEffect(manager.channels.previewCurrent).SetStrength(v)
	}
	if manager.channels.previewPrev != nil {
		getVolumeEffect(manager.channels.previewPrev).SetStrength(v)
	}
}

func getStream(path string) (io.ReadSeeker, error) {
	b, err := assets.GetAudio(path)
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(b)
	ext := assets.AudioExtFromPath(path)
	var stream io.ReadSeeker
	switch ext {
	case assets.Wav:
		stream, err = wav.DecodeWithSampleRate(sampleRate, reader)
	case assets.Ogg:
		stream, err = vorbis.DecodeWithSampleRate(sampleRate, reader)
	case assets.Mp3:
		stream, err = mp3.DecodeWithSampleRate(sampleRate, reader)
	}
	if err != nil {
		return nil, err
	}

	return stream, nil
}

func getPlayer(path string, c *resound.DSPChannel) (*resound.Player, error) {
	stream, err := getStream(path)
	if err != nil {
		return nil, err
	}

	volume := getVolumeEffect(c)
	if volume != nil {
		prop, err := manager.audioProperties.Get(path).Analyze(stream, 0)
		if err != nil {
			logger.Error("Error analyzing song (%s): %v", path, err)
		} else {
			volume.SetNormalizationFactor(prop.Normalization)
		}
	}

	player, err := resound.NewPlayer(stream)
	if err != nil {
		return nil, err
	}

	player.SetDSPChannel(c)
	return player, nil
}

func PlaySFX(sfx SFXCode) {
	if manager.volume.SFX == 0 {
		return
	}
	manager.sfxCache.PlaySound(sfx)
}

func PlayTrackSFX(trackName types.TrackName) {
	if manager.volume.SFX == 0 {
		return
	}
	manager.sfxCache.PlayTrackSound(trackName)
}

func GetBGM() *resound.Player {
	return manager.players.bgm
}

func PlayBGM(bgmCode BGMCode) {
	var volume *effects.Volume

	if manager.players.bgm == nil {
		path := bgmCode.Path()
		stream, err := getStream(path)
		if err != nil {
			logger.Error("Error getting stream for BGM: %s", path)
			return
		}
		volume = getVolumeEffect(manager.channels.bgm)

		if volume != nil {
			prop, err := manager.audioProperties.Get(path).Analyze(stream, 0)
			if err != nil {
				logger.Error("Error analyzing song (%s): %v", path, err)
			} else {
				volume.SetNormalizationFactor(prop.Normalization)
			}
			volume.StartFade(0, manager.volume.BGM, config.AUDIO_FADE_S*2)
		}
		length, err := stream.Seek(0, io.SeekEnd)
		if err != nil {
			logger.Error("Error getting length of BGM: %s", path)
			return
		}
		stream.Seek(0, io.SeekStart)

		loop := audio.NewInfiniteLoop(stream, length)
		player, err := resound.NewPlayer(loop)
		if err != nil {
			logger.Error("Error creating player for BGM: %s", path)
			return
		}

		player.SetDSPChannel(manager.channels.bgm)
		manager.players.bgm = player
	}

	if !manager.players.bgm.IsPlaying() {
		if volume == nil {
			volume := getVolumeEffect(manager.channels.bgm)
			if volume != nil {
				volume.StartFade(0, manager.volume.BGM, config.AUDIO_FADE_S*2)
			}
		}
		manager.players.bgm.Play()
	}
}

func GetBGMPositionMS() int64 {
	if manager.players.bgm != nil {
		return manager.players.bgm.Position().Milliseconds()
	}
	return 0
}

func FadeInBGM() {
	if manager.players.bgm != nil {
		effect := getVolumeEffect(manager.channels.bgm)
		if effect != nil {
			effect.StartFade(0, manager.volume.BGM, config.AUDIO_FADE_S)
		}
		if !manager.players.bgm.IsPlaying() {
			manager.players.bgm.Play()
		}
	}
}

func FadeOutBGM() {
	if manager.players.bgm != nil {
		effect := getVolumeEffect(manager.channels.bgm)
		if effect != nil {
			effect.StartFade(-1, 0, config.AUDIO_FADE_S)
		}
	}
}

func fadeCurrentPreview(start float64, end float64) {
	currentChannel := manager.channels.previewCurrent
	if currentChannel == nil {
		return
	}

	effect := getVolumeEffect(currentChannel)
	if effect == nil {
		return
	}
	effect.StartFade(start, end, config.AUDIO_FADE_S)
}

func PlaySongPreview(s *types.Song) {
	startPosition := s.PreviewStart
	endPosition := s.PreviewStart + config.SONG_PREVIEW_LENGTH
	player, err := getPlayer(s.AudioPath, manager.channels.previewCurrent)
	if err != nil {
		logger.Error("Error playing song preview: %v\n", err)
		return
	}

	// If there is a current song preview, move it to previous and fade out
	current := manager.players.previewCurrent
	if current != nil {
		// just stop it for now
		current.Close()

		//// TODO: fade out
		// if manager.players.songPreview.previous != nil {
		// 	manager.players.songPreview.previous.Close()
		// }
		// previousChannel := manager.channels.previewPrev
		// current.SetDSPChannel(previousChannel)
		// previousVolume := previousChannel.Effects["volume"].(*effects.Volume)
		// previousVolume.StartFade(current.Volume(), 0, config.SONG_PREVIEW_FADE/1000)
		// manager.players.songPreview.previous = current
	}

	// Fade in the new song preview
	player.SetDSPChannel(manager.channels.previewCurrent)
	manager.players.previewCurrent = player
	player.SetPosition(time.Duration(startPosition) * time.Millisecond)
	fadeCurrentPreview(-1, manager.volume.Song)
	manager.previewStatus = &previewStatus{
		start: startPosition,
		end:   endPosition,
	}
	player.Play()
}

func GetSongPreviewPositionMS() int64 {
	if manager.players.previewCurrent != nil {
		return manager.players.previewCurrent.Position().Milliseconds()
	}
	return 0
}

// Song
func InitSong(s *types.Song) {
	player, err := getPlayer(s.AudioPath, manager.channels.song)
	if err != nil {
		logger.Error("Error initializing song (%s): %v", s.AudioPath, err)
		return
	}
	manager.players.song = player
}

func CurrentSongPositionMS() int64 {
	if manager.players.song != nil {
		return manager.players.song.Position().Milliseconds()
	}
	return 0
}

func IsSongPlaying() bool {
	return manager.players.song != nil && manager.players.song.IsPlaying()
}

func PlaySong() {
	if manager.players.song != nil {
		if manager.players.song.IsPlaying() {
			ResumeSong()
		} else {
			manager.players.song.Play()
		}
	}
}

func PauseSong() {
	if manager.players.song != nil && manager.players.song.IsPlaying() {
		manager.players.song.Pause()
	}
}

func ResumeSong() {
	if manager.players.song != nil && !manager.players.song.IsPlaying() {
		manager.players.song.Play()
	}
}

func SetSongPositionMS(ms int) {
	if ms < 0 {
		return
	}
	if manager.players.song != nil {
		manager.players.song.SetPosition(time.Duration(ms) * time.Millisecond)
	}
}

// Stop
func StopSFX() {
	manager.sfxCache.StopAll()
}

func StopBGM() {
	if manager.players.bgm != nil && manager.players.bgm.IsPlaying() {
		manager.players.bgm.Pause()
		manager.players.bgm.Close()
	}
}

func StopSong() {
	if manager.players.song != nil && manager.players.song.IsPlaying() {
		manager.players.song.Close()
	}
}

func StopSongPreview() {
	players := []*resound.Player{
		manager.players.previewCurrent,
		manager.players.previewPrev,
	}
	for _, player := range players {
		if player != nil && player.IsPlaying() {
			player.Close()
		}
	}
	manager.previewStatus = nil
}

func StopAll() {
	StopSFX()
	StopBGM()
	StopSong()
	StopSongPreview()
}

// func FadeOutAll() {
// 	channels := []*resound.DSPChannel{
// 		manager.channels.sfx,
// 		manager.channels.bgm,
// 		manager.channels.song,
// 		manager.channels.previewCurrent,
// 		manager.channels.previewPrev,
// 	}
// 	for _, c := range channels {
// 		if c != nil {
// 			effect := getVolumeEffect(c)
// 			if effect != nil {
// 				effect.StartFade(-1, 0, config.AUDIO_FADE_S)
// 			}
// 		}
// 	}
// }

func updateSongPreview() {
	current := manager.players.previewCurrent
	if current == nil {
		return
	}

	status := manager.previewStatus
	position := current.Position().Milliseconds()
	fadeEnd := status.end - int64(config.AUDIO_FADE_MS)
	if position >= status.end {
		status.restarting = false
		current.SetPosition(time.Duration(status.start) * time.Millisecond)
		fadeCurrentPreview(0, manager.volume.Song)
	} else if (position >= fadeEnd || !current.IsPlaying()) && !status.restarting {
		fadeCurrentPreview(manager.volume.Song, 0)
		status.restarting = true
	}
}

func Update() {
	if manager == nil {
		return
	}

	if manager.previewStatus != nil {
		updateSongPreview()
	}
	if manager.players.bgm != nil {
		if manager.players.bgm.Volume() == 0 {
			manager.players.bgm.Pause()
		}
	}

	players := []*resound.Player{
		manager.players.bgm,
		manager.players.song,
		manager.players.previewCurrent,
		manager.players.previewPrev,
	}
	for _, player := range players {
		if player != nil && !player.IsPlaying() {
			player.Close()
		}
	}
}
