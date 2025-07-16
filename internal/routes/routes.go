package routes

import (
	"github.com/arseniizyk/internal/handlers"
	"github.com/gin-gonic/gin"
)

func Setup(h *handlers.Handler) *gin.Engine {
	r := gin.Default()

	api := r.Group("/v1")
	api.GET("/employee/", h.ListEmployees)
	api.POST("/employee", h.CreateEmployee)
	api.GET("/employee/:id", h.GetEmployee)
	api.PUT("/employee/:id", h.UpdateEmployee)
	api.DELETE("/employee/:id", h.DeleteEmployee)

	return r
}

// GET /employee/ - получение всех зарегистрированных сотрудников
// POST /employee/ - cоздание сотрудника
// GET /employee/<id> - получение информации о сотруднике
// PUT /employee/<id> - редактирование информации о сотруднике
// DELETE /employee/<id> - удаление сотрудника
