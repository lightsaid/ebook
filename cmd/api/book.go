package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

func (app *Application) PostBookHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *Application) GetBookHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIntParam(r, "id")
	if err != nil {
		app.FAIL(w, r, err)
		return
	}

	fmt.Println(id)
}

func (app *Application) PutBookHandler(w http.ResponseWriter, r *http.Request)    {}
func (app *Application) DeleteBookHandler(w http.ResponseWriter, r *http.Request) {}

func (app *Application) ListBookHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info(r.URL.Path)
	ff := app.readPageQuery(r)
	list, err := app.Db.BookRepo.ListWithCategory(r.Context(), ff)
	if err != nil {
		slog.Info(err.Error())
		// TODO:
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}
