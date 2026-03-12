package database

import (
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func NewPostgresConnection(dbURL string) (*sqlx.DB, error) {
	if dbURL == "" {
		return nil, fmt.Errorf("database url is required")
	}
	db, err := sqlx.Connect("pgx", dbURL)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	return db, nil
}
