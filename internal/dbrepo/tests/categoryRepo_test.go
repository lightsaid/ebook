package tests

import (
	"context"
	"testing"
	"time"

	"github.com/lightsaid/ebook/internal/models"
	"github.com/lightsaid/ebook/pkg/random"
	"github.com/stretchr/testify/require"
)

func createCategory(t *testing.T) *models.Category {
	var c = models.Category{
		CategoryName: random.RandomString(4),
		Icon:         random.RandomString(10),
		Sort:         random.RandomInt(1, 100),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, err := tRepo.CategoryRepo.Create(ctx, c)
	require.NoError(t, err)
	require.True(t, id > 0)

	c2, err := tRepo.CategoryRepo.Get(ctx, id)

	require.NoError(t, err)
	require.True(t, id == c2.ID)
	require.Equal(t, c.CategoryName, c2.CategoryName)
	require.Equal(t, c.Icon, c2.Icon)
	require.Equal(t, c.Sort, c2.Sort)
	require.WithinDuration(t, time.Now(), c2.CreatedAt.Time, time.Second)
	require.WithinDuration(t, time.Now(), c2.UpdatedAt.Time, time.Second)

	return c2
}

func TestCreateCategory(t *testing.T) {
	_ = createCategory(t)
}

func TestUpdateCategory(t *testing.T) {
	c1 := createCategory(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	c2 := models.Category{
		CategoryName: random.RandomString(4),
		Icon:         random.RandomString(10),
		Sort:         random.RandomInt(1, 100),
	}

	c2.ID = c1.ID

	err := tRepo.CategoryRepo.Update(ctx, c2)
	require.NoError(t, err)

	c3, err := tRepo.CategoryRepo.Get(ctx, c2.ID)
	require.NoError(t, err)

	require.Equal(t, c3.CategoryName, c2.CategoryName)
	require.Equal(t, c3.Icon, c2.Icon)
	require.Equal(t, c3.Sort, c2.Sort)
}
