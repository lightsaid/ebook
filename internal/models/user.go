package models

import (
	"time"

	"github.com/lightsaid/ebook/internal/types"
)

type User struct {
	ID        uint64        `db:"id" json:"id"`
	Email     string        `db:"email" json:"email"`
	Password  string        `db:"password" json:"-"`
	Nickname  string        `db:"nickname" json:"nickname"`
	Avatar    string        `db:"avatar" json:"avatar"`
	Role      int           `db:"role" json:"role"`
	LoginAt   *types.GxTime `db:"login_at" json:"loginAt"`
	LoginIP   *string       `db:"login_ip" json:"loginIp"`
	CreatedAt types.GxTime  `db:"created_at" json:"createdAt"`
	UpdatedAt types.GxTime  `db:"updated_at" json:"updatedAt"`
	DeletedAt *time.Time    `db:"deleted_at" json:"-"`
}
