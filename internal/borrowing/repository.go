package borrowing

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

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
		(
			br.status = 'borrowed'
			AND br.due_date IS NOT NULL
			AND br.due_date < CURRENT_DATE
		) AS is_overdue,
		CASE
			WHEN br.status = 'borrowed'
				AND br.due_date IS NOT NULL
				AND br.due_date < CURRENT_DATE
			THEN CURRENT_DATE - br.due_date
			ELSE 0
		END AS days_overdue,
		br.created_at,
		br.updated_at,
		br.deleted_at,
		m.name AS member_name,
		m.email AS member_email,
		b.title AS book_title,
		b.author AS book_author
	FROM borrowings br
	JOIN members m ON m.id = br.member_id AND m.deleted_at IS NULL
	JOIN books b ON b.id = br.book_id AND b.deleted_at IS NULL
`

func (r *Repository) FindAll(limit, offset int, filter BorrowingFilter) ([]Borrowing, error) {
	filterQuery, args := buildBorrowingFilterQuery(filter)

	query := borrowingDetailQuery + filterQuery + fmt.Sprintf(`
		ORDER BY br.id DESC
		LIMIT $%d OFFSET $%d
	`, len(args)+1, len(args)+2)

	args = append(args, limit, offset)

	var borrowings []Borrowing
	if err := r.db.Select(&borrowings, query, args...); err != nil {
		return nil, err
	}

	return borrowings, nil
}

func (r *Repository) Count(filter BorrowingFilter) (int64, error) {
	filterQuery, args := buildBorrowingFilterQuery(filter)

	var count int64
	query := `
		SELECT COUNT(*)
		FROM borrowings br
		JOIN members m ON m.id = br.member_id AND m.deleted_at IS NULL
		JOIN books b ON b.id = br.book_id AND b.deleted_at IS NULL
	` + filterQuery
	if err := r.db.Get(&count, query, args...); err != nil {
		return 0, err
	}

	return count, nil
}

func (r *Repository) FindByID(id int64) (*Borrowing, error) {
	query := borrowingDetailQuery + `
		WHERE br.id = $1
			AND br.deleted_at IS NULL
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
			AND deleted_at IS NULL
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
			AND deleted_at IS NULL
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

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM members
			WHERE id = $1
				AND deleted_at IS NULL
		)
	`

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
			AND deleted_at IS NULL
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
			AND br.deleted_at IS NULL
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

func (r *Repository) Delete(id int64) error {
	query := `
		UPDATE borrowings
		SET
			deleted_at = CURRENT_TIMESTAMP,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
			AND deleted_at IS NULL
			AND status = 'returned'
	`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrBorrowingNotFound
	}

	return nil
}

func buildBorrowingFilterQuery(filter BorrowingFilter) (string, []any) {
	conditions := []string{"br.deleted_at IS NULL"}
	args := []any{}

	addCondition := func(condition string, value any) {
		args = append(args, value)
		placeholder := fmt.Sprintf("$%d", len(args))
		conditions = append(conditions, fmt.Sprintf(condition, placeholder))
	}

	if filter.Status != "" {
		addCondition("br.status = %s", filter.Status)
	}

	if filter.Overdue != nil {
		overdueCondition := `
			br.status = 'borrowed'
			AND br.due_date IS NOT NULL
			AND br.due_date < CURRENT_DATE
		`
		if *filter.Overdue {
			conditions = append(conditions, overdueCondition)
		} else {
			conditions = append(conditions, "NOT ("+overdueCondition+")")
		}
	}

	if filter.MemberID != nil {
		addCondition("br.member_id = %s", *filter.MemberID)
	}

	if filter.BookID != nil {
		addCondition("br.book_id = %s", *filter.BookID)
	}

	return " WHERE " + strings.Join(conditions, " AND "), args
}
