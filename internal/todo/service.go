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
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(&ResponseError{Message: "error parsing the json"})
		return
	}

	validationErr := ValidateNewTodo(&newTodo)
	if validationErr != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(validationErr.Error())
		return
	}

	responseTodo, createErr := repository.CreateTodo(ctx, &newTodo)
	if createErr != nil {
		w.WriteHeader(400)
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
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(&ResponseError{Message: err.Error()})
		return
	}
	json.NewEncoder(w).Encode(res)
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
