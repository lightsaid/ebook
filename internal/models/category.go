package models

import "time"

type Category struct {
	ID           uint64     `db:"id" json:"id"`
	CategoryName string     `db:"category_name" json:"categoryName"`
	Icon         string     `db:"icon" json:"icon"`
	Created_at   time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time  `db:"updated_at" json:"UpdatedAt"`
	DeletedAt    *time.Time `db:"deleted_at" json:"-"`
}

type bookCategory struct {
	BookID     uint64 `db:"book_id" json:"bookId"`
	CategoryID uint64 `db:"icategory_idd" json:"categoryId"`
}
