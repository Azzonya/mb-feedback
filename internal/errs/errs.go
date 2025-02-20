// Package errs defines custom error types for the application, providing
// consistent error messages for common scenarios like invalid input,
// missing data, and service unavailability.
package errs

// Err represents a custom error type that implements the error interface.
// It allows for defining string constants as specific errors.
type Err string

// Error returns the string representation of the Err type,
// implementing the error interface.
func (e Err) Error() string {
	return string(e)
}

const (
	InvalidInput   = Err("invalid_input")
	BadStatusCode  = Err("bad_status_code")
	ObjectNotFound = Err("object_not_found")
)
