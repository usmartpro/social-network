package v1

// HTTPError ...
type HTTPError struct {
	err        Error
	statusCode int
}

// Error ...
func (err *HTTPError) Error() string {
	return err.err.Message
}

// StatusCode ...
func (err *HTTPError) StatusCode() int {
	return err.statusCode
}

// HTTPError ...
func (err *HTTPError) HTTPError() Error {
	return err.err
}

// NewHTTPError ...
func NewHTTPError(statusCode int, err Error) *HTTPError {
	return &HTTPError{
		statusCode: statusCode,
		err:        err,
	}
}
