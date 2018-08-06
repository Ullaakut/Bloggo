package errortype

import "errors"

// Error definitions
var (
	ErrNotFound            = errors.New("resource not found")
	ErrConflict            = errors.New("datamodel conflict")
	ErrDuplicateEntry      = errors.New("duplicate entry")
	ErrUnprocessableEntity = errors.New("unprocessable entity")
)
