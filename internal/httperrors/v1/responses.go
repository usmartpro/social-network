package v1

const (
	// CodeValidationError ...
	CodeValidationError = "validation_errors"
)

// Error ...
type Error struct {
	Code          string `json:"code,omitempty"`
	Message       string `json:"message,omitempty"`
	Error         string `json:"error,omitempty"`
	RetryAfterSec int    `json:"retryAfterSec,omitempty"`
}
