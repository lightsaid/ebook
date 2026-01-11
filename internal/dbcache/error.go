package dbcache

import (
	"errors"

	"github.com/lightsaid/ebook/pkg/errs"
	"github.com/lightsaid/gotk"
	"github.com/redis/go-redis/v9"
)

func ConvertToApiError(err error) *gotk.ApiError {
	if errors.Is(err, redis.Nil) {
		return errs.ErrNotFound
	}

	return errs.ErrServerError
}
