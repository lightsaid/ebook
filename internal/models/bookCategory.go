package models

type BookCategory struct {
	BookID     uint64 `db:"book_id" json:"bookId"`
	CategoryID uint64 `db:"category_id" json:"categoryId"`
}
