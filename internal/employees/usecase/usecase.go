package usecase

import (
	"context"

	"github.com/arseniizyk/internal/employees"
	repo "github.com/arseniizyk/internal/employees"
	"github.com/arseniizyk/internal/models"
)

type employeeUsecase struct {
	repo repo.EmployeeRepo
}

func New(repo repo.EmployeeRepo) employees.EmployeeUsecase {
	return &employeeUsecase{repo: repo}
}

func (uc *employeeUsecase) GetAll(ctx context.Context) ([]models.Employee, error) {
	return uc.repo.GetAll(ctx)
}

func (uc *employeeUsecase) Add(ctx context.Context, e *models.Employee) error {
	return uc.repo.Insert(ctx, e)
}

func (uc *employeeUsecase) Update(ctx context.Context, id int, e *models.Employee) error {
	return uc.repo.Update(ctx, id, e)
}

func (uc *employeeUsecase) Delete(ctx context.Context, id int) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *employeeUsecase) GetByID(ctx context.Context, id int) (*models.Employee, error) {
	return uc.repo.GetByID(ctx, id)
}
