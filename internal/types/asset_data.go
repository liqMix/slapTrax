package types

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

// SongData represents raw song data loaded from assets
type SongData struct {
	hash       string
	FolderName string
	Meta       []byte
	Art        *ebiten.Image
	AudioPath  string
	Charts     map[int][]byte
}

// GetHash calculates and returns the hash for this song data
func (sd *SongData) GetHash() string {
	if sd.hash != "" {
		return sd.hash
	}
	hasher := sha256.New()

	// Hash the metadata
	hasher.Write(sd.Meta)

	// Hash chart data in a deterministic order
	difficulties := make([]int, 0, len(sd.Charts))
	for diff := range sd.Charts {
		difficulties = append(difficulties, diff)
	}
	sort.Ints(difficulties)

	for _, diff := range difficulties {
		// Include difficulty level in hash to ensure unique hashes even if chart data is identical
		diffBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(diffBytes, uint32(diff))
		hasher.Write(diffBytes)

		hasher.Write(sd.Charts[diff])
	}

	// Store the computed hash
	sd.hash = fmt.Sprintf("%x", hasher.Sum(nil))
	return sd.hash
}