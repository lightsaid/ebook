package models

import (
	"time"

	"github.com/lightsaid/ebook/internal/types"
)

type Banner struct {
	ID        uint64       `db:"id" json:"id"`
	Slogan    string       `db:"slogan" json:"slogan"`
	LinkType  int          `db:"link_type" json:"linkType"`
	LinkURL   int          `db:"link_url" json:"linkUrl"`
	ImageUrl  string       `db:"image_url" json:"imageUrl"`
	Enable    int          `db:"enable" json:"enable"`
	Sort      int          `db:"sort" json:"sort"`
	CreatedAt types.GxTime `db:"created_at" json:"createdAt"`
	UpdatedAt types.GxTime `db:"updated_at" json:"UpdatedAt"`
	DeletedAt *time.Time   `db:"deleted_at" json:"-"`
}
