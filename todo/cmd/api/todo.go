/*
	CMPS4191 - Quiz #3
	Rene Sanchez - 2018118383
*/
// Filename: cmd/api/todo.go

package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// createTodoHandler for the "POST /v1/todo" endpoint
func (app *application) createTodoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new todo..")
}

// showSchoolHandler for the "GET /v1/todo/:id" endpoint
func (app *application) showTodoHandler(w http.ResponseWriter, r *http.Request) {
	// Use the "ParamsFromContext()" function to get the request context as a slice
	params := httprouter.ParamsFromContext(r.Context())
	// Get the value of the "id" parameter
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	// Display the todo id
	fmt.Fprintf(w, "show the details for  %d\n", id)
}
