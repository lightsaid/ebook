package models

import (
	"github.com/lightsaid/ebook/internal/types"
)

type ShoppingCart struct {
	ID        uint64       `db:"id" json:"id"`
	UserID    uint64       `db:"user_id" json:"userId"`
	BookID    uint64       `db:"book_id" json:"bookId"`
	Quantity  uint         `db:"quantity" json:"quantity"`
	CreatedAt types.GxTime `db:"created_at" json:"createdAt" swaggertype:"string"`
	UpdatedAt types.GxTime `db:"updated_at" json:"updatedAt" swaggertype:"string"`
}

type SQLShoppingCart struct {
	Book         Book
	ShoppingCart ShoppingCart
}
