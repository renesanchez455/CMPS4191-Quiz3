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
	// Our target decode destination
	var input struct {
		Name     string `json:"name"`
		Details  string `json:"Details"`
		Priority string `json:"priority"`
		Status   string `json:"status"`
	}
	// Initialize a new json.Decoder instance
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	// Display the request
	fmt.Fprintf(w, "%+v\n", input)
}

// showSchoolHandler for the "GET /v1/todo/:id" endpoint
func (app *application) showTodoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
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
	err = app.writeJSON(w, http.StatusOK, envelope{"school": todo}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
