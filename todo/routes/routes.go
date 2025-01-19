package routes

import (
	"crud/todo-crap-app/todo/crud"
	"crud/todo-crap-app/todo/database"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

func HandleRoutes(db *database.Database) {
	http.HandleFunc("POST /todo", func(w http.ResponseWriter, r *http.Request) {
		crud.AddNewTodo(w, r, db)
	})
}
