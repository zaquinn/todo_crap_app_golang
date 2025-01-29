package todo

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type ResponseError struct {
	Message string `json:"message"`
}

func CreateTodoService(w http.ResponseWriter, r *http.Request, repository Repository) {
	ctx := r.Context()
	var newTodo CreateTodoDTO

	if err := json.NewDecoder(r.Body).Decode(&newTodo); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&ResponseError{Message: "error parsing the json"})
		return
	}

	validationErr := ValidateNewTodo(&newTodo)
	if validationErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validationErr.Error())
		return
	}

	responseTodo, createErr := repository.CreateTodo(ctx, &newTodo)
	if createErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(createErr.Error())
		return
	}

	json.NewEncoder(w).Encode(responseTodo)
}

func RetrieveTodoService(w http.ResponseWriter, r *http.Request, repository Repository) {
	ctx := r.Context()
	id := r.PathValue("id")
	res, err := repository.GetTodoByID(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&ResponseError{Message: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(res)
}

func UpdateTodoService(w http.ResponseWriter, r *http.Request, repository Repository) {
	const (
		ErrInvalidRequestBody = "invalid request body"
		ErrNoFieldsToUpdate   = "no fields to update"
		ErrTodoNotFound       = "todo not found"
		ErrInternalServer     = "internal server error"
	)

	ctx := r.Context()
	id := r.PathValue("id")

	if r.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&ResponseError{Message: ErrInvalidRequestBody})
		return
	}

	var updates map[string]string
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&ResponseError{Message: ErrInvalidRequestBody})
		return
	}

	if len(updates) == 0 || (updates["title"] == "" && updates["message"] == "") {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&ResponseError{Message: ErrNoFieldsToUpdate})
		return
	}

	err = repository.PatchTodoByID(ctx, id, updates)
	if err != nil {
		switch err.Error() {
		case ErrTodoNotFound:
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(&ResponseError{Message: ErrTodoNotFound})
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&ResponseError{Message: ErrInternalServer})
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": http.StatusText(http.StatusOK)})
}

func DeleteTodoService(w http.ResponseWriter, r *http.Request, repository Repository) {
	ctx := r.Context()
	id := r.PathValue("id")
	err := repository.DeleteTodoByID(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&ResponseError{Message: err.Error()})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func ListTodosService(w http.ResponseWriter, r *http.Request, repository Repository) {
	ctx := r.Context()

	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}

	offset := 0
	if offsetStr != "" {
		offset, _ = strconv.Atoi(offsetStr)
	}

	res, err := repository.GetTodos(ctx, limit, offset)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&ResponseError{Message: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
