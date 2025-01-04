package assets

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"path"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/liqmix/ebiten-holiday-2024/internal/cache"
	"github.com/liqmix/ebiten-holiday-2024/internal/config"
	"github.com/liqmix/ebiten-holiday-2024/internal/logger"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/solarlune/resound"
	"github.com/solarlune/resound/effects"
)

//go:embed sfx/*.mp3 songs/**/*.mp3
var sfxFS embed.FS

type BGMCode string

const (
	BGMIntro   BGMCode = "intro"
	BGMMenu    BGMCode = "menu"
	BGMResults BGMCode = "results"
)

func (b BGMCode) Path() string {
	return path.Join(config.BGM_DIR, string(b)+".mp3")
}

type SFXCode string

const (
	SFXOffset     SFXCode = "offset"
	SFXHat        SFXCode = "hat"
	SFXSelectUp   SFXCode = "selectup"
	SFXSelectDown SFXCode = "selectdown"
	// SFXNoteHit          SFXCode = "hit"
	// SFXSelect     SFXCode = "select"
)

func (s SFXCode) Path() string {
	return path.Join(config.SFX_DIR, string(s)+".mp3")
}

func AllSFX() []SFXCode {
	return []SFXCode{
		SFXOffset,
		SFXHat,
		SFXSelectUp,
		SFXSelectDown,
	}
}

type songPreviewChannels struct {
	current  *resound.DSPChannel
	previous *resound.DSPChannel
}
type channels struct {
	sfx         *resound.DSPChannel
	bgm         *resound.DSPChannel
	song        *resound.DSPChannel
	songPreview *songPreviewChannels
}

type songPreviewPlayers struct {
	current  *resound.Player
	previous *resound.Player
}
type players struct {
	sfx         []*resound.Player
	bgm         *resound.Player
	song        *resound.Player
	songPreview *songPreviewPlayers
}

type preview struct {
	start      int64
	end        int64
	restarting bool
}

type audioMan struct {
	cache    map[string][]byte
	players  players
	channels *channels

	songPreview *preview
	volume      *Volume
}

type Volume struct {
	Sfx  float64
	Bgm  float64
	Song float64
}

var (
	manager audioMan
)

func InitAudioManager(v *Volume) {
	logger.Info("Initializing audio manager...")
	audio.NewContext(config.SAMPLE_RATE)
	manager = audioMan{
		channels: newChannels(v),
		players: players{
			sfx: make([]*resound.Player, 0),
			songPreview: &songPreviewPlayers{
				current:  nil,
				previous: nil,
			},
		},
		volume: v,
	}
	logger.Debug("Loading sfx...")
	for _, sfx := range AllSFX() {
		logger.Debug("  Loading %s...", sfx.Path())
		_, err := readAudio(sfx.Path())
		if err != nil {
			logger.Error("    Error loading audio: %v", err)
		}
	}
}

func newChannels(v *Volume) *channels {
	sfx := resound.NewDSPChannel()
	bgm := resound.NewDSPChannel()
	song := resound.NewDSPChannel()
	previewPrev := resound.NewDSPChannel()
	previewCurrent := resound.NewDSPChannel()

	sfx.AddEffect("volume", effects.NewVolume().SetStrength(v.Sfx))
	bgm.AddEffect("volume", effects.NewVolume().SetStrength(v.Bgm))
	song.AddEffect("volume", effects.NewVolume().SetStrength(v.Song))
	previewPrev.AddEffect("volume", effects.NewVolume().SetStrength(v.Song))
	previewCurrent.AddEffect("volume", effects.NewVolume().SetStrength(v.Song))

	return &channels{
		sfx:  sfx,
		bgm:  bgm,
		song: song,
		songPreview: &songPreviewChannels{
			current:  previewCurrent,
			previous: previewPrev,
		},
	}
}

// Just read the file into memory for now
func readAudio(path string) ([]byte, error) {
	if a, ok := cache.GetAudio(path); ok {
		return a, nil
	}
	data, err := sfxFS.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cache.SetAudio(path, data)
	return data, nil
}

// getAudioPlayer returns a player for the given audio file
//
//	supported formats: wav, ogg, mp3
func getPlayer(path string) (*resound.Player, error) {
	data, err := readAudio(path)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(data)
	audioExt := strings.Split(path, ".")[1]

	var stream io.ReadSeeker
	switch audioExt {
	case "wav":
		stream, err = wav.DecodeWithSampleRate(config.SAMPLE_RATE, reader)
	case "ogg":
		stream, err = vorbis.DecodeWithSampleRate(config.SAMPLE_RATE, reader)
	case "mp3":
		stream, err = mp3.DecodeWithSampleRate(config.SAMPLE_RATE, reader)
	default:
		panic("Unsupported audio format")
	}
	if err != nil {
		panic(err)
	}

	player, err := resound.NewPlayer(stream)
	if err != nil {
		return nil, err
	}
	return player, nil
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

func PlayBGM(bgmCode string) {
	path := path.Join(config.BGM_DIR, bgmCode)
	if err := play(path, manager.channels.bgm, nil); err != nil {
		logger.Error("Error playing BGM: %v", err)
	}
}

func fadeCurrentPreview(start float64, end float64) {
	currentChannel := manager.channels.songPreview.current
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
	current := manager.players.songPreview.current
	if current != nil {
		// just stop it for now
		current.Close()

		//// TODO: fade out
		// if manager.players.songPreview.previous != nil {
		// 	manager.players.songPreview.previous.Close()
		// }
		// previousChannel := manager.channels.songPreview.previous
		// current.SetDSPChannel(previousChannel)
		// previousVolume := previousChannel.Effects["volume"].(*effects.Volume)
		// previousVolume.StartFade(current.Volume(), 0, config.SONG_PREVIEW_FADE/1000)
		// manager.players.songPreview.previous = current
	}

	// Fade in the new song preview
	player.SetDSPChannel(manager.channels.songPreview.current)
	manager.players.songPreview.current = player
	player.SetPosition(time.Duration(startPosition) * time.Millisecond)
	fadeCurrentPreview(-1, manager.volume.Song)
	manager.songPreview = &preview{
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
		getVolumeEffect(manager.channels.bgm).SetStrength(v.Bgm)
	}
	if manager.channels.song != nil {
		getVolumeEffect(manager.channels.song).SetStrength(v.Song)
	}
	if manager.channels.sfx != nil {
		getVolumeEffect(manager.channels.sfx).SetStrength(v.Sfx)
	}
	if manager.channels.songPreview.current != nil {
		getVolumeEffect(manager.channels.songPreview.current).SetStrength(v.Song)
	}
	if manager.channels.songPreview.previous != nil {
		getVolumeEffect(manager.channels.songPreview.previous).SetStrength(v.Song)
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
		manager.players.song.Play()
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

func StopSong() {
	if manager.players.song != nil && manager.players.song.IsPlaying() {
		manager.players.song.Close()
	}
}

func StopSongPreview() {
	players := []*resound.Player{
		manager.players.songPreview.current,
		manager.players.songPreview.previous,
	}
	for _, player := range players {
		if player != nil && player.IsPlaying() {
			player.Close()
		}
	}
	manager.songPreview = nil
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
// 		manager.channels.songPreview.current,
// 		manager.channels.songPreview.previous,
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
	current := manager.players.songPreview.current
	if current == nil {
		return
	}

	preview := manager.songPreview
	position := current.Position().Milliseconds()
	fadeEnd := preview.end - int64(config.AUDIO_FADE_MS)
	if position >= preview.end {
		preview.restarting = false
		current.SetPosition(time.Duration(preview.start) * time.Millisecond)
		fadeCurrentPreview(0, manager.volume.Song)
	} else if position >= fadeEnd && !preview.restarting {
		fadeCurrentPreview(manager.volume.Song, 0)
		preview.restarting = true
	}
}

func Update() {
	if manager.songPreview != nil {
		updateSongPreview()
	}

	for _, player := range manager.players.sfx {
		if !player.IsPlaying() {
			player.Close()
		}
	}
	players := []*resound.Player{
		manager.players.bgm,
		manager.players.song,
		manager.players.songPreview.current,
		manager.players.songPreview.previous,
	}
	for _, player := range players {
		if player != nil && !player.IsPlaying() {
			player.Close()
		}
	}
}
