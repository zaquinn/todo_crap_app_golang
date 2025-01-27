package todo

import (
	"encoding/json"
	"net/http"
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
