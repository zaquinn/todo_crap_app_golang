package crud

import (
	"crud/todo-crap-app/todo/database"
	"crud/todo-crap-app/todo/models"
	"encoding/json"
	"fmt"
	"net/http"
)

func AddNewTodo(w http.ResponseWriter, r *http.Request, db *database.Database) {
	var todo models.Todo
	json.NewDecoder(r.Body).Decode(&todo)
	res, err := models.NewTodo(&todo)

	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(err.Error())
	}
	if res != nil {
		fmt.Println("{}", string(db.Conn))
		json.NewEncoder(w).Encode(res)
	}

}
