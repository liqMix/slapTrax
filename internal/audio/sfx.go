package audio

import (
	"fmt"
	"math/rand/v2"

	"github.com/liqmix/ebiten-holiday-2024/internal/logger"
	"github.com/liqmix/ebiten-holiday-2024/internal/types"
	"github.com/solarlune/resound"
	"github.com/solarlune/resound/effects"
)

const sfxPoolSize = 1
const pitchVariance = 0.05

type SFXPlayerCache struct {
	channels map[SFXCode]*resound.DSPChannel
	players  map[SFXCode][]*resound.Player
}

var TrackNameToPitch = map[types.TrackName]float64{
	types.TrackLeftBottom:   0.5,
	types.TrackRightBottom:  0.5,
	types.TrackCenterBottom: 0.3,

	types.TrackLeftTop:   1.0,
	types.TrackRightTop:  1.0,
	types.TrackCenterTop: 1.2,
}

func NewSFXPlayerCache(audioProperties resound.AudioProperties, volume float64) (*SFXPlayerCache, error) {
	cache := &SFXPlayerCache{
		players:  make(map[SFXCode][]*resound.Player),
		channels: make(map[SFXCode]*resound.DSPChannel),
	}

	// Pre-create several players for each sound
	for _, code := range AllSFXCodes() {
		// Create a channel for this sound
		channel := resound.NewDSPChannel()
		channel.AddEffect("volume", effects.NewVolume().SetStrength(volume))
		cache.channels[code] = channel

		// Get the stream
		stream, err := getStream(code.Path())
		if err != nil {
			return nil, fmt.Errorf("failed to decode %s: %v", code.Path(), err)
		}

		volumeEffect := getVolumeEffect(channel)
		if volumeEffect != nil {
			prop, err := audioProperties.Get(code.Path()).Analyze(stream, 0)
			if err != nil {
				logger.Error("Error analyzing song (%s): %v", code.Path(), err)
			} else {
				volumeEffect.SetNormalizationFactor(prop.Normalization)
			}
		}

		// Create a pool of players for this sound
		players := make([]*resound.Player, sfxPoolSize) // Adjust pool size as needed
		for i := 0; i < len(players); i++ {
			player, err := resound.NewPlayer(stream)
			player.SetBufferSize(64)
			player.SetDSPChannel(channel)
			if err != nil {
				return nil, fmt.Errorf("failed to create player for %s: %v", code.Path(), err)
			}

			// Add pitch variance
			if code != SFXHat {
				var pitch float64
				trackName := SFXTrack(code)
				if trackName != types.TrackUnknown {
					pitch = TrackNameToPitch[trackName]
				} else {
					pitch = 0.5
				}
				player.AddEffect("pitch", effects.NewPitchShift(256).SetSource(player).SetPitch(pitch))
			}

			players[i] = player
		}
		cache.players[code] = players
	}
	return cache, nil
}

func (c *SFXPlayerCache) StopAll() {
	for _, players := range c.players {
		for _, player := range players {
			player.Pause()
			player.Rewind()
		}
	}
}

func (c *SFXPlayerCache) SetVolume(v float64) {
	for _, channel := range c.channels {
		if vol := channel.Effects["volume"]; vol != nil {
			if volumeEffect, ok := vol.(*effects.Volume); ok {
				volumeEffect.SetStrength(v)
			}
		} else {
			channel.AddEffect("volume", effects.NewVolume().SetStrength(v))
		}
	}
}

func (c *SFXPlayerCache) PlayTrackSound(trackName types.TrackName) {
	code := TrackSFX(trackName)
	if code == SFXNone {
		return
	}
	players := c.players[code]
	if players == nil {
		return
	}

	// Find an available player
	for _, player := range players {
		// if !player.IsPlaying() {
		player.Rewind()
		pitchEffect := player.Effects["pitch"]
		if pitchEffect != nil {
			pitchEffect.(*effects.PitchShift).SetPitch(TrackNameToPitch[trackName] + pitchVariance*(rand.Float64()-0.5))
		}
		player.Play()
	}

	c.PlaySound(code)
}

func (c *SFXPlayerCache) PlaySound(code SFXCode) error {
	players := c.players[code]
	if players == nil {
		return fmt.Errorf("sound %s not found", code.Path())
	}

	// Find an available player
	for _, player := range players {
		// if !player.IsPlaying() {
		player.Rewind()
		player.Play()
	}

	// All players busy - either wait or skip
	// Could also dynamically create a new player if needed
	return fmt.Errorf("no available players for %s", code.Path())
}
