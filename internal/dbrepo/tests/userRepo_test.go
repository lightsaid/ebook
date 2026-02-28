package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/lightsaid/ebook/internal/dbrepo"
	"github.com/lightsaid/ebook/internal/models"
	"github.com/lightsaid/ebook/pkg/random"
	"github.com/stretchr/testify/require"
)

func randomUser() models.User {
	var email = random.RandomString(8) + "@" + random.RandomString(3)
	email += "." + random.RandomString(2)
	return models.User{
		Email:    email,
		Avatar:   "http://" + random.RandomString(10),
		Nickname: random.RandomString(4),
		Password: random.RandomString(8),
	}
}

func getUserByUq(t *testing.T, uq dbrepo.UserUq) *models.User {
	user, err := tRepo.UserRepo.GetByUqField(context.TODO(), uq)
	require.NoError(t, err)
	require.NotNil(t, user)
	return user
}

func createUser(t *testing.T) *models.User {
	user := randomUser()
	newID, err := tRepo.UserRepo.Create(context.TODO(), &user)
	require.NoError(t, err)
	require.True(t, newID > 0)

	user2 := getUserByUq(t, dbrepo.UserUq{ID: newID})
	require.Equal(t, user.Nickname, user2.Nickname)
	require.Equal(t, user.Email, user2.Email)
	require.Equal(t, user.Avatar, user2.Avatar)
	require.WithinDuration(t, time.Now(), user2.CreatedAt.Time, time.Second*3)
	require.WithinDuration(t, time.Now(), user2.UpdatedAt.Time, time.Second*3)

	return user2
}

func TestCreateUser(t *testing.T) {
	_ = createUser(t)
}

func TestUpdateUser(t *testing.T) {
	u := createUser(t)
	u2 := randomUser()

	u.Avatar = u2.Avatar
	u.Nickname = u2.Nickname
	u.Role = 1

	err := tRepo.UserRepo.Update(context.TODO(), u)
	require.NoError(t, err)

	u3, err := tRepo.UserRepo.GetByUqField(context.TODO(), dbrepo.UserUq{ID: u.ID})
	require.NoError(t, err)

	require.Equal(t, u3.Avatar, u.Avatar)
	require.Equal(t, u3.Nickname, u.Nickname)
	require.Equal(t, u3.Role, u.Role)
}

func TestListUser(t *testing.T) {
	f := dbrepo.Filters{
		PageNum:  1,
		PageSize: 10,
	}
	res, err := tRepo.UserRepo.List(context.TODO(), f)
	require.NoError(t, err)
	require.NotNil(t, res)

	by, _ := json.MarshalIndent(res, "", " ")
	fmt.Println(string(by))
}

// TODO: TEST crud
