package domain

type Todo struct {
	ID        int    `json:"id"`
	Task      string `json:"task"`
	Completed bool   `json:"completed"`
}

// Repository interfeysi - Ma'lumot bilan kim ishlasa shu metodlarni ta'minlashi shart
type TodoRepository interface {
	GetAll() []Todo
	GetByID(id int) (Todo, bool)
	Create(t Todo) Todo
	Update(id int, t Todo) (Todo, bool)
	Delete(id int) bool
}
