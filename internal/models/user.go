package models

import (
	"time"

	"github.com/lightsaid/ebook/internal/types"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint64        `db:"id" json:"id"`
	Email     string        `db:"email" json:"email"`
	Password  string        `db:"password" json:"-"`
	Nickname  string        `db:"nickname" json:"nickname"`
	Avatar    string        `db:"avatar" json:"avatar"`
	Role      int           `db:"role" json:"role"`
	LoginAt   *types.GxTime `db:"login_at" json:"loginAt" swaggertype:"string"`
	LoginIP   *string       `db:"login_ip" json:"loginIp"`
	CreatedAt types.GxTime  `db:"created_at" json:"createdAt" swaggertype:"string"`
	UpdatedAt types.GxTime  `db:"updated_at" json:"updatedAt" swaggertype:"string"`
	DeletedAt *time.Time    `db:"deleted_at" json:"-"`
}

func (u *User) SetHashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	if err != nil {
		return err
	}

	u.Password = string(hash)
	return nil
}

func (u *User) MatchesPassword(plaintext string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plaintext))
	if err != nil {
		return err
	}

	return nil
}
