package schema

import "fmt"

// Validation errors
var (
	ErrInvalidVersion      = fmt.Errorf("invalid schema version")
	ErrMissingTitle        = fmt.Errorf("missing song title")
	ErrMissingArtist       = fmt.Errorf("missing artist")
	ErrInvalidBPM          = fmt.Errorf("invalid BPM value")
	ErrNoCharts            = fmt.Errorf("no charts provided")
	ErrMissingChartName    = fmt.Errorf("missing chart name")
	ErrInvalidDifficulty   = fmt.Errorf("difficulty must be between 1 and 10")
	ErrInvalidTrackName    = fmt.Errorf("invalid track name")
	ErrNoteCountMismatch   = fmt.Errorf("note count doesn't match actual notes")
	ErrHoldCountMismatch   = fmt.Errorf("hold count doesn't match actual hold notes")
	ErrInvalidNoteTime     = fmt.Errorf("note time cannot be negative")
	ErrMissingHoldDuration = fmt.Errorf("hold note missing duration")
	ErrInvalidMultiNote    = fmt.Errorf("multi note must have at least 2 tracks")
	ErrInvalidNoteType     = fmt.Errorf("invalid note type")
)

// ValidationError provides context for validation failures
type ValidationError struct {
	Context string
	Field   string
	Err     error
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validation error in %s.%s: %v", e.Context, e.Field, e.Err)
}

func (e ValidationError) Unwrap() error {
	return e.Err
}

// NewValidationError creates a new validation error with context
func NewValidationError(context, field string, err error) ValidationError {
	return ValidationError{
		Context: context,
		Field:   field,
		Err:     err,
	}
}
