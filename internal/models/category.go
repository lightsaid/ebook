package models

import (
	"time"

	"github.com/lightsaid/ebook/internal/types"
)

type Category struct {
	ID           uint64       `db:"id" json:"id"`
	CategoryName string       `db:"category_name" json:"categoryName"`
	Icon         string       `db:"icon" json:"icon"`
	Sort         int          `db:"sort" json:"sort"`
	CreatedAt    types.GxTime `db:"created_at" json:"createdAt"`
	UpdatedAt    types.GxTime `db:"updated_at" json:"UpdatedAt"`
	DeletedAt    *time.Time   `db:"deleted_at" json:"-"`
}

// SQLBookCategory 方便查询映射
type SQLBookCategory struct {
	Category     Category     `db:"category"`
	BookCategory BookCategory `db:"book_category"`
}
