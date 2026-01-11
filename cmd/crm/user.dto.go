package main

import (
	"regexp"
	"strings"

	"github.com/lightsaid/gotk"
)

var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *SignInRequest) Verifiy(v *gotk.Validator) {
	u.Password = strings.TrimSpace(u.Password)
	u.Email = strings.TrimSpace(u.Email)
	v.Check(u.Email != "", "email", "邮箱不能为空")
	v.Check(gotk.Matches(u.Email, EmailRX), "email", "邮箱地址格式不正确")
	v.Check(u.Password != "", "password", "密码不能为空")
	v.Check(len([]rune(u.Password)) >= 6, "password", "密码长度必须>=6")
}

type RestPasswordRequest struct {
	ID               int64  `json:"id"`
	Password         string `json:"password"`
	NewPassword      string `json:"newPassword"`
	AgainNewPassword string `json:"againNewPassword"`
}

func (u *RestPasswordRequest) Verifiy(v *gotk.Validator) {
	u.Password = strings.TrimSpace(u.Password)
	u.NewPassword = strings.TrimSpace(u.NewPassword)
	u.AgainNewPassword = strings.TrimSpace(u.AgainNewPassword)

	v.Check(u.ID >= 0, "id", "请提供用户id")
	v.Check(u.Password != "", "password", "请输入旧密码")
	v.Check(u.NewPassword != "", "newPassword", "请输入新密码")
	v.Check(u.AgainNewPassword != "", "newPassword", "请输入确认密码")
	v.Check(u.NewPassword == u.AgainNewPassword, "newPassword", "两次密码不一致")
	v.Check(len([]rune(u.NewPassword)) >= 6, "newPassword", "密码长度必须>=6")
}

type UpdateProfileRequest struct {
	Nickname string `db:"nickname" json:"nickname"`
	Avatar   string `db:"avatar" json:"avatar"`
}

func (u *UpdateProfileRequest) Verifiy(v *gotk.Validator) {
	u.Avatar = strings.TrimSpace(u.Avatar)
	u.Nickname = strings.TrimSpace(u.Nickname)
	if u.Avatar == "" && u.Nickname == "" {
		v.AddError("nickname", "请填写用户昵称或头像地址")
	}
}

type RenewAccessTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func (u *RenewAccessTokenRequest) Verifiy(v *gotk.Validator) {
	v.Check(u.RefreshToken != "", "refreshToken", "请提供令牌")
}
