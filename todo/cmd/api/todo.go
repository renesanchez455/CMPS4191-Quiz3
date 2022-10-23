/*
	CMPS4191 - Quiz #3
	Rene Sanchez - 2018118383
*/
// Filename: cmd/api/todo.go

package main

import (
	"fmt"
	"net/http"
	"time"

	"todo.renesanchez455.net/internal/data"
)

// createTodoHandler for the "POST /v1/todo" endpoint
func (app *application) createTodoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new todo..")
}

// showSchoolHandler for the "GET /v1/todo/:id" endpoint
func (app *application) showTodoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Create a new instance of the School struct containing the ID we extracted
	// from our URL and some sample data
	todo := data.Todo{
		ID:        id,
		CreatedAt: time.Now(),
		Name:      "Laundry",
		Details:   "Wash white shirts",
		Priority:  "High",
		Status:    "Pending",
		Version:   1,
	}
	err = app.writeJSON(w, http.StatusOK, todo, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}

}
