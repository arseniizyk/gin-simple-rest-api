package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/arseniizyk/internal/storage"
	"github.com/arseniizyk/internal/utils"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type Handler struct {
	s storage.Storage
}

func New(s storage.Storage) *Handler {
	return &Handler{s: s}
}

func (h *Handler) CreateEmployee(c *gin.Context) {
	var employee storage.Employee

	if err := c.BindJSON(&employee); err != nil {
		log.Println("failed to bind employee", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   err.Error(),
			Message: "bad request",
		})
		return
	}

	h.s.Insert(&employee)

	c.JSON(http.StatusCreated, gin.H{"id": employee.ID})
}

func (h *Handler) GetEmployee(c *gin.Context) {
	id, err := utils.ParseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   err.Error(),
			Message: "cant parse id to int, probably you provide non-int value",
		})
		return
	}

	employee, err := h.s.Get(id)
	if err != nil {
		if errors.Is(err, storage.ErrNotExists) {
			log.Printf("employee with id %d not found\n", id)
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   err.Error(),
				Message: "employee not found",
			})
			return
		}
		log.Printf("cant get employee with id %d\n", id)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   err.Error(),
			Message: "unexpected error",
		})
		return
	}

	c.JSON(http.StatusOK, employee)
}

func (h *Handler) UpdateEmployee(c *gin.Context) {
	id, err := utils.ParseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   err.Error(),
			Message: "cant parse id to int, probably you provide non-int value",
		})
	}

	var employee storage.Employee

	if err := c.BindJSON(&employee); err != nil {
		log.Println("failed to bind employee", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   err.Error(),
			Message: "bad request",
		})
		return
	}

	if err := h.s.Update(id, employee); err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   err.Error(),
			Message: "employee not found",
		})
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *Handler) DeleteEmployee(c *gin.Context) {
	id, err := utils.ParseID(c)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   err.Error(),
			Message: "cant parse id to int, probably you provide non-int value",
		})
	}

	if err := h.s.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   err.Error(),
			Message: "employee not found",
		})
	}

	c.Status(http.StatusOK)
}
