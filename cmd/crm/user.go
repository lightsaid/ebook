package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/lightsaid/ebook/internal/dbcache"
	"github.com/lightsaid/ebook/internal/dbrepo"
	"github.com/lightsaid/ebook/internal/types"
	"github.com/lightsaid/ebook/pkg/errs"
	"github.com/lightsaid/gotk"
	"github.com/tomasen/realip"
)

const (
	adminRole = 1 // 管理员角色
)

// SignIn godoc
//
//	@Summary		管理员登录
//	@Description	管理员登录逻辑处理，role必须为1
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		SignInRequest	true	"登录入参"
//	@Success		200		{object}	int
//	@Router			/v1/signin [post]
func (app *Application) SignIn(w http.ResponseWriter, r *http.Request) {
	var input SignInRequest

	// 读取body参数并校验
	if ok := app.ReadJSONAndCheck(w, r, &input); !ok {
		return
	}

	// 根据email获取用户
	user, err := store.UserRepo.GetByUqField(r.Context(), dbrepo.UserUq{Email: input.Email})
	if err != nil {
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	// 判断权限
	if user.Role != adminRole {
		app.FAIL(w, r, errs.ErrForbidden)
		return
	}

	// 校验用户密码
	err = user.MatchesPassword(input.Password)
	if err != nil {
		app.FAIL(w, r, errs.ErrBadRequest.WithError(err).WithMessage("账户或密码不匹配"))
		return
	}

	// 更新登录信息
	ip := realip.FromRequest(r)
	if ip == "::1" {
		ip = "127.0.0.1"
	}
	user.LoginIP = &ip
	user.LoginAt = &types.GxTime{Time: time.Now()}
	err = store.UserRepo.Update(r.Context(), user)
	if err != nil {
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	// 生成 accessToken 和 refreshToken
	aPayload := gotk.NewTokenPayload(fmt.Sprintf("%d", user.ID), app.config.AccessToknExpires)
	rPayload := gotk.NewTokenPayload(fmt.Sprintf("%d", user.ID), app.config.RefreshToknExpires)
	aToken, err := app.jwt.GenToken(aPayload)
	if err != nil {
		app.FAIL(w, r, errs.ErrServerError.WithError(err))
		return
	}
	rToken, err := app.jwt.GenToken(rPayload)
	if err != nil {
		app.FAIL(w, r, errs.ErrServerError.WithError(err))
		return
	}

	// 缓存用户信息到redis
	err = cache.UserCache.SaveUser(r.Context(), user)
	if err != nil {
		a := dbcache.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	// 组合返回数据
	data := gotk.Map{"accessToken": aToken, "refreshToken": rToken, "user": user}

	app.SUCC(w, r, data)
}

// UpdateProfile 更新个人信息（仅对avatar、nickname修改）
func (app *Application) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	user := app.GetUserCtx(r)
	var input UpdateProfileRequest
	if ok := app.ReadJSONAndCheck(w, r, &input); !ok {
		return
	}
	if input.Avatar != "" {
		user.Avatar = input.Avatar
	}
	if input.Nickname != "" {
		user.Nickname = input.Nickname
	}
	err := store.UserRepo.Update(r.Context(), user)
	if err != nil {
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	app.SUCC(w, r, "修改成功")
}

// RenewAccessToken 根据refreshToken刷新accessToken
func (app *Application) RenewAccessToken(w http.ResponseWriter, r *http.Request) {
	var input RenewAccessTokenRequest
	if ok := app.ReadJSONAndCheck(w, r, &input); !ok {
		return
	}
	_, userId, ok := app.CheckToken(w, r, input.RefreshToken)
	if !ok {
		return
	}
	idText := strconv.Itoa(int(userId))
	payload := gotk.NewTokenPayload(idText, app.config.JWTConfig.AccessToknExpires)
	token, err := app.jwt.GenToken(payload)
	if err != nil {
		app.FAIL(w, r, errs.ErrServerError.WithError(err))
		return
	}

	app.SUCC(w, r, token)
}

// RestPassword 重置密码，通过旧密码修改
func (app *Application) RestPassword(w http.ResponseWriter, r *http.Request) {
	var input RestPasswordRequest
	if ok := app.ReadJSONAndCheck(w, r, &input); !ok {
		return
	}

	user, err := store.UserRepo.Get(r.Context(), uint64(input.ID))
	if err != nil {
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	err = user.MatchesPassword(user.Password)
	if err != nil {
		app.FAIL(w, r, errs.ErrBadRequest.WithMessage("密码不匹配"))
		return
	}

	user.Password = input.NewPassword

	err = store.UserRepo.Update(r.Context(), user)
	if err != nil {
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	app.SUCC(w, r, "修改密码成功")
}

// GetListUser 获取用户列表
func (app *Application) GetListUser(w http.ResponseWriter, r *http.Request) {
	filter := app.ReadPageQuery(r)
	vo, err := store.UserRepo.List(r.Context(), filter)
	if err != nil {
		a := dbrepo.ConvertToApiError(err)
		app.FAIL(w, r, a)
		return
	}

	app.SUCC(w, r, vo)
}
