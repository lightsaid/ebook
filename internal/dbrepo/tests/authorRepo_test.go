package tests

import (
	"context"
	"testing"
	"time"

	"github.com/lightsaid/ebook/internal/dbrepo"
	"github.com/lightsaid/ebook/internal/models"
	"github.com/lightsaid/ebook/pkg/random"
	"github.com/stretchr/testify/require"
)

func createAuthor(t *testing.T) *models.Author {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var name = random.RandomString(6)
	newID, err := tRepo.AuthorRepo.Create(ctx, name)
	require.NoError(t, err)
	require.True(t, newID > 0)

	author, err := tRepo.AuthorRepo.Get(ctx, newID)
	require.NoError(t, err)
	require.Equal(t, name, author.AuthorName)
	require.WithinDuration(t, author.CreatedAt.Time, time.Now(), time.Second*2)
	require.WithinDuration(t, author.UpdatedAt.Time, time.Now(), time.Second*2)

	return author
}

func TestCreateAuthor(t *testing.T) {
	_ = createAuthor(t)
}

// Get 和 Update 测试
func TestUpdateAuthor(t *testing.T) {
	a := createAuthor(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	newName := random.RandomString(7)
	time.Sleep(time.Second * 3)
	err := tRepo.AuthorRepo.Update(ctx, a.ID, newName)
	require.NoError(t, err)
	a2, err := tRepo.AuthorRepo.Get(ctx, a.ID)
	require.NoError(t, err)
	require.Equal(t, a.ID, a2.ID)
	require.Equal(t, newName, a2.AuthorName)
	require.WithinDuration(t, a2.UpdatedAt.Time, time.Now(), time.Second)
}

// TODO:
func TestAuthorDelete(t *testing.T) {}

func TestAuthorList(t *testing.T) {
	var limit = 10
	for range limit {
		_ = createAuthor(t)
	}

	var f = dbrepo.Filters{
		PageNum:      1,
		PageSize:     limit,
		SortFields:   []string{"-created_at", "id"},
		SortSafelist: []string{"id", "-id", "author_name", "-author_name", "created_at", "-created_at"},
	}
	vo, err := tRepo.AuthorRepo.List(context.TODO(), f)
	list := vo.List.([]*models.Author)
	require.NoError(t, err)
	require.Equal(t, len(list), limit)
	require.True(t, vo.Metadata.TotalCount >= limit)
	require.True(t, vo.Metadata.PageSize == limit)

	// by, _ := json.MarshalIndent(vo, "", " ")
	// fmt.Println(string(by))
}
