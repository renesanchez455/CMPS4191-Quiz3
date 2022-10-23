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
	"todo.renesanchez455.net/internal/validator"
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

	// Copy the values from the input struct to a new Todo struct
	todo := &data.Todo{
		Name:     input.Name,
		Details:  input.Details,
		Priority: input.Priority,
		Status:   input.Status,
	}

	// Initialize a new Validator instance
	v := validator.New()

	// Check the map to determine if there were any validation errors
	if data.ValidateTodo(v, todo); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
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

	// Create a new instance of the Todo struct containing the ID we extracted
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
	err = app.writeJSON(w, http.StatusOK, envelope{"todo": todo}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
