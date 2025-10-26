package dbrepo

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
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
	ctx, cancel := makeCtx()
	defer cancel()

	// 不存在分类
	if book.Categories != nil && len(book.Categories) == 0 {
		return r.Create(book)
	}
	var bookID uint64
	// 存在分类
	err := execTx(ctx, r.DB, func(r Repository) error {
		var err error
		bookID, err = r.BookRepo.Create(book)
		if err != nil {
			return err
		}
		list := make([]models.BookCategory, 0, len(book.Categories))
		for _, x := range book.Categories {
			list = append(list, models.BookCategory{BookID: bookID, CategoryID: x.ID})
		}

		return r.BookCategoryRepo.BatchInsert(list)
	})

	return bookID, err
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

	// NOTE: 方式1:
	// return updateErrorHandler(r.DB.NamedExecContext(ctx, query, args))

	// NOTE: 方式2: 方便查看sql和参数，便于日志输出
	query, argv, err := sqlx.Named(sql, args)
	if err != nil {
		log.Println("bookRepo.Update sqlx.Named falil ", err)
		return err
	}

	log.Println("query ", query)
	log.Println("argv  ", argv)

	query = r.DB.Rebind(query)

	return updateErrorHandler(r.DB.ExecContext(ctx, query, argv...))
}

// UpdateTx 通过事务更新图书，同时更新bookCategory表
func (r *bookRepo) UpdateTx(book *models.Book) error {
	err := execTx(context.Background(), r.DB, func(r Repository) error {
		// 更新图书
		err := r.BookRepo.Update(book)
		if err != nil {
			log.Println("[UpdateTx]->[r.BookRepo.Update] fail: ", err)
			return err
		}

		// 删除图书和分类关系（book_categories）
		err = r.BookCategoryRepo.DeleteByBookID(book.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			log.Println("[UpdateTx]->[r.BookCategoryRepo.DeleteByBookID] fail: ", err)
			return err
		}

		if len(book.Categories) == 0 {
			return nil
		}

		// 构建映射关系
		list := make([]models.BookCategory, 0, len(book.Categories))
		for _, x := range book.Categories {
			list = append(list, models.BookCategory{BookID: book.ID, CategoryID: x.ID})
		}

		// 添加新的对应关系
		return r.BookCategoryRepo.BatchInsert(list)
	})
	return err
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
