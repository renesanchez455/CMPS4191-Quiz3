/*
	CMPS4191 - Quiz #3
	Rene Sanchez - 2018118383
*/
// Filename: internal/data/todo.go

package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"todo.renesanchez455.net/internal/validator"
)

type Todo struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Name      string    `json:"name"`
	Details   string    `json:"details"`
	Priority  string    `json:"priority"`
	Status    string    `json:"status"`
	Version   int32     `json:"version"`
}

func ValidateTodo(v *validator.Validator, todo *Todo) {
	// Use the Check() method to execute our validation checks
	v.Check(todo.Name != "", "name", "must be provided")
	v.Check(len(todo.Name) <= 200, "name", "must not be more than 200 bytes long")

	v.Check(todo.Details != "", "details", "must be provided")
	v.Check(len(todo.Details) <= 300, "details", "must not be more than 300 bytes long")

	v.Check(todo.Priority != "", "priority", "must be provided")
	v.Check(len(todo.Priority) <= 100, "priority", "must not be more than 100 bytes long")

	v.Check(todo.Status != "", "status", "must be provided")
	v.Check(len(todo.Status) <= 100, "status", "must not be more than 100 bytes long")
}

// Define a TodoModel which wraps a sql.DB connection pool
type TodoModel struct {
	DB *sql.DB
}

// Insert() allows us  to create a new Todo
func (m TodoModel) Insert(todo *Todo) error {
	query := `
		INSERT INTO todo (name, details, priority, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version
	`
	// Collect the data fields into a slice
	args := []interface{}{
		todo.Name, todo.Details,
		todo.Priority, todo.Status,
	}
	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Cleanup to prevent memory leaks
	defer cancel()
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&todo.ID, &todo.CreatedAt, &todo.Version)
}

// Get() allows us to retrieve a specific Todo
func (m TodoModel) Get(id int64) (*Todo, error) {
	// Ensure that there is a valid id
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	// Create the query
	query := `
		SELECT id, created_at, name, details, priority, status, version
		FROM todo
		WHERE id = $1
	`
	// Declare a Todo variable to hold the returned data
	var todo Todo

	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Cleanup to prevent memory leaks
	defer cancel()

	// Execute the query using QueryRow()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&todo.ID,
		&todo.CreatedAt,
		&todo.Name,
		&todo.Details,
		&todo.Priority,
		&todo.Status,
		&todo.Version,
	)
	// Handle any errors
	if err != nil {
		// Check the type of error
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	// Success
	return &todo, nil
}

// Update() allows us to edit/alter a specific Todo
// Optimistic locking (version number)
func (m TodoModel) Update(todo *Todo) error {
	// Create a query
	query := `
		UPDATE todo
		SET name = $1, details = $2, priority = $3,
		    status = $4, version = version + 1
		WHERE id = $5
		AND version = $6
		RETURNING version
	`
	args := []interface{}{
		todo.Name,
		todo.Details,
		todo.Priority,
		todo.Status,
		todo.ID,
		todo.Version,
	}

	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Cleanup to prevent memory leaks
	defer cancel()
	// Check for edit conflicts
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&todo.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

// Delete() removes a specific Todo
func (m TodoModel) Delete(id int64) error {
	// Ensure that there is a valid id
	if id < 1 {
		return ErrRecordNotFound
	}
	// Create the delete query
	query := `
		DELETE FROM todo
		WHERE id = $1
	`

	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Cleanup to prevent memory leaks
	defer cancel()

	// Execute the query
	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	// Check how many rows were affected by the delete operation. We
	// call the RowsAffected() method on the result variable
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// Check if no rows were affected
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

// The GetAll() method retuns a list of all the todos sorted by id
func (m TodoModel) GetAll(name string, priority string, filters Filters) ([]*Todo, Metadata, error) {

	// Construct the query

	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, created_at, name, details,
		    priority, status, version
		FROM todo
		WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (to_tsvector('simple', priority) @@ plainto_tsquery('simple', $2) OR $2 = '')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortOrder())

	// Create a 3-second-timout context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Execute the query
	args := []interface{}{name, priority, filters.limit(), filters.offset()}
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	// Close the resultset
	defer rows.Close()
	totalRecords := 0
	// Initialize an empty slice to hold the Todo data
	todos := []*Todo{}
	// Iterate over the rows in the resultset
	for rows.Next() {
		var todo Todo
		// Scan the values from the row into todo
		err := rows.Scan(
			&totalRecords,
			&todo.ID,
			&todo.CreatedAt,
			&todo.Name,
			&todo.Details,
			&todo.Priority,
			&todo.Status,
			&todo.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		// Add the Todo to our slice
		todos = append(todos, &todo)
	}
	// Check for errors after looping through the resultset
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	// Return the slice of Todos
	return todos, metadata, nil
}
