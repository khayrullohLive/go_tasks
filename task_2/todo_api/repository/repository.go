package repository

import (
	"database/sql"
	"errors"
	"log/slog"
	"sync"
	"tasks/task_2/todo_api/domain"
)

type postgresRepo struct {
	mu    sync.Mutex
	todos []domain.Todo
	db    *sql.DB
}

func NewTodoRepo(db *sql.DB) domain.TodoRepository {
	return &postgresRepo{
		todos: []domain.Todo{{ID: 1, Task: "Go o'rganish", Completed: false}},
		db:    db,
	}
}

func (r *postgresRepo) GetAll() []domain.Todo {
	rows, _ := r.db.Query("SELECT id, task, completed FROM todos")
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			slog.Error("Error closing rows:")
			return
		}
	}(rows)

	var todos []domain.Todo
	for rows.Next() {
		var t domain.Todo
		err := rows.Scan(&t.ID, &t.Task, &t.Completed)
		if err != nil {
			return nil
		}
		todos = append(todos, t)
	}
	return todos
}

func (r *postgresRepo) Create(t domain.Todo) domain.Todo {
	err := r.db.QueryRow(
		"INSERT INTO todos (task, completed) VALUES ($1, $2) RETURNING id",
		t.Task, t.Completed,
	).Scan(&t.ID)

	if err != nil {
		return domain.Todo{}
	}
	return t
}

// GetByID - ID bo'yicha bitta
func (r *postgresRepo) GetByID(id int) (domain.Todo, bool) {
	var t domain.Todo
	err := r.db.QueryRow("SELECT id, task, completed FROM todos WHERE id = $1", id).
		Scan(&t.ID, &t.Task, &t.Completed)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Todo{}, false
		}
		slog.Error("GetByID xatosi:", "error", err)
		return domain.Todo{}, false
	}
	return t, true
}

// Update - Mavjud todoni yangilash
func (r *postgresRepo) Update(id int, t domain.Todo) (domain.Todo, bool) {
	// UPDATE so'rovi va o'zgargan qatorni qaytarib olish
	res, err := r.db.Exec(
		"UPDATE todos SET task = $1, completed = $2 WHERE id = $3",
		t.Task, t.Completed, id,
	)
	if err != nil {
		slog.Error("Update xatosi:", "error", err)
		return domain.Todo{}, false
	}

	// Chindan ham biror qator o'zgardimi yoki bu ID yo'qmi?
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return domain.Todo{}, false
	}

	t.ID = id
	return t, true
}

// Delete - ID bo'yicha o'chirish
func (r *postgresRepo) Delete(id int) bool {
	res, err := r.db.Exec("DELETE FROM todos WHERE id = $1", id)
	if err != nil {
		slog.Error("Delete xatosi:", "error", err)
		return false
	}

	rowsAffected, _ := res.RowsAffected()
	return rowsAffected > 0
}
