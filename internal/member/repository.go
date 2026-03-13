package member

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

func (r *Repository) FindAll() ([]Member, error) {
	query := `
		SELECT
			id,
			name,
			email,
			phone,
			created_at,
			updated_at
		FROM members
		ORDER BY id DESC
	`

	var members []Member
	if err := r.db.Select(&members, query); err != nil {
		return nil, err
	}

	return members, nil
}

func (r *Repository) FindByID(id int64) (*Member, error) {
	query := `
		SELECT
			id,
			name,
			email,
			phone,
			created_at,
			updated_at
		FROM members
		WHERE id = $1
	`

	var member Member
	if err := r.db.Get(&member, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMemberNotFound
		}
		return nil, err
	}

	return &member, nil
}

func (r *Repository) Create(req CreateMemberRequest) (*Member, error) {
	query := `
		INSERT INTO members (
			name,
			email,
			phone
		)
		VALUES ($1, $2, $3)
		RETURNING
			id,
			name,
			email,
			phone,
			created_at,
			updated_at
	`

	var member Member
	err := r.db.Get(
		&member,
		query,
		req.Name,
		req.Email,
		req.Phone,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.SQLState() == "23505" {
			return nil, ErrDuplicateEmail
		}
		return nil, err
	}

	return &member, nil
}

func (r *Repository) Update(id int64, req UpdateMemberRequest) (*Member, error) {
	query := `
		UPDATE members
		SET
			name = $1,
			email = $2,
			phone = $3,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
		RETURNING
			id,
			name,
			email,
			phone,
			created_at,
			updated_at
	`

	var member Member
	err := r.db.Get(
		&member,
		query,
		req.Name,
		req.Email,
		req.Phone,
		id,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMemberNotFound
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.SQLState() == "23505" {
			return nil, ErrDuplicateEmail
		}

		return nil, err
	}

	return &member, nil
}

func (r *Repository) Delete(id int64) error {
	query := `DELETE FROM members WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrMemberNotFound
	}

	return nil
}