package delivery

import (
	"github.com/arseniizyk/internal/employees"
	"github.com/gin-gonic/gin"
)

// GET /employees/ - получение всех зарегистрированных сотрудников
// POST /employees/ - cоздание сотрудника
// GET /employees/{id} - получение информации о сотруднике
// PUT /employees/{id} - редактирование информации о сотруднике
// DELETE /employees/{id} - удаление сотрудника
func RegisterEmployeesEndpoints(r *gin.RouterGroup, uc employees.EmployeeUsecase) {
	h := New(uc)

	r.GET("/employees", h.GetAllEmployees)
	r.POST("/employees", h.CreateEmployee)
	r.GET("/employees/:id", h.GetEmployee)
	r.PUT("/employees/:id", h.UpdateEmployee)
	r.DELETE("/employees/:id", h.DeleteEmployee)
}
