package borrowing

import "errors"

var (
	ErrBorrowingNotFound       = errors.New("borrowing not found")
	ErrBorrowingAlreadyReturned = errors.New("borrowing already returned")
	ErrMemberNotFound          = errors.New("member not found")
	ErrBookNotFound            = errors.New("book not found")
	ErrBookOutOfStock          = errors.New("book is out of stock")
	ErrInvalidDueDate          = errors.New("invalid due date format")
)