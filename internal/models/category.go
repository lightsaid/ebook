package models

import (
	"strings"
	"time"

	"github.com/lightsaid/ebook/internal/types"
	"github.com/lightsaid/gotk"
)

type Category struct {
	ID           uint64       `db:"id" json:"id"`
	CategoryName string       `db:"category_name" json:"categoryName"`
	Icon         string       `db:"icon" json:"icon"`
	Sort         int          `db:"sort" json:"sort"`
	CreatedAt    types.GxTime `db:"created_at" json:"createdAt" swaggertype:"string"`
	UpdatedAt    types.GxTime `db:"updated_at" json:"updatedAt" swaggertype:"string"`
	DeletedAt    *time.Time   `db:"deleted_at" json:"-"`
}

// Verifiy 实现validator.Verifiyer校验接口
func (c Category) Verifiy(v *gotk.Validator) {
	c.CategoryName = strings.Trim(c.CategoryName, "")
	v.Check(len(c.CategoryName) > 0, "categoryName", "分类名称不能为空")
}

// SQLBookCategory 方便查询映射, 不存库
type SQLBookCategory struct {
	Category     Category     `db:"category"`
	BookCategory BookCategory `db:"book_category"`
}
