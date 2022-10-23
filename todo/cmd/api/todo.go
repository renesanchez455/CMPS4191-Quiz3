/*
	CMPS4191 - Quiz #3
	Rene Sanchez - 2018118383
*/
// Filename: cmd/api/todo.go

package main

import (
	"errors"
	"fmt"
	"net/http"

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

	// Create a Todo
	err = app.models.Todos.Insert(todo)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	// Create a Location header for the newly created resource/Todo
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/todo/%d", todo.ID))
	// Write the JSON response with 201 - Created status code with the body
	// being the Todo data and the header being the headers map
	err = app.writeJSON(w, http.StatusCreated, envelope{"todo": todo}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// showTodoHandler for the "GET /v1/todo/:id" endpoint
func (app *application) showTodoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Fetch the specific todo
	todo, err := app.models.Todos.Get(id)
	// Handle errors
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Write the data returned by Get()
	err = app.writeJSON(w, http.StatusOK, envelope{"todo": todo}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateTodoHandler(w http.ResponseWriter, r *http.Request) {
	// This method does a partial replacement
	// Get the id for the todo that needs updating
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Fetch the orginal record from the database
	todo, err := app.models.Todos.Get(id)
	// Handle errors
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Create an input struct to hold data read in from the client
	// We update input struct to use pointers because pointers have a
	// default value of nil
	// If a field remains nil then we know that the client did not update it
	var input struct {
		Name     *string `json:"name"`
		Details  *string `json:"Details"`
		Priority *string `json:"priority"`
		Status   *string `json:"status"`
	}

	// Initialize a new json.Decoder instance
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Check for updates
	if input.Name != nil {
		todo.Name = *input.Name
	}
	if input.Details != nil {
		todo.Details = *input.Details
	}
	if input.Priority != nil {
		todo.Priority = *input.Priority
	}
	if input.Status != nil {
		todo.Status = *input.Status
	}

	// Perform validation on the updated Todo. If validation fails, then
	// we send a 422 - Unprocessable Entity respose to the client
	// Initialize a new Validator instance
	v := validator.New()

	// Check the map to determine if there were any validation errors
	if data.ValidateTodo(v, todo); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Pass the updated Todo record to the Update() method
	err = app.models.Todos.Update(todo)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Write the data returned by Get()
	err = app.writeJSON(w, http.StatusOK, envelope{"todo": todo}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	// Get the id for the todo that needs updating
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Delete the Todo from the database. Send a 404 Not Found status code to the
	// client if there is no matching record
	err = app.models.Todos.Delete(id)
	// Handle errors
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Return 200 Status OK to the client with a success message
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "todo item successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// The listTodoHandler() allows the client to see a listing of todos
// based on a set of criteria
func (app *application) listTodoHandler(w http.ResponseWriter, r *http.Request) {
	// Create an input struct to hold our query parameters
	var input struct {
		Name  string
		Level string
		Mode  []string
		data.Filters
	}
	// Initialize a validator
	v := validator.New()
	// Get the URL values map
	qs := r.URL.Query()
	// Use the helper methods to extract the values
	input.Name = app.readString(qs, "name", "")
	input.Level = app.readString(qs, "level", "")
	input.Mode = app.readCSV(qs, "mode", []string{})
	// Get the page information
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	// Get the sort information
	input.Filters.Sort = app.readString(qs, "sort", "id")
	// Specific the allowed sort values
	input.Filters.SortList = []string{"id", "name", "level", "-id", "-name", "-level"}
	// Check for validation errors
	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Results dump
	fmt.Fprintf(w, "%+v\n", input)
}
