package borrowing

import "time"

type CreateBorrowingRequest struct {
	MemberID int64 `json:"member_id" binding:"required,gt=0"`
	BookID   int64 `json:"book_id" binding:"required"`
	DueDate  *string `json:"due_date"`
}

type CreateBorrowingParams struct {
	MemberID int64
	BookID   int64
	DueDate  *time.Time
}