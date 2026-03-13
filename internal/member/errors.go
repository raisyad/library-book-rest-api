package member

import "errors"

var (
	ErrMemberNotFound = errors.New("member not found")
	ErrDuplicateEmail = errors.New("email already exists")
	ErrInvalidPhone   = errors.New("phone must be a valid number")
	ErrInvalidName    = errors.New("name must be a valid name")
	ErrInvalidEmail   = errors.New("email must be a valid email")
	ErrInvalidPassword = errors.New("password must be a valid password")
	ErrInvalidConfirmPassword = errors.New("confirm password must be a valid confirm password")
)