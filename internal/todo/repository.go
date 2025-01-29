package todo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Repository interface {
	CreateTodo(ctx context.Context, newTodo *CreateTodoDTO) (*CreateTodoDTO, error)
	GetTodoByID(ctx context.Context, id string) (*Todo, error)
	GetTodos(ctx context.Context, limit, offset int) ([]Todo, error)
	DeleteTodoByID(ctx context.Context, id string) error
	PatchTodoByID(ctx context.Context, id string, updates map[string]string) error
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

func (r *repository) DeleteTodoByID(ctx context.Context, id string) error {
	query := `
        DELETE FROM todos
        WHERE id = $1
    `
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("todo not found")
	}

	return nil
}

func (r *repository) PatchTodoByID(ctx context.Context, id string, updates map[string]string) error {
	if len(updates) == 0 {
		return fmt.Errorf("no updates provided")
	}

	validFields := map[string]bool{"title": true, "message": true}
	for key := range updates {
		if !validFields[key] {
			return fmt.Errorf("invalid field: %s", key)
		}
	}

	var queryBuilder strings.Builder
	queryBuilder.WriteString("UPDATE todos SET updated_at = NOW()")

	var args []interface{}
	args = append(args, id)
	paramIndex := 2

	if title, ok := updates["title"]; ok {
		queryBuilder.WriteString(fmt.Sprintf(", title = $%d", paramIndex))
		args = append(args, title)
		paramIndex++
	}
	if message, ok := updates["message"]; ok {
		queryBuilder.WriteString(fmt.Sprintf(", message = $%d", paramIndex))
		args = append(args, message)
		paramIndex++
	}

	queryBuilder.WriteString(" WHERE id = $1 RETURNING id")
	query := queryBuilder.String()

	var updatedID string
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&updatedID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("todo not found")
		}
		return fmt.Errorf("failed to update todo: %w", err)
	}

	return nil
}
