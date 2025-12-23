package main

import "net/http"

func (app *Application) PostOrderHandler(w http.ResponseWriter, r *http.Request)   {}
func (app *Application) PayOrderHandler(w http.ResponseWriter, r *http.Request)    {}
func (app *Application) GetOrderHandler(w http.ResponseWriter, r *http.Request)    {}
func (app *Application) PutOrderHandler(w http.ResponseWriter, r *http.Request)    {}
func (app *Application) DeleteOrderHandler(w http.ResponseWriter, r *http.Request) {}
func (app *Application) ListOrderHandler(w http.ResponseWriter, r *http.Request)   {}
