package models

import (
	"regexp"
	"strings"
	"time"

	"github.com/lightsaid/ebook/internal/types"
	"github.com/lightsaid/gotk"
)

var (
	isbn13 = regexp.MustCompile("^(?:[0-9]{13})$")
	isbn10 = regexp.MustCompile("^(?:[0-9]{9}X|[0-9]{10})$")
)

type Book struct {
	ID          uint64       `db:"id" json:"id"`
	ISBN        string       `db:"isbn" json:"isbn"`
	Title       string       `db:"title" json:"title"`
	Subtitle    string       `db:"subtitle" json:"subtitle"`
	AuthorID    uint64       `db:"author_id" json:"authorId"`
	CoverUrl    string       `db:"cover_url" json:"coverUrl"`
	PublisherID uint64       `db:"publisher_id" json:"publisherId"`
	Pubdate     types.GxTime `db:"pubdate" json:"pubdate"`
	// 价格,单位分
	Price uint `db:"price" json:"price"`
	// 0-下架,1-上架
	Status int `db:"status" json:"status"`
	//1-电子书,2-实体,3-电子书+实体
	Type        int          `db:"type" json:"type"`
	Stock       uint         `db:"stock" json:"stock"`
	SourceUrl   string       `db:"source_url" json:"sourceUrl"`
	Description string       `db:"description" json:"description"`
	Version     string       `db:"version" json:"version"`
	CreatedAt   types.GxTime `db:"created_at" json:"createdAt"`
	UpdatedAt   types.GxTime `db:"updated_at" json:"updatedAt"`
	DeletedAt   *time.Time   `db:"deleted_at" json:"-"`

	Author     *Author     `json:"author"`
	Publisher  *Publisher  `json:"publisher"`
	Categories []*Category `json:"categories"`
}

type SQLBoook struct {
	Book
	// 查询一本图书对应所有分类的sql字段映射
	CategoryJSON string `db:"category_json"`
}

// IsISBN 检查是否是 isbn
func IsISBN(isbn string) bool {
	if isbn10.MatchString(isbn) || isbn13.MatchString(isbn) {
		return true
	}
	return false
}

// Verifiy 实现validator.Verifiyer校验接口
func (b Book) Verifiy(v *gotk.Validator) {
	b.Title = strings.Trim(b.Title, "")
	b.Subtitle = strings.Trim(b.Subtitle, "")

	v.Check(b.Title != "", "title", "标题不能为空")
	v.Check(b.Subtitle != "", "subtitle", "副标题不能为空")
	v.Check(b.ISBN != "", "isbn", "isbn不能为空")
	v.Check(IsISBN(b.ISBN), "isbn", "请输入合法的ISBN")
	v.Check(b.CoverUrl != "", "coverUrl", "请上传封面图")
	v.Check(b.AuthorID > 0, "authorId", "请选择作者")
	v.Check(b.PublisherID > 0, "publisherID", "请选择出版社")
	v.Check(!b.Pubdate.IsZero(), "pubdate", "出版日期必填")
	v.Check(gotk.OneOf(b.Status, 0, 1), "status", "状态: 0-下架,1-上架")
	v.Check(gotk.OneOf(b.Type, 1, 2, 3), "status", "类型: 1-电子书,2-实体,3-电子书+实体")
}
