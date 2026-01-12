package models

import (
	"time"

	"github.com/lightsaid/ebook/internal/types"
	"github.com/lightsaid/gotk"
)

type Publisher struct {
	ID            uint64       `db:"id" json:"id"`
	PublisherName string       `db:"publisher_name" json:"publisherName"`
	CreatedAt     types.GxTime `db:"created_at" json:"createdAt" swaggertype:"string"`
	UpdatedAt     types.GxTime `db:"updated_at" json:"updatedAt" swaggertype:"string"`
	DeletedAt     *time.Time   `db:"deleted_at" json:"-"`
}

func (p Publisher) Verifiy(v *gotk.Validator) {
	v.Check(p.PublisherName != "", "publisherName", "出版社名称必填")
}
