package todo

import (
	"crud/todo-crap-app/pkg/middleware"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, repository Repository) {
	mux.HandleFunc("POST /todo", middleware.AuthorizationMiddleware([]string{"admin"}, createTodoHandler(repository)))
	mux.HandleFunc("GET /todo/{id}", middleware.AuthorizationMiddleware([]string{"admin"}, retrieveTodoHandler(repository)))
	mux.HandleFunc("GET /todo", middleware.AuthorizationMiddleware([]string{"admin"}, listTodosHandler(repository)))
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

func listTodosHandler(repository Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ListTodosService(w, r, repository)
	}
}
