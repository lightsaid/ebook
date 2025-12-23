package models

import (
	"time"

	"github.com/lightsaid/ebook/internal/types"
)

type Order struct {
	ID          uint64        `db:"id" json:"id"`
	OrderNo     uint64        `db:"order_no" json:"orderNo"`
	UserID      uint64        `db:"user_id" json:"userId"`
	OrderStatus int           `db:"order_status" json:"orderStatus"`
	OrderAmount uint          `db:"order_amount" json:"orderAmount"`
	PaidAt      *types.GxTime `db:"paid_at" json:"paidAt"`
	CreatedAt   types.GxTime  `db:"created_at" json:"createdAt"`
	UpdatedAt   types.GxTime  `db:"updated_at" json:"updatedAt"`
	DeletedAt   *time.Time    `db:"deleted_at" json:"-"`
}

type OrderItem struct {
	ID        uint64       `db:"id" json:"id"`
	OrderID   uint64       `db:"order_id" json:"orderId"`
	BookID    uint64       `db:"book_id" json:"bookId"`
	Quantity  uint         `db:"quantity" json:"quantity"`
	UnitPrice uint         `db:"unit_price" json:"unitPrice"`
	CreatedAt types.GxTime `db:"created_at" json:"createdAt"`
	DeletedAt *time.Time   `db:"deleted_at" json:"-"`
}
