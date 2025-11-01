package models

import (
	"time"

	"github.com/lightsaid/ebook/internal/types"
)

type Publisher struct {
	ID            uint64       `db:"id" json:"id"`
	PublisherName string       `db:"publisher_name" json:"publisherName"`
	CreatedAt     types.GxTime `db:"created_at" json:"createdAt"`
	UpdatedAt     types.GxTime `db:"updated_at" json:"UpdatedAt"`
	DeletedAt     *time.Time   `db:"deleted_at" json:"-"`
}
