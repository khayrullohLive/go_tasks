package delivery

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"tasks/task_2/todo_api/domain"
	"tasks/task_2/todo_api/usecase"
)

type TodoHandler struct {
	UC *usecase.TodoUseCase
}

func (h *TodoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Juda oddiy router logikasi
	path := strings.TrimPrefix(r.URL.Path, "/todos")
	id, _ := strconv.Atoi(strings.TrimPrefix(path, "/"))

	switch {
	case r.Method == http.MethodGet && path == "":
		err := json.NewEncoder(w).Encode(h.UC.List())
		if err != nil {
			return
		}
	case r.Method == http.MethodPost:
		var t domain.Todo
		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			return
		}
		err = json.NewEncoder(w).Encode(h.UC.Create(t))
		if err != nil {
			return
		}
	case r.Method == http.MethodPut && id != 0:
		var t domain.Todo
		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			return
		}
		updated, ok := h.UC.Update(id, t)
		if !ok {
			w.WriteHeader(404)
			return
		}
		err = json.NewEncoder(w).Encode(updated)
		if err != nil {
			return
		}
	case r.Method == http.MethodDelete && id != 0:
		if ok := h.UC.Delete(id); !ok {
			w.WriteHeader(404)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
