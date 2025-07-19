package delivery

import (
	"github.com/arseniizyk/internal/employees"
	"github.com/gin-gonic/gin"
)

// GET /employee/ - получение всех зарегистрированных сотрудников
// POST /employee/ - cоздание сотрудника
// GET /employee/<id> - получение информации о сотруднике
// PUT /employee/<id> - редактирование информации о сотруднике
// DELETE /employee/<id> - удаление сотрудника
func RegisterEmployeesEndpoints(r *gin.RouterGroup, uc employees.EmployeeUsecase) {
	h := New(uc)

	r.GET("/employees", h.GetAllEmployees)
	r.POST("/employees", h.CreateEmployee)
	r.GET("/employees/:id", h.GetEmployee)
	r.PUT("/employees/:id", h.UpdateEmployee)
	r.DELETE("/employees/:id", h.DeleteEmployee)
}
