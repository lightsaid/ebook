package dbrepo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/lightsaid/ebook/internal/models"
)

type BookRepo interface {
	Create(ctx context.Context, book *models.Book) (uint64, error)
	CreateTx(ctx context.Context, book *models.Book) (uint64, error)
	Get(ctx context.Context, id uint64) (*models.Book, error)
	Update(ctx context.Context, book *models.Book) error   // 仅更新books表
	UpdateTx(ctx context.Context, book *models.Book) error // 更新图书和与之关联的分类、出版社、作者
	List(context.Context, Filters) ([]*models.Book, error)
	ListByCategory(ctx context.Context, categoryID uint64) ([]*models.Book, error)
	ListWithCategory(ctx context.Context, filter Filters) (list []*models.Book, err error)
	Delete(ctx context.Context, id uint64) error
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

func (r *bookRepo) Create(ctx context.Context, book *models.Book) (uint64, error) {
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

	slog.InfoContext(ctx, spaceRex.ReplaceAllString(sql, " "), "args", slog.AnyValue(args))

	result, err := r.DB.NamedExecContext(ctx, sql, args)
	return insertErrorHandler(result, err)
}

func (r *bookRepo) CreateTx(ctx context.Context, book *models.Book) (uint64, error) {
	ctx, cancel := timeoutCtx(ctx)
	defer cancel()

	// 不存在分类
	if book.Categories != nil && len(book.Categories) == 0 {
		return r.Create(ctx, book)
	}
	var bookID uint64
	// 存在分类
	err := execTx(ctx, r.DB, func(r Repository) error {
		var err error
		bookID, err = r.BookRepo.Create(ctx, book)
		if err != nil {
			return err
		}
		list := make([]models.BookCategory, 0, len(book.Categories))
		for _, x := range book.Categories {
			list = append(list, models.BookCategory{BookID: bookID, CategoryID: x.ID})
		}

		return r.BookCategoryRepo.BatchInsert(ctx, list)
	})

	return bookID, err
}

func (r *bookRepo) Get(ctx context.Context, id uint64) (book *models.Book, err error) {
	// NOTE: 使用 group_concat 默认长度是1024个字符
	// 临时增加 GROUP_CONCAT
	// SET SESSION group_concat_max_len = 8192;

	// NOTE: GROUP_CONCAT 的对齐问题
	// 比如下面的 category_ids 可能有3个，但是 category_icons 只有2个，就会空值错位了
	sql := `
	select 
		b.*, 
		a.id as "author.id",
		a.author_name as "author.author_name", 
		p.id as "publisher.id",
		p.publisher_name as "publisher.publisher_name",
		-- group_concat(distinct c.id) as category_ids,
		-- group_concat(distinct c.category_name) as category_names,
		-- group_concat(distinct c.sort) as category_sorts,
		-- group_concat(distinct c.icon) as category_icons

		-- 拼接为JSON字符串,就不存在错位问题
		-- json字段命名需要安装Category tag命名，方便使用json.Unmarshal转换
		group_concat(
		  json_object( 
		    'id', c.id,
		    'categoryName', c.category_name,
		    'icon', IFNULL(c.icon, ''),
		    'sort', IFNULL(c.sort, 0)
		  )
		) as category_json
	from books b
	left join author a on a.id=b.author_id
	left join publisher p on p.id=b.publisher_id
	left join book_categories bc ON b.id = bc.book_id
	left join category c ON bc.category_id = c.id
	where b.id=? and b.deleted_at is null
	
	-- 使用了group_concat，避免查询出错，指定分组
	group by b.id;
	`
	book = new(models.Book)
	queryBook := new(models.SQLBoook)

	ctx, cancel := timeoutCtx(ctx)
	defer cancel()

	slog.InfoContext(ctx, sql, "id", slog.Int64Value(int64(id)))

	err = r.DB.GetContext(ctx, queryBook, sql, id)
	if err != nil {
		slog.ErrorContext(ctx, "r.DB.Get fail ", slog.String("err", err.Error()))
		return book, err
	}

	// 把SQLBoook数据同步book上
	var categories []*models.Category
	err = json.Unmarshal([]byte("["+queryBook.CategoryJSON+"]"), &categories)
	if err != nil {
		slog.ErrorContext(ctx, "bookRepo.Get json.Unmarshal faile ", slog.String("err", err.Error()))
		return book, err
	}

	book = &queryBook.Book
	book.Categories = categories

	// by, _ := json.MarshalIndent(book, "", "\t")
	// fmt.Println(string(by))

	// Get 方法，查询到数据，都是零值问题说明，err 也会说明是那个字段没法映射
	/*
		Get() 只取第一行。
		如果 SQL 结果为空，sqlx.Get() 会返回 sql.ErrNoRows。
		如果列名不匹配，不会报错，只是结构体字段保持零值。
	*/

	return book, err
}

func (r *bookRepo) Update(ctx context.Context, book *models.Book) error {
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
	ctx, cancel := timeoutCtx(ctx)
	defer cancel()
	args := bookFieldToSQLArgs(book)

	// NOTE: 方式1:
	// return updateErrorHandler(r.DB.NamedExecContext(ctx, query, args))

	// NOTE: 方式2: 方便查看sql和参数，便于日志输出
	// query, argv, err := sqlx.Named(sql, args)
	// if err != nil {
	// 	log.Println("bookRepo.Update sqlx.Named falil ", err)
	// 	return err
	// }

	// query = r.DB.Rebind(query)

	query, argv, err := debugSQL(ctx, r.DB, sql, args)
	if err != nil {
		slog.ErrorContext(ctx, "bookRepo.Update sqlx.Named falil ", slog.String("err", err.Error()))
		return err
	}

	return updateErrorHandler(r.DB.ExecContext(ctx, query, argv...))
}

// UpdateTx 通过事务更新图书，同时更新bookCategory表
func (r *bookRepo) UpdateTx(ctx context.Context, book *models.Book) error {
	ctx, cancel := timeoutCtx(ctx)
	defer cancel()

	err := execTx(context.Background(), r.DB, func(r Repository) error {
		// 更新图书
		err := r.BookRepo.Update(ctx, book)
		if err != nil {
			slog.ErrorContext(ctx, "[UpdateTx]->[r.BookRepo.Update] fail: ", slog.String("err", err.Error()))
			return err
		}

		// 删除图书和分类关系（book_categories）
		err = r.BookCategoryRepo.DeleteByBookID(ctx, book.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			slog.ErrorContext(ctx, "[UpdateTx]->[r.BookCategoryRepo.DeleteByBookID] fail: ", slog.String("err", err.Error()))
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
		return r.BookCategoryRepo.BatchInsert(ctx, list)
	})
	return err
}

func (r *bookRepo) List(ctx context.Context, filter Filters) ([]*models.Book, error) {
	/* NOTE: sqlx 查询嵌套结构体字段语法
	author_name as "author.author_name",
	publisher_name as "publisher.publisher_name"
	*/
	sql := r.DB.Rebind(`
	select 
		b.*, 
		a.id as "author.id",
		author_name as "author.author_name", 
		p.id as "publisher.id",
		publisher_name as "publisher.publisher_name" 
	from books b
	left join author a on a.id=b.author_id
	left join publisher p on p.id=b.publisher_id
	where b.deleted_at is null order by b.id desc limit ? offset ?`)
	ctx, cancel := timeoutCtx(ctx)
	defer cancel()

	list := make([]*models.Book, 0, filter.limit())

	// TODO: 返回 PageQueryVo 结构类型数据

	slog.InfoContext(
		ctx, spaceRex.ReplaceAllString(sql, " "),
		"limit", slog.IntValue(filter.limit()),
		"offset", slog.IntValue(filter.offset()))

	err := r.DB.SelectContext(ctx, &list, sql, filter.limit(), filter.offset())
	return list, err
}

// ListByCategory 根据分类查询图书 TODO: 分页
func (r *bookRepo) ListByCategory(ctx context.Context, categoryID uint64) (list []*models.Book, err error) {
	sql := `
		select 
			b.*, 
			a.id as "author.id",
			author_name as "author.author_name", 
			p.id as "publisher.id",
			publisher_name as "publisher.publisher_name"
		from books b 
		left join book_categories bc on b.id = bc.book_id
		left join category c on c.id = bc.category_id
		left join author a on a.id = b.author_id
		left join publisher p on p.id = b.publisher_id
		where bc.category_id = ? and b.deleted_at is null;`
	list = make([]*models.Book, 0, 10)

	ctx, cancel := timeoutCtx(ctx)
	defer cancel()

	err = r.DB.SelectContext(ctx, &list, sql, categoryID)
	return list, err
}

// ListWithCategory 查询图书列表和分类
func (r *bookRepo) ListWithCategory(ctx context.Context, filter Filters) (list []*models.Book, err error) {
	ctx, cancel := timeoutCtx(ctx)
	defer cancel()

	list, err = r.List(ctx, filter)
	if err != nil {
		return list, err
	}

	var bookIDs []uint64
	for _, x := range list {
		bookIDs = append(bookIDs, x.ID)
	}

	if len(bookIDs) <= 0 {
		return list, nil
	}

	sql, arg, err := sqlx.In(
		`select 
				c.id as "category.id",
				c.category_name as "category.category_name",
				c.icon as "category.icon",
				c.sort as "category.sort",
				c.created_at as "category.created_at",
				c.updated_at as "category.updated_at",
			 	bc.book_id AS "book_category.book_id",
    		bc.category_id AS "book_category.category_id"
			from category as c 
			left join book_categories bc on c.id = bc.category_id
			where bc.book_id in (?);
		`,
		bookIDs,
	)
	if err != nil {
		return list, err
	}

	var cc []models.SQLBookCategory

	sql = r.DB.Rebind(sql)

	slog.InfoContext(ctx, sql, "args", slog.AnyValue(arg))

	err = r.DB.Select(&cc, sql, arg...)
	if err != nil {
		return list, err
	}

	// 组合数据

	// 方式1:
	for i := range list {
		for _, x := range cc {
			if list[i].ID == x.BookCategory.BookID {
				list[i].Categories = append(list[i].Categories, &x.Category)
			}
		}
	}

	// 方式2:
	// bookMap := make(map[uint64]*models.Book)
	// // 图书list转map
	// for i := range list {
	// 	bookMap[list[i].ID] = list[i]
	// }
	// // 查找对应关系
	// for _, bc := range cc {
	// 	book := bookMap[bc.BookCategory.BookID]
	// 	if book != nil {
	// 		book.Categories = append(book.Categories, &bc.Category)
	// 	}
	// }

	return list, nil
}

func (r *bookRepo) Delete(ctx context.Context, id uint64) error {
	sql := `update books set deleted_at=now() where id=?`
	ctx, cancel := timeoutCtx(ctx)
	defer cancel()

	query := r.DB.Rebind(sql)
	slog.DebugContext(ctx, query, slog.Int64("id", int64(id)))

	return updateErrorHandler(r.DB.ExecContext(ctx, query, id))
}
