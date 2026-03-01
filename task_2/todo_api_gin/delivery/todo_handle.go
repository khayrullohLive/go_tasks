package handler

import (
	"errors"
	"tasks/task_2/todo_api_gin/data/request"
	"tasks/task_2/todo_api_gin/data/response"
	"tasks/task_2/todo_api_gin/domain/repository"
	"tasks/task_2/todo_api_gin/domain/usecase"

	"github.com/gin-gonic/gin"
)

// TodoHandler - todo HTTP handler
type TodoHandler struct {
	usecase usecase.TodoUseCase
}

// NewTodoHandler - yangi handler yaratadi
func NewTodoHandler(svc usecase.TodoUseCase) *TodoHandler {
	return &TodoHandler{usecase: svc}
}

// RegisterRoutes - route'larni ro'yxatdan o'tkazadi
func (h *TodoHandler) RegisterRoutes(router *gin.RouterGroup) {
	todos := router.Group("/todos")
	{
		todos.GET("", h.GetAll)
		todos.GET("/:id", h.GetByID)
		todos.POST("", h.Create)
		todos.PATCH("/:id", h.Update)
		todos.DELETE("/:id", h.Delete)
	}
}

// GetAll godoc
// @Summary      Barcha todo'larni olish
// @Description  Barcha todo'larni ro'yxatini qaytaradi
// @Tags         todos
// @Produce      json
// @Success      200  {object}  response.Response
// @Router       /todos [get]
func (h *TodoHandler) GetAll(c *gin.Context) {
	todos, err := h.usecase.GetAll()
	if err != nil {
		response.InternalError(c, "Todo'larni olishda xato yuz berdi")
		return
	}
	response.OK(c, "Todo'lar muvaffaqiyatli olindi", todos)
}

// GetByID godoc
// @Summary      Todo'ni ID bo'yicha olish
// @Description  Berilgan ID bo'yicha bitta todo'ni qaytaradi
// @Tags         todos
// @Produce      json
// @Param        id   path      string  true  "Todo ID"
// @Success      200  {object}  response.Response
// @Failure      404  {object}  response.ErrorResponse
// @Router       /todos/{id} [get]
func (h *TodoHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	todo, err := h.usecase.GetByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.NotFound(c, "Todo topilmadi")
			return
		}
		response.InternalError(c, "Todo'ni olishda xato yuz berdi")
		return
	}
	response.OK(c, "Todo muvaffaqiyatli olindi", todo)
}

// Create godoc
// @Summary      Yangi todo yaratish
// @Description  Yangi todo yaratadi va qaytaradi
// @Tags         todos
// @Accept       json
// @Produce      json
// @Param        request  body      request.CreateTodoRequest  true  "Todo ma'lumotlari"
// @Success      201      {object}  response.Response
// @Failure      400      {object}  response.ErrorResponse
// @Router       /todos [post]
func (h *TodoHandler) Create(c *gin.Context) {
	var req request.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Noto'g'ri so'rov ma'lumotlari: "+err.Error())
		return
	}
	if len(req.Title) < 3 {
		response.BadRequest(c, "Noto'g'ri so'rov ma'lumotlari: "+" title must be at least 3 characters long")
		return
	}

	todo, err := h.usecase.Create(req)
	if err != nil {
		response.InternalError(c, "Todo yaratishda xato yuz berdi")
		return
	}
	response.Created(c, "Todo muvaffaqiyatli yaratildi", todo)
}

// Update godoc
// @Summary      Todo'ni yangilash
// @Description  Mavjud todo'ni yangilaydi
// @Tags         todos
// @Accept       json
// @Produce      json
// @Param        id       path      string                 true  "Todo ID"
// @Param        request  body      request.UpdateTodoRequest  true  "Yangilanayotgan ma'lumotlar"
// @Success      200      {object}  response.Response
// @Failure      400      {object}  response.ErrorResponse
// @Failure      404      {object}  response.ErrorResponse
// @Router       /todos/{id} [patch]
func (h *TodoHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req request.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Noto'g'ri so'rov ma'lumotlari: "+err.Error())
		return
	}
	if len(*req.Title) < 3 {
		response.BadRequest(c, "Noto'g'ri so'rov ma'lumotlari: "+" title must be at least 3 characters long")
		return
	}

	todo, err := h.usecase.Update(id, req)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.NotFound(c, "Todo topilmadi")
			return
		}
		response.InternalError(c, "Todo yangilashda xato yuz berdi")
		return
	}
	response.OK(c, "Todo muvaffaqiyatli yangilandi", todo)
}

// Delete godoc
// @Summary      Todo'ni o'chirish
// @Description  Berilgan ID bo'yicha todo'ni o'chiradi
// @Tags         todos
// @Produce      json
// @Param        id   path      string  true  "Todo ID"
// @Success      200  {object}  response.Response
// @Failure      404  {object}  response.ErrorResponse
// @Router       /todos/{id} [delete]
func (h *TodoHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.usecase.Delete(id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			response.NotFound(c, "Todo topilmadi")
			return
		}
		response.InternalError(c, "Todo o'chirishda xato yuz berdi")
		return
	}
	response.OK(c, "Todo muvaffaqiyatli o'chirildi", nil)
}
