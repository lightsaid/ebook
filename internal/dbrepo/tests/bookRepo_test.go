package tests

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/lightsaid/ebook/internal/models"
	"github.com/lightsaid/ebook/pkg/random"
	"github.com/stretchr/testify/require"
)

func createBook(t *testing.T) *models.Book {
	c := createAuthor(t)
	p := createPublisher(t)
	var status = 0
	if random.RandomInt(1, 10) > 5 {
		status = 1
	}
	b1 := &models.Book{
		ISBN:        random.RandomString(11),
		Title:       random.RandomString(8),
		Subtitle:    random.RandomString(20),
		AuthorID:    c.ID,
		CoverUrl:    random.RandomString(32),
		PublisherID: p.ID,
		Pubdate:     time.Now(),
		Price:       uint(random.RandomInt(100, 300)),
		Status:      status,
		Type:        random.RandomInt(1, 3),
		Stock:       uint(random.RandomInt(100, 3000)),
		SourceUrl:   random.RandomString(32),
		Description: random.RandomString(64),
	}

	id, err := tRepo.BookRepo.Create(b1)
	require.NoError(t, err)
	require.True(t, id > 0)

	b2, err := tRepo.BookRepo.Get(id)
	require.NoError(t, err)
	require.NotEmpty(t, b2)
	require.True(t, b2.ID > 0)

	return b2
}

func TestCreateBook(t *testing.T) {
	_ = createBook(t)
}

func TestUpdateBook(t *testing.T) {
	b1 := createBook(t)
	c1 := createAuthor(t)
	p1 := createPublisher(t)

	var status = 0
	if random.RandomInt(1, 10) > 5 {
		status = 1
	}

	b1.ISBN = random.RandomString(11)
	b1.Title = random.RandomString(8)
	b1.Subtitle = random.RandomString(20)
	b1.AuthorID = c1.ID
	b1.CoverUrl = random.RandomString(32)
	b1.PublisherID = p1.ID
	b1.Pubdate = time.Now()
	b1.Price = uint(random.RandomInt(100, 300))
	b1.Status = status
	b1.Type = random.RandomInt(1, 3)
	b1.Stock = uint(random.RandomInt(100, 3000))
	b1.SourceUrl = random.RandomString(32)
	b1.Description = random.RandomString(64)

	time.Sleep(time.Second * 2)

	err := tRepo.BookRepo.Update(b1)
	require.NoError(t, err)
}

func TestDeleted(t *testing.T) {
	b := createBook(t)
	require.Empty(t, b.DeletedAt)

	err := tRepo.BookRepo.Delete(b.ID)
	require.NoError(t, err)

	_, err = tRepo.BookRepo.Get(b.ID)
	require.ErrorIs(t, err, sql.ErrNoRows)
}

func TestListBook(t *testing.T) {
	var size = 10
	for range size {
		_ = createBook(t)
	}

	list, err := tRepo.BookRepo.List(10, 0)
	require.NoError(t, err)
	require.True(t, len(list) == 10)

	by, err := json.MarshalIndent(list, "", "\t")
	require.NoError(t, err)
	fmt.Println(string(by))
}
