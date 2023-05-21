package dbrepo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/lightsaid/ebook/internal/models"
	"github.com/lightsaid/ebook/pkg/logger"
)

type BookRepo interface {
	Insert(b *models.Book) (uint, error)
	Get(id uint) (*models.Book, error)
	List(page Pager) ([]*models.Book, error)
	Update(id uint, book *models.Book) error
	Del(id uint) error
	BulkInsert(books []*models.Book) (int64, error)
	TestTx(id uint) error
}

var _ BookRepo = (*bookRepo)(nil)

type bookRepo struct {
	dbBase
	DB Queryable
}

func NewBookRepo(db Queryable) *bookRepo {
	return &bookRepo{
		DB: db,
	}
}

func (b *bookRepo) Insert(book *models.Book) (uint, error) {
	insertSQL := `insert into book(
		isbn,
		title,
		poster,
		pages,
		price,
		published_at
	)values(
		?,?,?,?,?,?
	)`

	// ret, err := b.DB.Exec(insertSQL, book.ISBN, book.Title, book.Poster, book.Pages, book.Price, book.PublishedAt)
	// if err != nil {
	// 	return 0, err
	// }

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt, err := b.DB.PrepareContext(ctx, insertSQL)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	ret, err := stmt.Exec(book.ISBN, book.Title, book.Poster, book.Pages, book.Price, book.PublishedAt)
	if err != nil {
		return 0, handleMySQLError(err)
	}

	insertID, err := ret.LastInsertId()
	if err != nil {
		return 0, handleMySQLError(err)
	}

	return uint(insertID), nil
}

func (b *bookRepo) Get(id uint) (*models.Book, error) {
	querySQL := `select 
		id,
		isbn,
		title,
		poster,
		pages,
		price,
		published_at
	from book where id = ?`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt, err := b.DB.PrepareContext(ctx, querySQL)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var book = new(models.Book)
	err = stmt.QueryRow(id).Scan(&book.ID, &book.ISBN, &book.Title, &book.Poster, &book.Pages, &book.Price, &book.PublishedAt)
	return book, handleMySQLError(err)
}

func (b *bookRepo) List(page Pager) ([]*models.Book, error) {
	querySQL := `select 
		id,
		isbn,
		title,
		poster,
		pages,
		price,
		published_at
	from book limit ? offset ?`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt, err := b.DB.PrepareContext(ctx, querySQL)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(page.Limit(), page.Offset())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books = make([]*models.Book, 0)
	for rows.Next() {
		var book = new(models.Book)
		err = rows.Scan(&book.ID, &book.ISBN, &book.Title, &book.Poster, &book.Pages, &book.Price, &book.PublishedAt)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return books, nil
}

func (b *bookRepo) Update(id uint, book *models.Book) error {
	updateSQL := `update book set
		isbn = ?,
		title = ?,
		poster = ?,
		pages = ?,
		price = ?,
		published_at = ?,
		updated_at = now()
		where id = ?
	`

	qBook, err := b.Get(id)
	if err != nil {
		return handleMySQLError(err)
	}

	if book.ISBN != "" {
		qBook.ISBN = book.ISBN
	}

	if book.Title != "" {
		qBook.Title = book.Title
	}

	if book.Poster != "" {
		qBook.Poster = book.Poster
	}

	if book.Pages != 0 {
		qBook.Pages = book.Pages
	}

	if book.Price >= 0 {
		qBook.Price = book.Price
	}

	if !book.PublishedAt.IsZero() {
		qBook.PublishedAt = book.PublishedAt
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt, err := b.DB.PrepareContext(ctx, updateSQL)
	if err != nil {
		return err
	}
	defer stmt.Close()

	ret, err := stmt.Exec(
		qBook.ISBN,
		qBook.Title, qBook.Poster,
		qBook.Pages,
		qBook.Price,
		qBook.PublishedAt,
		qBook.ID,
	)
	if err != nil {
		return err
	}

	aff, err := ret.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w:%v", ErrUpdateFailed, err)
	}

	// 最终更新成功
	if aff > 0 {
		return nil
	}

	return ErrUpdateFailed
}

func (b *bookRepo) Del(id uint) error {
	delSQL := `delete from book where id = ?`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt, err := b.DB.PrepareContext(ctx, delSQL)
	if err != nil {
		return err
	}
	defer stmt.Close()

	ret, err := stmt.Exec(id)
	if err != nil {
		return handleMySQLError(err)
	}

	aff, err := ret.RowsAffected()
	if err != nil {
		return err
	}

	if aff > 0 {
		return nil
	}
	// SQL 执行成功，但是没有影响行数，那就是数据已经不存在了
	return ErrNotFound
}

// BulkInsert 批量插入, 并没有对数量做限制，因此需要注意插入数据的数量，以防超出SQL语句的大小限制
func (b *bookRepo) BulkInsert(books []*models.Book) (int64, error) {
	// 批量插入，采用 insert into tbName(field1, field2) values(?, ?),values(?,?)...
	insertSQL := `insert into book(
		isbn,
		title,
		poster,
		pages,
		price,
		published_at
	) values `

	var values []string
	var params []interface{}
	for _, book := range books {
		values = append(values, "(?,?,?,?,?,?)")
		params = append(params, book.ISBN, book.Title, book.Poster, book.Pages, book.Price, book.PublishedAt)
	}

	// 拼接 values => "insert into ... values (?,?...),(?,?...)"
	insertSQL += strings.Join(values, ",")

	logger.InfoLog.Println("book BulkInsert sql: ", insertSQL)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt, err := b.DB.PrepareContext(ctx, insertSQL)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	ret, err := stmt.Exec(params...)
	if err != nil {
		return 0, handleMySQLError(err)
	}

	aff, err := ret.RowsAffected()
	if err != nil {
		return 0, err
	}
	return aff, nil
}

// TestTx NOTE: 测试事物，获取图书和更新图书，不会修改任何数据
func (b *bookRepo) TestTx(id uint) error {
	err := b.dbBase.execTx(context.Background(), b.DB, func(r *Repository) error {
		// 获取
		book, err := r.Book.Get(id)
		if err != nil {
			return err
		}
		// 更新
		return r.Book.Update(id, book)
	})
	return err
}

// // execTx 执行事务的方法
// func (b *bookRepo) execTx(ctx context.Context, fn func(*Repository) error) error {
// 	// b.DB => sql.DB/sql.Tx
// 	db, ok := b.DB.(*sql.DB)
// 	if !ok {
// 		return errors.New("b.DB not is *sql.DB")
// 	}
// 	// 1. 开启事务
// 	tx, err := db.BeginTx(ctx, nil)
// 	if err != nil {
// 		return err
// 	}
// 	r := NewRepository(tx)
// 	if _, err = r.Book.Get(6); err != nil {
// 		return err
// 	}

// 	q := NewRepository(tx)
// 	if err = fn(q); err != nil {
// 		// 2. 执行不成功，则 Rollback
// 		if rbErr := tx.Rollback(); rbErr != nil {
// 			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
// 		}
// 		return err
// 	}

// 	// 3. 执行成功 Commit
// 	return tx.Commit()
// }
