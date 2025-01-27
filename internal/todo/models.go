package todo

import (
	"crud/todo-crap-app/pkg/utils/errors"
	"time"
)

type Todo struct {
	Id        int        `json:"id"`
	Title     string     `json:"title"`
	Message   string     `json:"message"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func ValidateNewTodo(todo *CreateTodoDTO) *errors.ValidationErrorsList {
	var validationErrors errors.ValidationErrorsList

	if err := errors.ValidateStrLen(todo.Title, 3); err != nil {
		validationErrors.Errors = append(validationErrors.Errors, errors.ValidationError{
			Field: "title",
			Error: err.Error(),
		})
	}

	if err := errors.ValidateStrLen(todo.Message, 3); err != nil {
		validationErrors.Errors = append(validationErrors.Errors, errors.ValidationError{
			Field: "message",
			Error: err.Error(),
		})
	}

	if len(validationErrors.Errors) > 0 {
		return &validationErrors
	}

	return nil
}
