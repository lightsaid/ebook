package tests

import (
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
