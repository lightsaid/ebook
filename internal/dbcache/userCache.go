package dbcache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/lightsaid/ebook/internal/models"
	"github.com/redis/go-redis/v9"
)

type UserCache interface {
	SaveUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, userID uint64) (*models.User, error)
}

type userCache struct {
}

var _ UserCache = (*userCache)(nil)

func NewUserCache() *userCache {
	return &userCache{}
}

// SaveUser 保存一个用户，有效时间5分钟
func (cache *userCache) SaveUser(ctx context.Context, user *models.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	key := userIDKey(user.ID)

	d := 5 * time.Minute

	err = rdb.SetEx(ctx, key, string(data), d).Err()
	if err != nil {
		return err
	}

	return nil
}

// GetUser 获取一个用户
func (cache *userCache) GetUser(ctx context.Context, userID uint64) (*models.User, error) {
	key := userIDKey(userID)
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return &models.User{}, err
	}

	if len(val) == 0 {
		return &models.User{}, redis.Nil
	}

	var user models.User
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		return &models.User{}, err
	}

	return &user, nil
}
