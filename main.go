package main

import (
	"crud/todo-crap-app/internal/todo"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://root:root@localhost:5432/todo_crap_app?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := runMigrations(db); err != nil {
		log.Fatalf("Could not run migrations: %v", err)
	}

	mux := http.NewServeMux()

	todoRepository := todo.NewRepository(db)

	todo.RegisterRoutes(mux, todoRepository)

	fmt.Println("Starting server on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver,
	)
	if err != nil {
		return fmt.Errorf("could not create migration instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run migrations: %w", err)
	}

	fmt.Println("Migrations ran successfully!")
	return nil
}
