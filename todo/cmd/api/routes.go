/*
	CMPS4191 - Quiz #3
	Rene Sanchez - 2018118383
*/
// Filename: cmd/api/routes.go

package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	// Create a new httprouter instance
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/todo", app.listTodoHandler)
	router.HandlerFunc(http.MethodPost, "/v1/todo", app.createTodoHandler)
	router.HandlerFunc(http.MethodGet, "/v1/todo/:id", app.showTodoHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/todo/:id", app.updateTodoHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/todo/:id", app.deleteTodoHandler)
	return router
}
