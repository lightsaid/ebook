package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/lightsaid/ebook/internal/dbrepo"
	"github.com/lightsaid/ebook/pkg/logger"
	"github.com/lightsaid/gotk"
)

type Application struct {
	Db dbrepo.Repository
}

type envelope map[string]any

// @title EBook API
// @version 1.0
// @description EBook API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /api
func main() {
	app := Application{}

	instance := logger.NewLogger(os.Stdout, "DEBUG", gotk.TextType)
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
