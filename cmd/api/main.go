package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/lightsaid/ebook/internal/dbrepo"
	"github.com/lightsaid/ebook/pkg/logger"
)

type Application struct {
	Db dbrepo.Repository
}

type envelope map[string]any

func main() {
	app := Application{}

	instance := logger.NewLogger(os.Stdout, "DEBUG", logger.TextStyle)
	slog.SetDefault(instance)

	conn, err := dbrepo.Open()
	if err != nil {
		panic(err)
	}

	app.Db = dbrepo.NewRepository(conn)

	if err := app.serve(instance); err != nil {
		log.Fatalln(err)
	}

}
