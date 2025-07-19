package employees

import (
	"context"

	"github.com/arseniizyk/internal/models"
)

type EmployeeUsecase interface {
	GetAll(ctx context.Context) ([]models.Employee, error)
	Add(ctx context.Context, e *models.Employee) error
	GetByID(ctx context.Context, id int) (*models.Employee, error)
	Update(ctx context.Context, id int, e *models.Employee) error
	Delete(ctx context.Context, id int) error
}
