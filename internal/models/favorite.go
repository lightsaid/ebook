package models

import "github.com/lightsaid/ebook/internal/types"

type Favorite struct {
	ID        uint64       `db:"id" json:"id"`
	UserID    uint64       `db:"user_id" json:"userId"`
	BookID    uint64       `db:"book_id" json:"bookId"`
	UpdatedAt types.GxTime `db:"updated_at" json:"updatedAt" swaggertype:"string"`
}
