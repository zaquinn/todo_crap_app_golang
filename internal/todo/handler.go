package todo

import (
	"crud/todo-crap-app/pkg/middleware"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, repository Repository) {
	mux.HandleFunc("POST /todo", middleware.AuthorizationMiddleware([]string{"admin"}, createTodoHandler(repository)))
	mux.HandleFunc("GET /todo/{id}", middleware.AuthorizationMiddleware([]string{"admin"}, retrieveTodoHandler(repository)))
}

func createTodoHandler(repository Repository) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		CreateTodoService(w, r, repository)
	}
}

func retrieveTodoHandler(repository Repository) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		RetrieveTodoService(w, r, repository)
	}
}
