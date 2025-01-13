package types

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"github.com/liqmix/slaptrax/internal/logger"
)

const (
	// MIDI event types
	noteOff = 0x80
	noteOn  = 0x90
)

var noteToTrack = map[uint8]TrackName{
	74: TrackLeftTop,      // D6
	73: TrackLeftBottom,   // C#6
	72: TrackCenterTop,    // C6
	71: TrackCenterBottom, // B5
	70: TrackRightTop,     // A#5
	69: TrackRightBottom,  // A5
}

type Chart struct {
	TotalNotes     int
	TotalHoldNotes int
	Tracks         []*Track
}

func NewChart(song *Song, data []byte) (*Chart, error) {
	logger.Debug("Parsing chart for %s", song.Title)
	chart := &Chart{}
	chart.Tracks = make([]*Track, 0)
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
	// logger.Debug("MIDI Header: length %d bytes", headerChunk.Length)
	// logger.Debug("Format: %d", headerChunk.Format)
	// logger.Debug("Tracks: %d", headerChunk.Tracks)
	// logger.Debug("Division: %d", headerChunk.Division)

	// Midi junk to get timing
	ticksPerQuarter := float64(headerChunk.Division)
	msPerTick := (60000.0 / float64(song.BPM)) / ticksPerQuarter

	// A note must be held for at least this duration to be considered a hold note
	minHoldDuration := int64(msPerTick * (ticksPerQuarter / 2))

	// Vars to track current state
	var currentTs int64 = 0
	var running uint8 = 0

	activeTracks := make(map[TrackName]int64)

	for trackNum := uint16(0); trackNum < headerChunk.Tracks; trackNum++ {
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

		// logger.Debug("Track %d: length %d bytes", trackNum, trackHeader.Length)

		// Track 0 is usually metadata/tempo - skip it (for now... maybe we can use this instead of songmeta for bpm)
		if trackNum == 0 {
			reader.Seek(int64(trackHeader.Length), io.SeekCurrent)
			continue
		}

		// Track event reading loop
		startPos, _ := reader.Seek(0, io.SeekCurrent) // Get current position
		running = 0

		for {
			currentPos, _ := reader.Seek(0, io.SeekCurrent)
			if currentPos-startPos >= int64(trackHeader.Length) {
				break // We've read all data in this track
			}

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
				if running == 0 {
					logger.Debug("Warning: Got data byte 0x%02X with no running status!", eventType)
				}
				eventType = running
				reader.Seek(-1, io.SeekCurrent)
			} else {
				running = eventType
			}

			// First handle meta events (before the command masking)
			if eventType == 0xFF {
				var metaType uint8
				if err := binary.Read(reader, binary.BigEndian, &metaType); err != nil {
					if err == io.EOF {
						break
					}
					return nil, err
				}

				length, err := readVariableLength(reader)
				if err != nil {
					if err == io.EOF {
						break
					}
					return nil, err
				}
				reader.Seek(int64(length), io.SeekCurrent)
				running = 0 // Reset running status after meta event
				continue    // Skip to next event
			}

			// Handle event type
			var eventReachedEOF bool
			command := eventType & 0xF0
			switch command {
			case noteOn, noteOff:
				var note, velocity uint8
				if err := binary.Read(reader, binary.BigEndian, &note); err != nil {
					if err == io.EOF {
						eventReachedEOF = true
						break
					}
					return nil, err
				}
				if err := binary.Read(reader, binary.BigEndian, &velocity); err != nil {
					if err == io.EOF {
						eventReachedEOF = true
						break
					}
					return nil, err
				}

				// Ignore notes not mapped to a track
				trackName, ok := noteToTrack[note]
				if ok {
					// Note on with velocity > 0 starts a note
					if command == noteOn && velocity > 0 {
						activeTracks[trackName] = currentTs
					} else {
						// Otherwise end the note (if it's active)
						if active, ok := activeTracks[trackName]; ok {
							release := currentTs
							if (currentTs - active) <= minHoldDuration {
								release = 0
							}
							notes[trackName] = append(notes[trackName], NewNote(trackName, active, release))
							delete(activeTracks, trackName)
						}
					}
				}
			}

			if eventReachedEOF {
				break
			}
		}
	}

	// End any remaining notes
	for trackName, target := range activeTracks {
		release := currentTs
		if (currentTs - target) <= minHoldDuration {
			release = 0
		}
		notes[trackName] = append(notes[trackName], NewNote(trackName, target, release))
	}

	// Iterate through all tracks and identify notes that have the same start time.
	// mark them as non-solo notes
	noteCounts := make(map[int64]int)
	for _, trackNotes := range notes {
		for _, note := range trackNotes {
			noteCounts[note.Target]++
		}
	}

	for _, trackNotes := range notes {
		for _, note := range trackNotes {
			if noteCounts[note.Target] > 1 {
				note.SetSolo(false)
			}
		}
	}

	// Create tracks from notes
	beatInterval := int64(msPerTick * 4)
	for _, name := range TrackNames() {
		track := NewTrack(name, notes[name], beatInterval)
		chart.Tracks = append(chart.Tracks, track)
		chart.TotalNotes += len(notes[name])
		for _, note := range notes[name] {
			if note.IsHoldNote() {
				chart.TotalHoldNotes++
			}
		}

	}

	if chart.TotalNotes == 0 {
		return nil, errors.New("no notes found in chart")
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
