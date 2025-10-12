package models

import "time"

type Author struct {
	ID         uint64     `db:"id" json:"id"`
	AuthorName string     `db:"author_name" json:"authorName"`
	Created_at time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time  `db:"updated_at" json:"UpdatedAt"`
	DeletedAt  *time.Time `db:"deleted_at" json:"-"`
}
