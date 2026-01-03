package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/lightsaid/ebook/internal/dbrepo"
	"github.com/lightsaid/ebook/internal/types"
	"github.com/lightsaid/ebook/pkg/errs"
	"github.com/lightsaid/gotk"
	"github.com/tomasen/realip"
)

// SignIn 处理登录逻辑，role必须为1
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
	if user.Role != 1 {
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

	// TODO: 存储 用户信息到redis
	// user::id

	// 组合返回数据
	data := gotk.Map{"accessToken": aToken, "refreshToken": rToken, "user": user}

	app.SUCC(w, r, data)
}

// UpdateProfile 更新个人信息
func (app *Application) UpdateProfile(w http.ResponseWriter, r *http.Request) {

}

func (app *Application) RenewAccessToken(w http.ResponseWriter, r *http.Request) {

}

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
