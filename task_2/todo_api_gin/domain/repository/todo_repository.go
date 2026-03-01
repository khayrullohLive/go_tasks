package repository

import (
	"errors"
	"sync"
	"tasks/task_2/todo_api_gin/domain/model"
)

// ErrNotFound - todo topilmadi xatosi
var ErrNotFound = errors.New("todo not found")

// TodoRepository - todo ma'lumotlar bazasi interfeysi
type TodoRepository interface {
	GetAll() ([]*model.Todo, error)
	GetByID(id string) (*model.Todo, error)
	Create(todo *model.Todo) error
	Update(todo *model.Todo) error
	Delete(id string) error
}

// inMemoryTodoRepository - xotirada saqlanadigan implementatsiya
type inMemoryTodoRepository struct {
	mu    sync.RWMutex
	todos map[string]*model.Todo
}

// NewInMemoryTodoRepository - yangi repository yaratadi
func NewInMemoryTodoRepository() TodoRepository {
	return &inMemoryTodoRepository{
		todos: make(map[string]*model.Todo),
	}
}

func (r *inMemoryTodoRepository) GetAll() ([]*model.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	todos := make([]*model.Todo, 0, len(r.todos))
	for _, todo := range r.todos {
		todos = append(todos, todo)
	}
	return todos, nil
}

func (r *inMemoryTodoRepository) GetByID(id string) (*model.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	todo, exists := r.todos[id]
	if !exists {
		return nil, ErrNotFound
	}
	return todo, nil
}

func (r *inMemoryTodoRepository) Create(todo *model.Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.todos[todo.ID] = todo
	return nil
}

func (r *inMemoryTodoRepository) Update(todo *model.Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.todos[todo.ID]; !exists {
		return ErrNotFound
	}
	r.todos[todo.ID] = todo
	return nil
}

func (r *inMemoryTodoRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.todos[id]; !exists {
		return ErrNotFound
	}
	delete(r.todos, id)
	return nil
}
