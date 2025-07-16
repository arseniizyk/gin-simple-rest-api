package storage

import (
	"context"
	"errors"

	"github.com/arseniizyk/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotExists = errors.New("models.Employee with this id does not exist")
)

type Storage interface {
	Insert(ctx context.Context, e *models.Employee) error
	Get(ctx context.Context, id int) (*models.Employee, error)
	Update(ctx context.Context, id int, e *models.Employee) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) ([]models.Employee, error)
}

type PostgreSQL struct {
	pool *pgxpool.Pool
}

func New(p *pgxpool.Pool) Storage {
	return &PostgreSQL{
		pool: p,
	}
}

func (p *PostgreSQL) List(ctx context.Context) ([]models.Employee, error) {
	query := `SELECT id, name, sex, age, salary FROM employees`
	rows, err := p.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	employees := make([]models.Employee, 0)

	for rows.Next() {
		var e models.Employee
		err := rows.Scan(&e.ID, &e.Name, &e.Sex, &e.Age, &e.Salary)
		if err != nil {
			return nil, err
		}

		employees = append(employees, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return employees, nil
}

func (s *PostgreSQL) Insert(ctx context.Context, e *models.Employee) error {
	query := `INSERT INTO employees (name, sex, age, salary) VALUES ($1, $2, $3, $4) RETURNING id`

	err := s.pool.QueryRow(ctx, query, e.Name, e.Sex, e.Age, e.Salary).Scan(&e.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgreSQL) Get(ctx context.Context, id int) (*models.Employee, error) {
	query := `SELECT id, name, sex, age, salary FROM employees WHERE id = $1`

	e := new(models.Employee)
	err := s.pool.QueryRow(ctx, query, id).Scan(&e.ID, &e.Name, &e.Sex, &e.Age, &e.Salary)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}

	return e, nil
}

func (s *PostgreSQL) Update(ctx context.Context, id int, e *models.Employee) error {
	query := `UPDATE employees SET name = $1, sex = $2, age = $3, salary = $4 WHERE id = $5`

	cmdTag, err := s.pool.Exec(ctx, query, e.Name, e.Sex, e.Age, e.Salary, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrNotExists
	}

	return nil
}

func (s *PostgreSQL) Delete(ctx context.Context, id int) error {
	cmdTag, err := s.pool.Exec(ctx, `DELETE FROM employees WHERE id = $1`, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrNotExists
	}

	return nil
}
