package main

import (
	"context"
	"crud/todo-crap-app/internal/todo"
	"crud/todo-crap-app/pkg/middleware"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://root:root@localhost:5432/todo_crap_app?sslmode=disable")
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer db.Close()

	if err := runMigrations(db); err != nil {
		log.Fatalf("Could not run migrations: %v", err)
	}

	mux := http.NewServeMux()
	todoRepository := todo.NewRepository(db)
	todo.RegisterRoutes(mux, todoRepository)

	server := &http.Server{
		Addr:    ":8080",
		Handler: middleware.LoggingMiddleware(middleware.CorsMiddleware(mux)),
	}

	go func() {
		log.Printf("Starting server on http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
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
