package models

import (
	"time"

	"github.com/lightsaid/ebook/internal/types"
)

type Author struct {
	ID         uint64       `db:"id" json:"id"`
	AuthorName string       `db:"author_name" json:"authorName"`
	CreatedAt  types.GxTime `db:"created_at" json:"createdAt,omitempty"`
	UpdatedAt  types.GxTime `db:"updated_at" json:"UpdatedAt,omitempty"`
	DeletedAt  *time.Time   `db:"deleted_at" json:"-"`
}
