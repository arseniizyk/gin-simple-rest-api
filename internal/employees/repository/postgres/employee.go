package postgres

import (
	"context"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/arseniizyk/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotExists = errors.New("models.Employee with this id does not exist")
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type Repo struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repo {

	return &Repo{pool: pool}
}

func (r *Repo) GetAll(ctx context.Context) ([]models.Employee, error) {
	query, args, err := psql.Select("id", "name", "sex", "age", "salary").
		From("employees").ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	employees := make([]models.Employee, 0)

	for rows.Next() {
		var e models.Employee
		if err := rows.Scan(&e.ID, &e.Name, &e.Sex, &e.Age, &e.Salary); err != nil {
			return nil, err
		}

		employees = append(employees, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return employees, nil
}

func (r *Repo) Insert(ctx context.Context, e *models.Employee) error {
	query, args, err := psql.Insert("employees").
		Columns("name", "sex", "age", "salary").
		Values(e.Name, e.Sex, e.Age, e.Salary).
		Suffix("RETURNING \"id\"").ToSql()
	if err != nil {
		return err
	}

	if err := r.pool.QueryRow(ctx, query, args...).Scan(&e.ID); err != nil {
		return err
	}

	return nil
}

func (r *Repo) GetByID(ctx context.Context, id int) (*models.Employee, error) {
	query, args, err := psql.Select("id", "name", "sex", "age", "salary").From("employees").Where("id = $1", id).ToSql()
	if err != nil {
		return nil, err
	}

	var e models.Employee
	if err := r.pool.QueryRow(ctx, query, args...).Scan(&e.ID, &e.Name, &e.Sex, &e.Age, &e.Salary); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotExists
		}

		return nil, err
	}

	return &e, nil
}

func (r *Repo) Update(ctx context.Context, id int, e *models.Employee) error {
	query, args, err := psql.Update("employees").
		Set("name", e.Name).
		Set("sex", e.Sex).
		Set("age", e.Age).
		Set("salary", e.Salary).Where("id = ?", id).ToSql()
	if err != nil {
		return err
	}

	cmdTag, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrNotExists
	}

	return nil
}

func (r *Repo) Delete(ctx context.Context, id int) error {
	query, args, err := psql.Delete("employees").Where("id = $1", id).ToSql()
	if err != nil {
		return err
	}

	cmdTag, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrNotExists
	}

	return nil
}
