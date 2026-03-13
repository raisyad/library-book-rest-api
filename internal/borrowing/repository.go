package borrowing

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

const borrowingDetailQuery = `
	SELECT
		br.id,
		br.member_id,
		br.book_id,
		br.borrowed_at,
		br.due_date,
		br.returned_at,
		br.status,
		br.created_at,
		br.updated_at,
		m.name AS member_name,
		m.email AS member_email,
		b.title AS book_title,
		b.author AS book_author
	FROM borrowings br
	JOIN members m ON m.id = br.member_id
	JOIN books b ON b.id = br.book_id
`

func (r *Repository) FindAll() ([]Borrowing, error) {
	query := borrowingDetailQuery + `
		ORDER BY br.id DESC
	`

	var borrowings []Borrowing
	if err := r.db.Select(&borrowings, query); err != nil {
		return nil, err
	}

	return borrowings, nil
}

func (r *Repository) FindByID(id int64) (*Borrowing, error) {
	query := borrowingDetailQuery + `
		WHERE br.id = $1
	`

	var borrowing Borrowing
	if err := r.db.Get(&borrowing, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBorrowingNotFound
		}
		return nil, err
	}

	return &borrowing, nil
}

func (r *Repository) Create(params CreateBorrowingParams) (*Borrowing, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err := r.ensureMemberExists(tx, params.MemberID); err != nil {
		return nil, err
	}

	if err := r.lockBookAndCheckStock(tx, params.BookID); err != nil {
		return nil, err
	}

	var created struct {
		ID int64 `db:"id"`
	}

	insertQuery := `
		INSERT INTO borrowings (
			member_id,
			book_id,
			due_date,
			status
		)
		VALUES ($1, $2, $3, 'borrowed')
		RETURNING id
	`

	if err := tx.Get(&created, insertQuery, params.MemberID, params.BookID, params.DueDate); err != nil {
		return nil, err
	}

	updateStockQuery := `
		UPDATE books
		SET
			stock = stock - 1,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	if _, err := tx.Exec(updateStockQuery, params.BookID); err != nil {
		return nil, err
	}

	borrowing, err := r.findByIDTx(tx, created.ID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return borrowing, nil
}

func (r *Repository) Return(id int64) (*Borrowing, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var current struct {
		ID     int64  `db:"id"`
		BookID int64  `db:"book_id"`
		Status string `db:"status"`
	}

	lockBorrowingQuery := `
		SELECT
			id,
			book_id,
			status
		FROM borrowings
		WHERE id = $1
		FOR UPDATE
	`

	if err := tx.Get(&current, lockBorrowingQuery, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBorrowingNotFound
		}
		return nil, err
	}

	if current.Status == "returned" {
		return nil, ErrBorrowingAlreadyReturned
	}

	updateBorrowingQuery := `
		UPDATE borrowings
		SET
			status = 'returned',
			returned_at = CURRENT_TIMESTAMP,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	if _, err := tx.Exec(updateBorrowingQuery, id); err != nil {
		return nil, err
	}

	updateBookStockQuery := `
		UPDATE books
		SET
			stock = stock + 1,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	if _, err := tx.Exec(updateBookStockQuery, current.BookID); err != nil {
		return nil, err
	}

	borrowing, err := r.findByIDTx(tx, id)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return borrowing, nil
}

func (r *Repository) ensureMemberExists(tx *sqlx.Tx, memberID int64) error {
	var exists bool

	query := `SELECT EXISTS (SELECT 1 FROM members WHERE id = $1)`

	if err := tx.Get(&exists, query, memberID); err != nil {
		return err
	}

	if !exists {
		return ErrMemberNotFound
	}

	return nil
}

func (r *Repository) lockBookAndCheckStock(tx *sqlx.Tx, bookID int64) error {
	var book struct {
		ID    int64 `db:"id"`
		Stock int   `db:"stock"`
	}

	query := `
		SELECT
			id,
			stock
		FROM books
		WHERE id = $1
		FOR UPDATE
	`

	if err := tx.Get(&book, query, bookID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrBookNotFound
		}
		return err
	}

	if book.Stock <= 0 {
		return ErrBookOutOfStock
	}

	return nil
}

func (r *Repository) findByIDTx(tx *sqlx.Tx, id int64) (*Borrowing, error) {
	query := borrowingDetailQuery + `
		WHERE br.id = $1
	`

	var borrowing Borrowing
	if err := tx.Get(&borrowing, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBorrowingNotFound
		}
		return nil, err
	}

	return &borrowing, nil
}