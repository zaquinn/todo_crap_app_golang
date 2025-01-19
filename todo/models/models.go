package models

import "fmt"

type Todo struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type TodoValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type TodoValidationErrorsList struct {
	Errors []TodoValidationError `json:"errors"`
}

func (v TodoValidationErrorsList) Error() []TodoValidationError {
	return v.Errors
}

func NewTodo(todo *Todo) (*Todo, *TodoValidationErrorsList) {
	var validationErrors TodoValidationErrorsList

	if err := validateTodoTitleLen(todo.Title); err != nil {
		validationErrors.Errors = append(validationErrors.Errors, TodoValidationError{
			Field: "title",
			Error: err.Error(),
		})
	}

	if err := validateTodoMessageLen(todo.Message); err != nil {
		validationErrors.Errors = append(validationErrors.Errors, TodoValidationError{
			Field: "message",
			Error: err.Error(),
		})
	}

	if len(validationErrors.Errors) > 0 {
		return nil, &validationErrors
	}

	return todo, nil
}

func validateTodoTitleLen(title string) error {
	if len(title) < 3 {
		return fmt.Errorf("min length is 3")
	}
	return nil
}

func validateTodoMessageLen(message string) error {
	if len(message) < 3 {
		return fmt.Errorf("min length is 3")
	}
	return nil
}
