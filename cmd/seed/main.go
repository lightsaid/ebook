package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/lightsaid/ebook/internal/config"
	"github.com/lightsaid/ebook/internal/dbrepo"
	"github.com/lightsaid/ebook/internal/models"
)

var store dbrepo.Repository

func main() {
	var conf config.DbConfig
	config.Load(&conf, "./configs/develop.env")
	db, err := dbrepo.Open(conf)
	if err != nil {
		panic(err)
	}

	store = dbrepo.NewRepository(db)
	createUser()
}

func createUser() {
	var email = "lightsaid@foxmail.com"
	_, err := store.UserRepo.GetByUqField(context.TODO(), dbrepo.UserUq{Email: email})
	if errors.Is(err, sql.ErrNoRows) {
		user := models.User{
			Nickname: "lightsaid",
			Email:    email,
			Password: "123456",
			Avatar:   "http://",
		}

		if err := user.SetHashPassword(); err != nil {
			fmt.Println(err)
			return
		}

		newID, err := store.UserRepo.Create(context.TODO(), &user)
		if err != nil {
			fmt.Println(err)
			return
		}

		updateUser, err := store.UserRepo.Get(context.TODO(), newID)
		if err != nil {
			fmt.Println(err)
			return
		}

		updateUser.Role = 1

		err = store.UserRepo.Update(context.TODO(), updateUser)
		if err != nil {
			fmt.Println(err)
		}
	}

}
