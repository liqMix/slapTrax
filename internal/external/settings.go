package external

// MergeFrom merges settings from another settings object
func (s *Settings) MergeFrom(other *Settings) {
	if other == nil {
		return
	}

	// Only update if non-zero values exist
	if other.Version > 0 {
		s.Version = other.Version
	}
	if !other.LastModified.IsZero() {
		s.LastModified = other.LastModified
	}
	if !other.LastSync.IsZero() {
		s.LastSync = other.LastSync
	}
	if other.Locale != "" {
		s.Locale = other.Locale
	}
	if other.ScreenWidth > 0 {
		s.ScreenWidth = other.ScreenWidth
	}
	if other.ScreenHeight > 0 {
		s.ScreenHeight = other.ScreenHeight
	}
	if other.RenderWidth > 0 {
		s.RenderWidth = other.RenderWidth
	}
	if other.RenderHeight > 0 {
		s.RenderHeight = other.RenderHeight
	}

	// Boolean flags
	s.Fullscreen = other.Fullscreen
	s.FixedRenderScale = other.FixedRenderScale
	s.DisableHoldNotes = other.DisableHoldNotes
	s.DisableHitEffects = other.DisableHitEffects
	s.DisableLaneEffects = other.DisableLaneEffects

	// Only update if positive values
	if other.BGMVolume > 0 {
		s.BGMVolume = other.BGMVolume
	}
	if other.SFXVolume > 0 {
		s.SFXVolume = other.SFXVolume
	}
	if other.SongVolume > 0 {
		s.SongVolume = other.SongVolume
	}
	if other.LaneSpeed > 0 {
		s.LaneSpeed = other.LaneSpeed
	}

	// Update offsets if non-zero
	if other.AudioOffset != 0 {
		s.AudioOffset = other.AudioOffset
	}
	if other.InputOffset != 0 {
		s.InputOffset = other.InputOffset
	}

	// Update colors if non-empty
	if other.NoteColorTheme != "" {
		s.NoteColorTheme = other.NoteColorTheme
	}
	if other.CenterNoteColor != "" {
		s.CenterNoteColor = other.CenterNoteColor
	}
	if other.CornerNoteColor != "" {
		s.CornerNoteColor = other.CornerNoteColor
	}
}

// Clone creates a deep copy of settings
func (s *Settings) Clone() *Settings {
	if s == nil {
		return nil
	}

	clone := *s // Shallow copy

	return &clone
}
