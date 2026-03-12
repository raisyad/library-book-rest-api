package book

import "time"

type Book struct {
	ID            int64     `db:"id" json:"id"`
	Title         string    `db:"title" json:"title"`
	Author        string    `db:"author" json:"author"`
	ISBN          string    `db:"isbn" json:"isbn"`
	PublishedYear *int      `db:"published_year" json:"published_year"`
	Stock         int       `db:"stock" json:"stock"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}
