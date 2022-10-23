// Filename: internal/data/todo.go

package data

import (
	"database/sql"
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
	return m.DB.QueryRow(query, args...).Scan(&todo.ID, &todo.CreatedAt, &todo.Version)
}

// Get() allows us to retrieve a specific Todo
func (m TodoModel) Get(id int64) (*Todo, error) {
	return nil, nil
}

// Update() allows us to edit/alter a specific Todo
func (m TodoModel) Update(todo *Todo) error {
	return nil
}

// Delete() removes a specific Todo
func (m TodoModel) Delete(id int64) error {
	return nil
}
