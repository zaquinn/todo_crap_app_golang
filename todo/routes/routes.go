package routes

import (
	"crud/todo-crap-app/todo/crud"
	"crud/todo-crap-app/todo/database"
	"crud/todo-crap-app/todo/middleware"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

func HandleRoutes(db *database.Database) {
	http.HandleFunc("POST /todo", middleware.AuthorizationMiddleware([]string{"admin", "user"}, func(w http.ResponseWriter, r *http.Request) {
		crud.AddNewTodo(w, r, db)
	}))
}
