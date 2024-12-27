package song

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"github.com/liqmix/ebiten-holiday-2024/internal/l"
)

type Difficulty int

func (d Difficulty) String() string {
	if d < 5 {
		return l.String(l.DIFFICULTY_EASY)
	}
	if d < 8 {
		return l.String(l.DIFFICULTY_MEDIUM)
	}
	if d <= 10 {
		return l.String(l.DIFFICULTY_HARD)
	}
	return l.String(l.DIFFICULTY_UNKNOWN)
}

type Chart struct {
	Difficulty Difficulty
	TotalNotes int
	Tracks     []Track
}

const (
	// MIDI event types
	noteOff   = 0x80
	noteOn    = 0x90
	metaEvent = 0xFF
)

// MIDI note numbers to track mapping
var noteToTrack = map[uint8]TrackName{
	62: LeftTop,     // D5
	60: LeftBottom,  // C5
	57: Center,      // A4
	55: RightBottom, // G4
	53: RightTop,    // F4

	// Edge taps
	50: EdgeTop,  // D4
	48: EdgeTap1, // C4
	47: EdgeTap2, // B3
	46: EdgeTap3, // A#3
}

// Parse the chart file into a set of tracks and associated notes
func ParseChart(song *Song, data []byte) (*Chart, error) {
	chart := &Chart{}
	chart.Tracks = make([]Track, 0)
	notes := make(map[TrackName][]*Note)
	for _, name := range TrackNames() {
		notes[name] = []*Note{}
	}

	reader := bytes.NewReader(data)

	// Read header chunk
	var headerChunk struct {
		ID       [4]byte
		Length   uint32
		Format   uint16
		Tracks   uint16
		Division uint16
	}
	if err := binary.Read(reader, binary.BigEndian, &headerChunk); err != nil {
		return nil, err
	}
	if string(headerChunk.ID[:]) != "MThd" {
		return nil, errors.New("invalid MIDI file: missing MThd header")
	}

	// Midi junk to get timing
	ticksPerQuarter := float64(headerChunk.Division)
	msPerTick := (60000.0 / float64(song.BPM)) / ticksPerQuarter

	// Vars to track current state
	var currentTs int64
	activeTracks := make(map[TrackName]int64)

	// Only one track for now
	var trackHeader struct {
		ID     [4]byte
		Length uint32
	}
	if err := binary.Read(reader, binary.BigEndian, &trackHeader); err != nil {
		if err == io.EOF {
			return nil, errors.New("invalid MIDI file: missing track header")
		}
		return nil, err
	}

	if string(trackHeader.ID[:]) != "MTrk" {
		return nil, errors.New("invalid MIDI file: missing MTrk header")
	}

	var running uint8
	// A note must be held for at least this duration to be considered a hold note
	// - 1/32 note at song.BPM
	MIN_HOLD_DURATION := int64(msPerTick * 4)

	for {
		// Read delta time
		delta, err := readVariableLength(reader)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		currentTs += int64(float64(delta) * msPerTick)

		// Read event type
		var eventType uint8
		if err := binary.Read(reader, binary.BigEndian, &eventType); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		// Handle running status
		if eventType&0x80 == 0 {
			eventType = running
			reader.Seek(-1, io.SeekCurrent)
		} else {
			running = eventType
		}

		command := eventType & 0xF0
		switch command {
		case noteOn, noteOff:
			var note, velocity uint8
			if err := binary.Read(reader, binary.BigEndian, &note); err != nil {
				if err == io.EOF {
					break
				}
				return nil, err
			}
			if err := binary.Read(reader, binary.BigEndian, &velocity); err != nil {
				if err == io.EOF {
					break
				}
				return nil, err
			}

			// Ignore notes not mapped to a track
			trackName, ok := noteToTrack[note]
			if !ok {
				continue
			}

			// Note on with velocity > 0 starts a note
			if command == noteOn && velocity > 0 {
				activeTracks[trackName] = currentTs
			} else {
				// Otherwise end the note (if it's active)
				if active, ok := activeTracks[trackName]; ok {
					duration := currentTs - active
					if duration <= MIN_HOLD_DURATION {
						duration = 0
					}
					notes[trackName] = append(notes[trackName], NewNote(active, currentTs))
					delete(activeTracks, trackName)
				}
			}
		}
	}

	// End any remaining notes
	for trackName, target := range activeTracks {
		release := currentTs
		if (currentTs - target) <= MIN_HOLD_DURATION {
			release = 0
		}
		notes[trackName] = append(notes[trackName], NewNote(target, release))
	}

	// Create tracks from notes
	for _, name := range TrackNames() {
		track := NewTrack(name, notes[name])
		chart.Tracks = append(chart.Tracks, *track)
		chart.TotalNotes += len(notes[name])
	}

	// Can probably do some chart metadata (?) here since we have all info
	return chart, nil
}

func readVariableLength(reader io.Reader) (uint32, error) {
	var result uint32
	for {
		var b uint8
		if err := binary.Read(reader, binary.BigEndian, &b); err != nil {
			return 0, err
		}
		result = (result << 7) | uint32(b&0x7F)
		if b&0x80 == 0 {
			break
		}
	}
	return result, nil
}
