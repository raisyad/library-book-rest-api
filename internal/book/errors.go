package book

import "errors"

var (
	ErrBookNotFound  = errors.New("book not found")
	ErrDuplicateISBN = errors.New("book isbn already exists")
	ErrInvalidStock  = errors.New("stock must be greater than 0")
	ErrInvalidPublishedYear = errors.New("published year must be greater than 0")
)