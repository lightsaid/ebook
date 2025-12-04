package models

import (
	"strings"
	"time"

	"github.com/lightsaid/ebook/internal/types"
	"github.com/lightsaid/gotk"
)

type Banner struct {
	ID        uint64       `db:"id" json:"id"`
	Slogan    string       `db:"slogan" json:"slogan"`
	LinkType  int          `db:"link_type" json:"linkType"`
	LinkUrl   string       `db:"link_url" json:"linkUrl"`
	ImageUrl  string       `db:"image_url" json:"imageUrl"`
	Enable    int          `db:"enable" json:"enable"`
	Sort      int          `db:"sort" json:"sort"`
	CreatedAt types.GxTime `db:"created_at" json:"createdAt"`
	UpdatedAt types.GxTime `db:"updated_at" json:"UpdatedAt"`
	DeletedAt *time.Time   `db:"deleted_at" json:"-"`
}

// Verifiy 实现validator.Verifiyer校验接口
func (b Banner) Verifiy(v *gotk.Validator) {
	b.Slogan = strings.Trim(b.Slogan, "")
	v.Check(len(b.Slogan) > 0, "slogan", "标题不能为空")
	v.Check(gotk.OneOf(b.Enable, 0, 1), "enable", "enable仅支持1启用、0停用")
	v.Check(gotk.OneOf(b.LinkType, 0, 1), "enable", "enable仅支持0内链、1外链")
	v.Check(len(b.ImageUrl) > 0, "imageUrl", "请上传图片")
	v.Check(len(b.LinkUrl) > 0, "linkType", "请填写跳转链接")
}
