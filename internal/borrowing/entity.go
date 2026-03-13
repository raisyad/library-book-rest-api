package borrowing

import "time"

type Borrowing struct {
	ID        int64     `db:"id" json:"id"`
	MemberID  int64     `db:"member_id" json:"member_id"`
	BookID    int64     `db:"book_id" json:"book_id"`
	BorrowedAt time.Time `db:"borrowed_at" json:"borrowed_at"`
	DueDate   *time.Time `db:"due_date" json:"due_date,omitempty"`
	ReturnedAt *time.Time `db:"returned_at" json:"returned_at,omitempty"`
	Status    string    `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	
	MemberName  string `db:"member_name" json:"member_name"`
	MemberEmail string `db:"member_email" json:"member_email"`
	BookTitle   string `db:"book_title" json:"book_title"`
	BookAuthor  string `db:"book_author" json:"book_author"`
}