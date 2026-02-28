package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"tasks/task_2/todo_api/delivery"
	"tasks/task_2/todo_api/repository"
	"tasks/task_2/todo_api/usecase"

	_ "github.com/lib/pq"
)

func main() {
	// Agar bazangiz nomi 'todo_db' bo'lsa:
	connStr := "user=macbook dbname=todo sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	repo := repository.NewTodoRepo(db)
	uc := usecase.NewTodoUseCase(repo)
	handler := &delivery.TodoHandler{UC: uc}

	http.Handle("/todos", handler)
	http.Handle("/todos/", handler)

	fmt.Println("Clean Todo API running on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
