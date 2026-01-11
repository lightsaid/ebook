package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/go-chi/cors"
	"github.com/lightsaid/ebook/internal/dbrepo"
	"github.com/lightsaid/ebook/internal/models"
	"github.com/lightsaid/ebook/pkg/errs"
	"github.com/lightsaid/gotk"
	"github.com/pascaldekloe/jwt"
	"github.com/redis/go-redis/v9"
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

// CheckToken 校验token是否有效，校验不通过会调用app.FAIL写入响应
func (app *Application) CheckToken(w http.ResponseWriter, r *http.Request, tokenValue string) (*gotk.TokenPayload, uint64, bool) {
	payload, err := app.jwt.ParseToken(tokenValue)
	if err != nil {
		if errors.Is(err, gotk.ErrInvalidToken) {
			app.FAIL(w, r, errs.ErrUnauthorized.WithMessage("无效的认证令牌"))
			return nil, 0, false
		}
		if errors.Is(err, gotk.ErrExpiredToken) {
			app.FAIL(w, r, errs.ErrUnauthorized.WithMessage("认证令牌已过期"))
			return nil, 0, false

		}

		if errors.Is(err, jwt.ErrSigMiss) || errors.Is(err, jwt.ErrUnsecured) {
			app.FAIL(w, r, errs.ErrUnauthorized.WithMessage("认证令牌签名无效"))
			return nil, 0, false

		}

		app.FAIL(w, r, errs.ErrUnauthorized.WithError(err).WithMessage("请提供合法的令牌"))
		return nil, 0, false
	}

	userId, err := strconv.ParseInt(payload.Data, 10, 64)
	if err != nil {
		app.FAIL(w, r, errs.ErrUnauthorized.WithError(err).WithMessage("认证令牌数据无效"))
		return nil, 0, false
	}
	return payload, uint64(userId), true
}

// RequiredAuth 解析Authorization请求头，
// 获取AccessToken，解析user id从redis/mysql获取用户信息并设置到ctx
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
		// 大小写不敏感对比
		if !strings.EqualFold(authType, authTypeBearer) {
			app.FAIL(w, r, errs.ErrUnauthorized.WithMessage("认证令牌类型不对"))
			return
		}
		_, userId, ok := app.CheckToken(w, r, fields[1])
		if !ok {
			return
		}

		// 从 redis 获取
		user, err := cache.UserCache.GetUser(r.Context(), uint64(userId))
		if err != nil && !errors.Is(err, redis.Nil) {
			slog.ErrorContext(r.Context(), "中间件从redis获取用户信息失败", "err", err)
			app.FAIL(w, r, errs.ErrServerError.WithError(err))
			return
		}

		if errors.Is(err, redis.Nil) {
			// 缓存不存在从mysql中获取
			user, err = store.UserRepo.Get(r.Context(), uint64(userId))
			if err != nil {
				slog.ErrorContext(r.Context(), "中间件从mysql获取用户信息失败", "err", err)
				a := dbrepo.ConvertToApiError(err)
				app.FAIL(w, r, a)
				return
			}

			// 回溯/回填信息到redis
			err = cache.UserCache.SaveUser(r.Context(), user)
			if err != nil {
				slog.ErrorContext(r.Context(), "中间件存用户信息到redis失败", "err", err)
				app.FAIL(w, r, errs.ErrServerError.WithError(err))
				return
			}
		}

		// 验证权限相关
		if user.Role != adminRole {
			app.FAIL(w, r, errs.ErrForbidden.WithMessage("请联系管理员"))
			return
		}

		app.SetUserCtx(r, user)

		next.ServeHTTP(w, r)
	})
}

// EnableCORS 支持跨域
func (app *Application) EnableCORS(next http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})(next)
}

// RecoverPanic 错误恢复
func (app *Application) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// 记录panic日志
				slog.ErrorContext(
					r.Context(),
					"panic recovered",
					"error", err,
					"stack", string(debug.Stack()),
				)

				w.Header().Set("Connection", "Close")

				if a, ok := err.(*gotk.ApiError); ok {
					app.FAIL(w, r, a)
					return
				}
				app.FAIL(w, r, errs.ErrServerError.WithError(fmt.Errorf("%v", err)))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// AccessLog TODO: 待完善
func (app *Application) AccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		info := fmt.Sprintf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL)
		slog.InfoContext(r.Context(), info)

		next.ServeHTTP(w, r)
	})
}
