package customerrors

import (
	"errors"
)

var (
	ErrNotFound            = errors.New("not found")
	ErrAlreadyExists       = errors.New("already exists")
	ErrValidation          = errors.New("validation error")
	ErrInternalServerError = errors.New("internal server error")
)
