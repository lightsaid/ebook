package models

import (
	"time"

	"github.com/lightsaid/ebook/internal/types"
	"github.com/lightsaid/gotk"
)

type Author struct {
	ID         uint64       `db:"id" json:"id"`
	AuthorName string       `db:"author_name" json:"authorName"`
	CreatedAt  types.GxTime `db:"created_at" json:"createdAt" swaggertype:"string"`
	UpdatedAt  types.GxTime `db:"updated_at" json:"updatedAt" swaggertype:"string"`
	DeletedAt  *time.Time   `db:"deleted_at" json:"-"`
}

func (p Author) Verifiy(v *gotk.Validator) {
	v.Check(p.AuthorName != "", "authorName", "作者名称称必填")
}
