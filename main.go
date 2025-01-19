package main

import (
	"crud/todo-crap-app/routes"
	"fmt"
	"log"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

func main() {
	routes.HandleRoutes()
	fmt.Println("Starting server on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
