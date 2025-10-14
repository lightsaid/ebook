package tests

import (
	"testing"
	"time"

	"github.com/lightsaid/ebook/internal/models"
	"github.com/lightsaid/ebook/pkg/random"
	"github.com/stretchr/testify/require"
)

func createAuthor(t *testing.T) *models.Author {
	var name = random.RandomString(6)
	newID, err := tRepo.AuthorRepo.Create(name)
	require.NoError(t, err)
	require.True(t, newID > 0)

	author, err := tRepo.AuthorRepo.Get(newID)
	require.NoError(t, err)
	require.Equal(t, name, author.AuthorName)
	require.WithinDuration(t, author.CreatedAt, time.Now(), time.Second)
	require.WithinDuration(t, author.UpdatedAt, time.Now(), time.Second)

	return author
}

func TestCreateAuthor(t *testing.T) {
	_ = createAuthor(t)
}

func TestUpdateAuthor(t *testing.T) {
	a := createAuthor(t)
	newName := random.RandomString(7)
	time.Sleep(time.Second * 3)
	err := tRepo.AuthorRepo.Update(a.ID, newName)
	require.NoError(t, err)
	a2, err := tRepo.AuthorRepo.Get(a.ID)
	require.NoError(t, err)
	require.Equal(t, a.ID, a2.ID)
	require.Equal(t, newName, a2.AuthorName)
	require.WithinDuration(t, a2.UpdatedAt, time.Now(), time.Second)
}
