package main

// import (
// 	"errors"
// 	"log"

// 	"github.com/lightsaid/ebook/internal/app"
// )

// var ErrrNotAllowExt = errors.New("不支持文件类型")

// func main() {

// 	if err := app.Serve(); err != nil {
// 		log.Fatal(err)
// 	}
// }

// NOTE: 例子1
// // 第三方的接口定义
// type Worker interface {
// 	Context() context.Context
// 	DoSomeThing()
// 	// ...
// }

// type wrapWorker interface {
// 	Worker
// 	SetContext(context.Context)
// }
// type wrap struct {
// 	Worker
// 	ctx context.Context
// }

// func newWrap(w Worker) wrapWorker {
// 	return &wrap{
// 		w,
// 		w.Context(),
// 	}
// }
// func (wp *wrap) SetContext(ctx context.Context) {
// 	wp.ctx = ctx
// }

// func (wp wrap) Context() context.Context {
// 	return wp.ctx
// }

// type contextKey string

// func work(w Worker) {
// 	wp := newWrap(w)
// 	ctx := context.WithValue(w.Context(), contextKey("greet"), "Hello")
// 	wp.SetContext(ctx)

// 	next(wp)
// }

// func next(w Worker) {
// 	v := w.Context().Value(contextKey("greet"))
// 	fmt.Printf("-> %v \n", v)
// 	w.DoSomeThing()
// }

// type person string

// func (person) Context() context.Context {
// 	return context.Background()
// }
// func (p person) DoSomeThing() {
// 	fmt.Println(string(p), "吃饭睡觉打代码～")
// }

// func main() {
// 	var p person = "张三"
// 	work(p)
// }

// NOTE： 例子2

// NOTE: storage.go
// 存储文件接口
// type Storage interface {
// 	Save(data []byte) (string, error)
// }

// // 开发时候使用本地存储
// type localStorage struct {
// 	config *struct{ base string }
// }

// func NewLocalStorage(cfg struct{ base string }) *localStorage {
// 	return &localStorage{config: &cfg}
// }

// func (store localStorage) Save(data []byte) (string, error) {
// 	return "http://127.0.0.1/" + store.config.base + "/a.png", nil
// }

// // 生产使用oss存储
// type ossStorage struct{}

// func NewOSSStorage() *ossStorage {
// 	return &ossStorage{}
// }

// func (store ossStorage) Save(data []byte) (string, error) {
// 	return "https://abc.com/oss/a.png", nil
// }

// // NOTE: upload.go
// // 保存文件
// func saveFile(store Storage) {
// 	var data []byte // 伪代码
// 	url, _ := store.Save(data)
// 	fmt.Println(url)
// }

// func main() {
// 	var store Storage
// 	if os.Getenv("ENV") == "prod" {
// 		store = NewOSSStorage()
// 	} else {
// 		store = NewLocalStorage(struct{ base string }{base: "/static"})
// 	}
// 	saveFile(store)
// }

// NOTE: 例子3

// type Repository interface {
// 	Insert()
// 	Update()
// }

// type repository struct {
// 	DB *sql.DB
// }

// // NOTE: 接口检查
// var _ Repository = (*repository)(nil)

// func NewRepository(db *sql.DB) *repository {
// 	return &repository{DB: db}
// }

// func (repo *repository) Insert() {
// 	// do some thing
// }

// func (repo *repository) Update() {
// 	// do some thing
// }

// type Service struct {
// 	Repo Repository
// }

// func NewService(r Repository) *Service {
// 	return &Service{
// 		Repo: r,
// 	}
// }

// func main() {
// 	// 假设是真实的链接 &sql.DB{}
// 	r := NewRepository(&sql.DB{})
// 	srv := NewService(r)
// 	srv.Repo.Insert()
// 	srv.Repo.Update()
// }
