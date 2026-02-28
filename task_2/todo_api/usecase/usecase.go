package usecase

import "tasks/task_2/todo_api/domain"

type TodoUseCase struct {
	repo domain.TodoRepository
}

func NewTodoUseCase(r domain.TodoRepository) *TodoUseCase {
	return &TodoUseCase{repo: r}
}

func (uc *TodoUseCase) List() []domain.Todo              { return uc.repo.GetAll() }
func (uc *TodoUseCase) Create(t domain.Todo) domain.Todo { return uc.repo.Create(t) }
func (uc *TodoUseCase) Get(id int) (domain.Todo, bool)   { return uc.repo.GetByID(id) }
func (uc *TodoUseCase) Update(id int, t domain.Todo) (domain.Todo, bool) {
	return uc.repo.Update(id, t)
}
func (uc *TodoUseCase) Delete(id int) bool { return uc.repo.Delete(id) }
