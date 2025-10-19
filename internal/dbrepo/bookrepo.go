package dbrepo

import (
	"context"

	"github.com/lightsaid/ebook/internal/models"
)

type BookRepo interface {
	Create(book *models.Book) (uint64, error)
	CreateTx(book *models.Book) (uint64, error)
	Get(id uint64) (*models.Book, error)
	Update(book *models.Book) error   // 仅更新books表
	UpdateTx(book *models.Book) error // 更新图书和与之关联的分类、出版社、作者
	List(limit, offset int) ([]*models.Book, error)
	ListByCategory(categoryID uint64) ([]*models.Book, error)
	Delete(id uint64) error
}

var _ BookRepo = (*bookRepo)(nil)

type bookRepo struct {
	DB Queryable
}

func NewBookRepo(db Queryable) *bookRepo {
	var repo = &bookRepo{
		DB: db,
	}
	return repo
}

func bookFieldToSQLArgs(book *models.Book) map[string]any {
	args := map[string]any{
		"isbn":         book.ISBN,
		"title":        book.Title,
		"subtitle":     book.Subtitle,
		"author_id":    book.AuthorID,
		"cover_url":    book.CoverUrl,
		"publisher_id": book.PublisherID,
		"pubdate":      book.Pubdate,
		"price":        book.Price,
		"status":       book.Status,
		"type":         book.Type,
		"stock":        book.Stock,
		"source_url":   book.SourceUrl,
		"description":  book.Description,
	}

	if book.ID > 0 {
		args["id"] = book.ID
	}

	return args
}

func (r *bookRepo) Create(book *models.Book) (uint64, error) {
	sql := `
	insert into books set 
		isbn=:isbn, 
		title=:title, 
		subtitle=:subtitle,
		author_id=:author_id,
		cover_url=:cover_url, 
		publisher_id=:publisher_id, 
		pubdate=:pubdate,
		price=:price, 
		status=:status, 
		type=:type,
		stock=:stock,
		source_url=:source_url,
		description=:description
	`

	ctx, cancel := makeCtx()
	defer cancel()
	args := bookFieldToSQLArgs(book)
	result, err := r.DB.NamedExecContext(ctx, sql, args)
	return insertErrorHandler(result, err)
}

func (r *bookRepo) CreateTx(book *models.Book) (uint64, error) {
	return 0, nil
}

func (r *bookRepo) Get(id uint64) (book *models.Book, err error) {
	sql := "select * from books where id=? and deleted_at is null;"
	book = new(models.Book)
	err = r.DB.Get(book, sql, id)
	return book, err
}

func (r *bookRepo) Update(book *models.Book) error {
	sql := `
	update books set 
		isbn=:isbn, 
		title=:title, 
		subtitle=:subtitle,
		author_id=:author_id,
		cover_url=:cover_url, 
		publisher_id=:publisher_id, 
		pubdate=:pubdate,
		price=:price, 
		status=:status, 
		type=:type,
		stock=:stock,
		source_url=:source_url,
		description=:description
	where id=:id and deleted_at is null;
	`
	ctx, cancel := makeCtx()
	defer cancel()
	args := bookFieldToSQLArgs(book)
	return updateErrorHandler(r.DB.NamedExecContext(ctx, sql, args))
}

func (r *bookRepo) UpdateTx(book *models.Book) error {
	execTx(context.Background(), r.DB, func(r Repository) error {
		// r.BookRepo.Create(nil)
		return nil
	})
	return nil
}

func (r *bookRepo) List(limit, offset int) ([]*models.Book, error) {
	sql := `select * from books where deleted_at is null order by id desc limit ? offset ?`
	ctx, cancel := makeCtx()
	defer cancel()
	list := make([]*models.Book, 0, limit)
	err := r.DB.SelectContext(ctx, &list, sql, limit, offset)
	return list, err
}

func (r *bookRepo) ListByCategory(categoryID uint64) (list []*models.Book, err error) {
	sql := `
		select b.*, author_name, publisher_name from books b 
		left join book_categories bc on b.id = bc.book_id
		left join category c on c.id = bc.category_id
		left join author a on a.id = b.author_id
		left join publisher p on p.id = b.publisher_id
		where bc.category_id = ? and b.deleted_at is null;`
	list = make([]*models.Book, 0, 10)
	err = r.DB.Select(&list, sql, categoryID)
	return list, err
}

func (r *bookRepo) Delete(id uint64) error {
	sql := `update books set deleted_at=now() where id=?`
	return updateErrorHandler(r.DB.Exec(sql, id))
}
