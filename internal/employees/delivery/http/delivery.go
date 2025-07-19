package delivery

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/arseniizyk/internal/employees"
	"github.com/arseniizyk/internal/employees/repository/postgres"
	"github.com/arseniizyk/internal/models"
	"github.com/arseniizyk/internal/utils"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc employees.EmployeeUsecase
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func New(uc employees.EmployeeUsecase) *Handler {
	return &Handler{
		uc: uc,
	}
}

func (h *Handler) CreateEmployee(c *gin.Context) {
	emp := new(models.Employee)

	if err := c.BindJSON(emp); err != nil {
		sendError(c, http.StatusBadRequest, "failed to bind employee", err)
		return
	}

	if err := h.uc.Add(c.Request.Context(), emp); err != nil {
		sendError(c, http.StatusInternalServerError, "failed to add employee", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": emp.ID})
}

func (h *Handler) GetAllEmployees(c *gin.Context) {
	employees, err := h.uc.GetAll(c.Request.Context())
	if err != nil {
		sendError(c, http.StatusInternalServerError, "failed to get all employees", err)
		return
	}

	c.JSON(http.StatusOK, employees)
}

func (h *Handler) GetEmployee(c *gin.Context) {
	id, err := utils.ParseID(&c.Params)
	if err != nil {
		sendError(c, http.StatusBadRequest, "cant parse id to int, probably you provide non-int value", err)
		return
	}

	employee, err := h.uc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, postgres.ErrNotExists) {
			sendError(c, http.StatusNotFound, "employee not found with id"+fmt.Sprint(id), err)
			return
		}
		sendError(c, http.StatusInternalServerError, "cant get employee with id"+fmt.Sprint(id), err)
		return
	}

	c.JSON(http.StatusOK, employee)
}

func (h *Handler) UpdateEmployee(c *gin.Context) {
	id, err := utils.ParseID(&c.Params)
	if err != nil {
		sendError(c, http.StatusBadRequest, "cant parse id to int, probably you provide non-int value", err)
		return
	}

	emp := new(models.Employee)

	if err := c.BindJSON(emp); err != nil {
		sendError(c, http.StatusBadRequest, "failed to bind employee", err)
		return
	}

	if err := h.uc.Update(c.Request.Context(), id, emp); err != nil {
		sendError(c, http.StatusNotFound, "failed to update employee, not found", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *Handler) DeleteEmployee(c *gin.Context) {
	id, err := utils.ParseID(&c.Params)
	if err != nil {
		sendError(c, http.StatusBadRequest, "cant parse id to int, probably you provide non-int value", err)
		return
	}

	if err := h.uc.Delete(c.Request.Context(), id); err != nil {
		sendError(c, http.StatusNotFound, "employee not found with id"+fmt.Sprint(id), err)
		return
	}

	c.Status(http.StatusOK)
}
