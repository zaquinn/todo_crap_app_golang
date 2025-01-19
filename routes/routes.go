package routes

import (
	db_todo "crud/todo-crap-app/todo/database"
	todo "crud/todo-crap-app/todo/routes"
)

func HandleRoutes() {
	todo.HandleRoutes(db_todo.PostgreConn())
}
