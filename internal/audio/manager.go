package audio

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/liqmix/ebiten-holiday-2024/internal/assets"
	"github.com/liqmix/ebiten-holiday-2024/internal/config"
	"github.com/liqmix/ebiten-holiday-2024/internal/logger"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/solarlune/resound"
	"github.com/solarlune/resound/effects"
)

type channels struct {
	sfx            *resound.DSPChannel
	bgm            *resound.DSPChannel
	song           *resound.DSPChannel
	previewPrev    *resound.DSPChannel
	previewCurrent *resound.DSPChannel
}

type players struct {
	sfx            []*resound.Player
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
	players  *players
	channels *channels

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

	manager = &audioMan{
		channels: newChannels(v),
		players: &players{
			sfx: make([]*resound.Player, 0),
		},
		volume: v,
	}
}

func newChannels(v *Volume) *channels {
	sfx := resound.NewDSPChannel()
	bgm := resound.NewDSPChannel()
	song := resound.NewDSPChannel()
	previewPrev := resound.NewDSPChannel()
	previewCurrent := resound.NewDSPChannel()

	sfx.AddEffect("volume", effects.NewVolume().SetStrength(v.SFX))
	bgm.AddEffect("volume", effects.NewVolume().SetStrength(v.BGM))
	song.AddEffect("volume", effects.NewVolume().SetStrength(v.Song))
	previewPrev.AddEffect("volume", effects.NewVolume().SetStrength(v.Song))
	previewCurrent.AddEffect("volume", effects.NewVolume().SetStrength(v.Song))

	return &channels{
		sfx:            sfx,
		bgm:            bgm,
		song:           song,
		previewPrev:    previewPrev,
		previewCurrent: previewCurrent,
	}
}

func getPlayer(path string) (*resound.Player, error) {
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

	return resound.NewPlayer(stream)
}

func play(path string, c *resound.DSPChannel, players *[]*resound.Player) error {
	player, err := getPlayer(path)
	if err != nil {
		return err
	}
	player.SetDSPChannel(c)
	player.Play()
	if players != nil {
		*players = append(*players, player)
	}
	return nil
}

func PlaySFX(sfx SFXCode) {
	if err := play(sfx.Path(), manager.channels.sfx, &manager.players.sfx); err != nil {
		logger.Error("Error playing SFX: %v", err)
	}
}

func PlaySFXWithOffset(sfx SFXCode, offset int64) {
	player, err := getPlayer(sfx.Path())
	if err != nil {
		logger.Error("Error playing SFX with offset: %v", err)
		return
	}
	player.SetDSPChannel(manager.channels.sfx)
	player.SetPosition(time.Duration(offset) * time.Millisecond)
	player.Play()
	manager.players.sfx = append(manager.players.sfx, player)
}

func PlayBGM(bgmCode BGMCode) {
	path := bgmCode.Path()
	if err := play(path, manager.channels.bgm, nil); err != nil {
		logger.Error("Error playing BGM: %v", err)
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

// TODO: fade?
func PlaySongPreview(s *types.Song) {
	startPosition := s.PreviewStart
	endPosition := s.PreviewStart + config.SONG_PREVIEW_LENGTH
	player, err := getPlayer(s.AudioPath)
	if err != nil {
		fmt.Printf("Error playing song preview: %v\n", err)
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

func getVolumeEffect(c *resound.DSPChannel) *effects.Volume {
	if effect, ok := c.Effects["volume"]; !ok {
		return nil
	} else {
		return effect.(*effects.Volume)
	}
}

func SetVolume(v *Volume) {
	manager.volume = v
	if manager.channels.bgm != nil {
		getVolumeEffect(manager.channels.bgm).SetStrength(v.BGM)
	}
	if manager.channels.song != nil {
		getVolumeEffect(manager.channels.song).SetStrength(v.Song)
	}
	if manager.channels.sfx != nil {
		getVolumeEffect(manager.channels.sfx).SetStrength(v.SFX)
	}
	if manager.channels.previewCurrent != nil {
		getVolumeEffect(manager.channels.previewCurrent).SetStrength(v.Song)
	}
	if manager.channels.previewPrev != nil {
		getVolumeEffect(manager.channels.previewPrev).SetStrength(v.Song)
	}
}

func StopSFX() {
	for _, player := range manager.players.sfx {
		player.Pause()
		player.Close()
	}
}

func StopBGM() {
	if manager.players.bgm != nil && manager.players.bgm.IsPlaying() {
		manager.players.bgm.Pause()
		manager.players.bgm.Close()
	}
}

// TODO: collect these into something like audio.Song.InitSong()...
func InitSong(s *types.Song) {
	player, err := getPlayer(s.AudioPath)
	if err != nil {
		fmt.Printf("Error preparing song: %v\n", err)
		return
	}
	player.SetDSPChannel(manager.channels.song)
	manager.players.song = player

}

func CurrentSongPositionMS() int64 {
	if manager.players.song != nil {
		return manager.players.song.Position().Milliseconds() + config.INHERENT_OFFSET
	}
	panic("No song playing!")
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

	for _, player := range manager.players.sfx {
		if player == nil {
			continue
		}

		if !player.IsPlaying() {
			player.Close()
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
