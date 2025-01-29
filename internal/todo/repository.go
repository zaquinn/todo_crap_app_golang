package todo

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository interface {
	CreateTodo(ctx context.Context, newTodo *CreateTodoDTO) (*CreateTodoDTO, error)
	GetTodoByID(ctx context.Context, id string) (*Todo, error)
	GetTodos(ctx context.Context, limit, offset int) ([]Todo, error)
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

func (r *repository) GetTodos(ctx context.Context, limit, offset int) ([]Todo, error) {
	var todos []Todo

	query := `
		SELECT id, title, message, created_at, updated_at
		FROM todos
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query todos: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var todo Todo
		err := rows.Scan(
			&todo.Id,
			&todo.Title,
			&todo.Message,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan todo: %w", err)
		}
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows: %w", err)
	}

	return todos, nil
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
