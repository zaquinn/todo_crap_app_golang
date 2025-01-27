package todo

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository interface {
	CreateTodo(ctx context.Context, newTodo *CreateTodoDTO) (*CreateTodoDTO, error)
	GetTodoByID(ctx context.Context, id string) (*Todo, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateTodo(ctx context.Context, newTodo *CreateTodoDTO) (*CreateTodoDTO, error) {
	createdTodo := &CreateTodoDTO{
		Title:   newTodo.Title,
		Message: newTodo.Message,
	}

	var id int
	err := r.db.QueryRowContext(
		ctx,
		"INSERT INTO todos (title, message) VALUES ($1, $2) RETURNING id",
		newTodo.Title,
		newTodo.Message,
	).Scan(&id)

	if err != nil {
		return nil, fmt.Errorf("failed to create todo: %w", err)
	}

	createdTodo.Id = &id
	return createdTodo, nil
}

func (r *repository) GetTodoByID(ctx context.Context, id string) (*Todo, error) {
	var todo Todo
	query := `
		SELECT id, title, message, created_at, updated_at
		FROM todos
		WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&todo.Id,
		&todo.Title,
		&todo.Message,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("todo not found")
		}
		return nil, err
	}
	return &todo, nil
}
