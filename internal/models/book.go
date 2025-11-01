package models

import (
	"time"

	"github.com/lightsaid/ebook/pkg/validator"
)

type Book struct {
	ID          uint64     `db:"id" json:"id"`
	ISBN        string     `db:"isbn" json:"isbn"`
	Title       string     `db:"title" json:"title"`
	Subtitle    string     `db:"subtitle" json:"subtitle"`
	AuthorID    uint64     `db:"author_id" json:"authorId"`
	CoverUrl    string     `db:"cover_url" json:"coverUrl"`
	PublisherID uint64     `db:"publisher_id" json:"publisherId"`
	Pubdate     time.Time  `db:"pubdate" json:"pubdate"`
	Price       uint       `db:"price" json:"price"`
	Status      int        `db:"status" json:"status"`
	Type        int        `db:"type" json:"type"`
	Stock       uint       `db:"stock" json:"stock"`
	SourceUrl   string     `db:"source_url" json:"sourceUrl"`
	Description string     `db:"description" json:"description"`
	Version     string     `db:"version" json:"version"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time  `db:"updated_at" json:"UpdatedAt"`
	DeletedAt   *time.Time `db:"deleted_at" json:"-"`

	Author     *Author     `json:"author"`
	Publisher  *Publisher  `json:"publisher"`
	Categories []*Category `json:"categories"`
}

type SQLBoook struct {
	Book
	// 查询一本图书对应所有分类的sql字段映射
	CategoryJSON string `db:"category_json"`
}

func ValidateBook(v *validator.Validator, book *Book) {
	v.Check(validator.IsISBN(book.ISBN), "isbn", "请输入合法的ISBN")
	v.Check(book.Title != "", "title", "title 不能为空")
}
