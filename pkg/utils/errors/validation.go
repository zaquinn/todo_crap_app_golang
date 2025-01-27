package errors

import "fmt"

type ValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type ValidationErrorsList struct {
	Errors []ValidationError `json:"errors"`
}

func (v ValidationErrorsList) Error() []ValidationError {
	return v.Errors
}

func ValidateStrLen(str string, minSize int) error {
	if len(str) < minSize {
		return fmt.Errorf("min length is %d", minSize)
	}
	return nil
}
