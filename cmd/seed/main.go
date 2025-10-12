package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/lightsaid/ebook/internal/dbrepo"
)

var repo dbrepo.Repository

func main() {
	db, err := dbrepo.Open()
	if err != nil {
		panic(err)
	}

	repo = dbrepo.NewRepository(db)
	id := createAuthor()
	getAuthor(id)

	list, err := repo.AuthorRepo.List()
	if err != nil {
		log.Println("repo.AuthorRepo.List error: ", err)
	}
	str, _ := json.MarshalIndent(list, "", "\t")
	fmt.Println(string(str))
}

func createAuthor() uint64 {
	src := rand.NewSource(time.Now().UnixMicro())
	randNumber := rand.New(src).Intn(2000)
	id, err := repo.AuthorRepo.Create(fmt.Sprintf("法外狂徒张三-%d", randNumber))
	if err != nil {
		log.Println("createAuthor: ", err)
	}
	fmt.Println("createAuthor id: ", id)
	return id
}

func getAuthor(id uint64) {
	author, err := repo.AuthorRepo.Get(id)
	if err != nil {
		log.Println("getAuthor error: ", err)
	}
	fmt.Println("getAuthor succ: ", author)
}
