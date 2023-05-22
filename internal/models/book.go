package models

import (
	"time"

	"github.com/lightsaid/ebook/pkg/validator"
)

type Book struct {
	ID          uint      `json:"id"`
	ISBN        string    `json:"isbn"`
	Title       string    `json:"title"`
	Poster      string    `json:"poster"`
	Pages       uint      `json:"pages"`
	Price       float32   `json:"price"`
	PublishedAt time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}

func ValidateBook(v *validator.Validator, book *Book) {
	v.Check(validator.IsISBN(book.ISBN), "isbn", "请输入合法的ISBN")
	v.Check(book.Title != "", "title", "title 不能为空")
	v.Check(book.Price >= 0, "price", "价格不能为负数")
	v.Check(book.Pages > 0, "pages", "页数小于0")
	v.Check(book.Poster != "", "poster", "封面不能为空")
	v.Check(!book.PublishedAt.IsZero(), "published_at", "published_at 不能为空")
}

// func NewBook(isbn, title, poster string, pages uint, price float32, publishedAt time.Time) Book {
// 	return Book{
// 		ISBN:        isbn,
// 		Title:       title,
// 		Poster:      poster,
// 		Pages:       pages,
// 		Price:       price,
// 		PublishedAt: publishedAt,
// 		CreatedAt:   time.Now(),
// 		UpdatedAt:   time.Now(),
// 	}
// }
