package models

import "time"

type Publisher struct {
	ID            uint64     `db:"id" json:"id"`
	PublisherName string     `db:"publisher_name" json:"publisherName"`
	CreatedAt     time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt     time.Time  `db:"updated_at" json:"UpdatedAt"`
	DeletedAt     *time.Time `db:"deleted_at" json:"-"`
}
