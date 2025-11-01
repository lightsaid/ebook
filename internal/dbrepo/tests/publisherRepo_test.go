package tests

import (
	"testing"
	"time"

	"github.com/lightsaid/ebook/internal/models"
	"github.com/lightsaid/ebook/pkg/random"
	"github.com/stretchr/testify/require"
)

func createPublisher(t *testing.T) *models.Publisher {
	name := random.RandomString(12)
	id, err := tRepo.PublisherRepo.Create(name)
	require.NoError(t, err)
	require.True(t, id > 0)

	p1, err := tRepo.PublisherRepo.Get(id)
	require.NoError(t, err)
	require.NotEmpty(t, p1)
	require.Equal(t, id, p1.ID)
	require.Equal(t, p1.PublisherName, name)
	require.WithinDuration(t, time.Now(), p1.CreatedAt.Time, time.Second)
	require.WithinDuration(t, time.Now(), p1.UpdatedAt.Time, time.Second)

	return p1
}

func TestCreatePublisher(t *testing.T) {
	_ = createPublisher(t)
}
