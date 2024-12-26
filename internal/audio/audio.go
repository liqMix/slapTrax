package audio

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/liqmix/ebiten-holiday-2024/internal/audio/sfx"
	"github.com/liqmix/ebiten-holiday-2024/internal/config"
	"github.com/liqmix/ebiten-holiday-2024/internal/song"
	"github.com/liqmix/ebiten-holiday-2024/internal/user"
	"github.com/solarlune/resound"
	"github.com/solarlune/resound/effects"
)

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

type audioMan struct {
	cache    map[string][]byte
	players  players
	channels channels
}

var (
	manager audioMan
)

func InitAudioManager() {
	audio.NewContext(config.SAMPLE_RATE)
	manager = audioMan{
		cache: prewarmCache(),
		channels: channels{
			sfx:         newSfxChannel(),
			bgm:         newMusicChannel(),
			song:        newSongChannel(),
			songPreview: newSongPreviewChannels(),
		},
		players: players{
			sfx: make([]*resound.Player, 0),
			songPreview: &songPreviewPlayers{
				current:  nil,
				previous: nil,
			},
		},
	}
}

// Preload SFX and BGM
func prewarmCache() map[string][]byte {
	cache := make(map[string][]byte, 0)
	return cache
}

func newSfxChannel() *resound.DSPChannel {
	sfxChannel := resound.NewDSPChannel()
	// sfx effects
	return sfxChannel
}

func newMusicChannel() *resound.DSPChannel {
	musicChannel := resound.NewDSPChannel()
	// music effects
	return musicChannel
}

func newSongChannel() *resound.DSPChannel {
	songChannel := resound.NewDSPChannel()
	// song effects
	return songChannel
}

func newSongPreviewChannels() *songPreviewChannels {
	previous := resound.NewDSPChannel()
	previous.AddEffect("volume", effects.NewVolume())

	current := resound.NewDSPChannel()
	current.AddEffect("volume", effects.NewVolume())

	// song preview effects
	return &songPreviewChannels{
		current:  current,
		previous: previous,
	}
}

// // Finds a file with the given name and any extension
// func findAudioPath(parentPath string, name string) (string, error) {
// 	// Find the first file with the given name
// 	files, err := os.ReadDir(parentPath)
// 	if err != nil {
// 		return "", err
// 	}
// 	for _, file := range files {
// 		if strings.HasPrefix(file.Name(), name) {
// 			return path.Join(parentPath, file.Name()), nil
// 		}
// 	}
// 	return "", fmt.Errorf("file not found: %s", name)
// }

// Just read the file into memory for now
func readAudio(path string) ([]byte, error) {
	if data, ok := manager.cache[path]; ok {
		return data, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	manager.cache[path] = data
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

func PlaySFX(sfxCode sfx.SFXCode) {
	path := path.Join(config.SFX_DIR, string(sfxCode), ".wav")
	player, err := getPlayer(path)
	if err != nil {
		fmt.Printf("Error playing SFX: %v\n", err)
		return
	}
	player.SetDSPChannel(manager.channels.sfx)
	manager.players.sfx = append(manager.players.sfx, player)
	player.Play()
}

func PlayBGM(bgmCode string) {
	path := path.Join(config.BGM_DIR, bgmCode)
	player, err := getPlayer(path)
	if err != nil {
		fmt.Printf("Error playing BGM: %v\n", err)
		return
	}
	player.SetDSPChannel(manager.channels.bgm)
	manager.players.bgm = player
	player.Play()
}

// TODO: get the fade working right for seamless transitions between selected songs
func PlaySongPreview(s *song.Song) {
	position := s.PreviewStart
	player, err := getPlayer(s.AudioPath)
	if err != nil {
		fmt.Printf("Error playing song preview: %v\n", err)
		return
	}

	// If there is a current song preview, move it to previous and fade out
	current := manager.players.songPreview.current
	if current != nil {
		if manager.players.songPreview.previous != nil {
			manager.players.songPreview.previous.Close()
		}
		previousChannel := manager.channels.songPreview.previous
		current.SetDSPChannel(previousChannel)
		previousVolume := previousChannel.Effects["volume"].(*effects.Volume)
		previousVolume.StartFade(current.Volume(), 0, config.SONG_PREVIEW_FADE/1000)
		manager.players.songPreview.previous = current
	}

	// Fade in the new song preview
	player.SetDSPChannel(manager.channels.songPreview.current)
	manager.players.songPreview.current = player
	currentChannel := manager.channels.songPreview.current
	currentVolume := currentChannel.Effects["volume"].(*effects.Volume)
	currentVolume.StartFade(0, user.Current.Settings.SongPreviewVolume, config.SONG_PREVIEW_FADE/1000)

	player.SetPosition(time.Duration(position) * time.Millisecond)
	manager.players.songPreview.current = player
	player.Play()
}

func StopSFX() {
	for _, player := range manager.players.sfx {
		if player.IsPlaying() {
			player.Pause()
		}
		player.Close()
	}
}

func StopBGM() {
	if manager.players.bgm != nil && manager.players.bgm.IsPlaying() {
		manager.players.bgm.Close()
	}
}

// TODO: collect these into something like audio.Song.InitSong()...
func InitSong(s *song.Song) {
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
		return manager.players.song.Position().Milliseconds()
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
	if manager.players.songPreview.current != nil && manager.players.songPreview.current.IsPlaying() {
		manager.players.songPreview.current.Close()
	}
	if manager.players.songPreview.previous != nil && manager.players.songPreview.previous.IsPlaying() {
		manager.players.songPreview.previous.Close()
	}
}

func StopAll() {
	StopSFX()
	StopBGM()
	StopSong()
	StopSongPreview()
}

func Update() {
	for _, player := range manager.players.sfx {
		if !player.IsPlaying() {
			player.Close()
		}
	}
	if manager.players.bgm != nil && !manager.players.bgm.IsPlaying() {
		fmt.Println("Closing bgm")
		manager.players.bgm.Close()
	}
	if manager.players.song != nil && !manager.players.song.IsPlaying() {
		fmt.Println("Closing song")
		manager.players.song.Close()
	}
	if manager.players.songPreview.current != nil && !manager.players.songPreview.current.IsPlaying() {
		manager.players.songPreview.current.Close()
	}
	if manager.players.songPreview.previous != nil && !manager.players.songPreview.previous.IsPlaying() {
		manager.players.songPreview.previous.Close()
	}
}
