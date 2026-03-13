package book

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindAll() ([]Book, error) {
	query := `
		SELECT
			id,
			title,
			author,
			isbn,
			published_year,
			stock,
			created_at,
			updated_at
		FROM books
		ORDER BY id DESC
	`

	var books []Book
	if err := r.db.Select(&books, query); err != nil {
		return nil, err
	}

	return books, nil
}

func (r *Repository) FindByID(id int64) (*Book, error) {
	query := `
		SELECT
			id,
			title,
			author,
			isbn,
			published_year,
			stock,
			created_at,
			updated_at
		FROM books
		WHERE id = $1
	`

	var book Book
	if err := r.db.Get(&book, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBookNotFound
		}
		return nil, err
	}

	return &book, nil
}

func (r *Repository) Create(req CreateBookRequest) (*Book, error) {
	query := `
		INSERT INTO books (
			title,
			author,
			isbn,
			published_year,
			stock
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING
			id,
			title,
			author,
			isbn,
			published_year,
			stock,
			created_at,
			updated_at
	`

	var book Book
	err := r.db.Get(
		&book,
		query,
		req.Title,
		req.Author,
		req.ISBN,
		req.PublishedYear,
		req.Stock,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.SQLState() == "23505" {
			return nil, ErrDuplicateISBN
		}
		return nil, err
	}

	return &book, nil
}

func (r *Repository) Update(id int64, req UpdateBookRequest) (*Book, error) {
	query := `
		UPDATE books
		SET
			title = $1,
			author = $2,
			isbn = $3,
			published_year = $4,
			stock = $5,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
		RETURNING
			id,
			title,
			author,
			isbn,
			published_year,
			stock,
			created_at,
			updated_at
	`

	var book Book
	err := r.db.Get(
		&book,
		query,
		req.Title,
		req.Author,
		req.ISBN,
		req.PublishedYear,
		req.Stock,
		id,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBookNotFound
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.SQLState() == "23505" {
			return nil, ErrDuplicateISBN
		}

		return nil, err
	}

	return &book, nil
}

func (r *Repository) Delete(id int64) error {
	query := `DELETE FROM books WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrBookNotFound
	}

	return nil
}