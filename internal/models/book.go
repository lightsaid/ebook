package models

import "time"

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

func NewBook(isbn, title, poster string, pages uint, price float32, publishedAt time.Time) Book {
	return Book{
		ISBN:        isbn,
		Title:       title,
		Poster:      poster,
		Pages:       pages,
		Price:       price,
		PublishedAt: publishedAt,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
