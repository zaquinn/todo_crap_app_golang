package todo

import (
	"crud/todo-crap-app/pkg/middleware"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, repository Repository) {
	mux.HandleFunc("POST /todo", middleware.AuthorizationMiddleware([]string{"admin"}, createTodoHandler(repository)))
	mux.HandleFunc("GET /todo/{id}", middleware.AuthorizationMiddleware([]string{"admin"}, retrieveTodoHandler(repository)))
	mux.HandleFunc("DELETE /todo/{id}", middleware.AuthorizationMiddleware([]string{"admin"}, deleteTodoHandler(repository)))
	mux.HandleFunc("PATCH /todo/{id}", middleware.AuthorizationMiddleware([]string{"admin"}, updateTodoHandler(repository)))
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

func deleteTodoHandler(repository Repository) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		DeleteTodoService(w, r, repository)
	}
}

func updateTodoHandler(repository Repository) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		UpdateTodoService(w, r, repository)
	}
}

func listTodosHandler(repository Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ListTodosService(w, r, repository)
	}
}
