package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/lightsaid/ebook/internal/dbrepo"
	"github.com/lightsaid/ebook/internal/models"
	"github.com/lightsaid/ebook/pkg/errs"
	"github.com/lightsaid/gotk"
	"github.com/pascaldekloe/jwt"
)

const (
	authTypeBearer = "Bearer"
	authHeaderKey  = "Authorization"
	userCtxKey     = gotk.CtxKey("user")
)

// SetUserCtx 设置用户信息到上下文
func (app *Application) SetUserCtx(r *http.Request, user *models.User) *http.Request {
	ctx := context.WithValue(r.Context(), userCtxKey, user)
	return r.WithContext(ctx)
}

// GetUserCtx 从上下文获取用户信息
func (app *Application) GetUserCtx(r *http.Request) *models.User {
	user, ok := r.Context().Value(userCtxKey).(*models.User)
	if !ok {
		panic(errs.ErrUnauthorized)
	}

	return user
}

// RequiredAuth 解析Authorization请求头，
// 获取AccessToken，解析user id获取用户信息并设置到ctx
func (app *Application) RequiredAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get(authHeaderKey)
		if len(authHeader) == 0 {
			app.FAIL(w, r, errs.ErrUnauthorized.WithMessage("请提供认证令牌"))
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) != 2 {
			app.FAIL(w, r, errs.ErrUnauthorized.WithMessage("认证令牌格式不对"))
			return
		}

		authType := fields[0]
		if authType != authTypeBearer {
			app.FAIL(w, r, errs.ErrUnauthorized.WithMessage("认证令牌类型不对"))
			return
		}
		tokenValue := fields[1]
		payolad, err := app.jwt.ParseToken(tokenValue)
		if err != nil {
			if errors.Is(err, gotk.ErrInvalidToken) {
				app.FAIL(w, r, errs.ErrUnauthorized.WithMessage("无效的认证令牌"))
				return
			}
			if errors.Is(err, gotk.ErrExpiredToken) {
				app.FAIL(w, r, errs.ErrUnauthorized.WithMessage("认证令牌已过期"))
				return
			}

			if errors.Is(err, jwt.ErrSigMiss) || errors.Is(err, jwt.ErrUnsecured) {
				app.FAIL(w, r, errs.ErrUnauthorized.WithMessage("认证令牌签名无效"))
				return
			}

			app.FAIL(w, r, errs.ErrUnauthorized.WithError(err).WithMessage("请提供合法的令牌"))
			return
		}

		userId, err := strconv.ParseInt(payolad.Data, 10, 64)

		// TODO: 从 redis 获取

		user, err := store.UserRepo.Get(r.Context(), uint64(userId))
		if err != nil {
			a := dbrepo.ConvertToApiError(err)
			app.FAIL(w, r, a)
			return
		}

		app.SetUserCtx(r, user)

		next.ServeHTTP(w, r)
	})
}
