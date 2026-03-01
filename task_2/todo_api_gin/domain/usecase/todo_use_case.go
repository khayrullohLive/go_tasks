package usecase

import (
	"tasks/task_2/todo_api_gin/data/request"
	"tasks/task_2/todo_api_gin/domain/model"
	"tasks/task_2/todo_api_gin/domain/repository"
	"time"

	"github.com/google/uuid"
)

// TodoUseCase - todo biznes mantiq interfeysi
type TodoUseCase interface {
	GetAll() ([]*model.Todo, error)
	GetByID(id string) (*model.Todo, error)
	Create(req request.CreateTodoRequest) (*model.Todo, error)
	Update(id string, req request.UpdateTodoRequest) (*model.Todo, error)
	Delete(id string) error
}

type todoUseCase struct {
	repo repository.TodoRepository
}

// NewTodoUseCase - yangi service yaratadi
func NewTodoUseCase(repo repository.TodoRepository) TodoUseCase {
	return &todoUseCase{repo: repo}
}

func (s *todoUseCase) GetAll() ([]*model.Todo, error) {
	return s.repo.GetAll()
}

func (s *todoUseCase) GetByID(id string) (*model.Todo, error) {
	return s.repo.GetByID(id)
}

func (s *todoUseCase) Create(req request.CreateTodoRequest) (*model.Todo, error) {
	now := time.Now()
	todo := &model.Todo{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Description: req.Description,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(todo); err != nil {
		return nil, err
	}
	return todo, nil
}

func (s *todoUseCase) Update(id string, req request.UpdateTodoRequest) (*model.Todo, error) {
	todo, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		todo.Title = *req.Title
	}
	if req.Description != nil {
		todo.Description = *req.Description
	}
	if req.Completed != nil {
		todo.Completed = *req.Completed
	}
	todo.UpdatedAt = time.Now()

	if err := s.repo.Update(todo); err != nil {
		return nil, err
	}
	return todo, nil
}

func (s *todoUseCase) Delete(id string) error {
	return s.repo.Delete(id)
}
