package storage

import (
	"errors"
	"sync"

	"github.com/arseniizyk/internal/models"
)

var (
	ErrNotExists = errors.New("models.Employee with this id does not exist")
)

type Storage interface {
	Insert(e *models.Employee)
	Get(id int) (*models.Employee, error)
	Update(id int, e models.Employee) error
	Delete(id int) error
	List() []models.Employee
}

type MemoryStorage struct {
	counter int
	data    map[int]models.Employee
	sync.RWMutex
}

func NewMemoryStorage() Storage {
	return &MemoryStorage{
		counter: 1,
		data:    make(map[int]models.Employee),
	}
}

func (s *MemoryStorage) List() []models.Employee {
	employees := make([]models.Employee, 0, len(s.data))
	for _, emp := range s.data {
		employees = append(employees, emp)
	}

	return employees
}

func (s *MemoryStorage) Insert(e *models.Employee) {
	s.Lock()
	defer s.Unlock()

	e.ID = s.counter
	s.data[e.ID] = *e

	s.counter++
}

func (s *MemoryStorage) Get(id int) (*models.Employee, error) {
	s.RLock()
	defer s.RUnlock()

	e, exists := s.data[id]
	if !exists {
		return nil, ErrNotExists
	}

	return &e, nil
}

func (s *MemoryStorage) Update(id int, e models.Employee) error {
	s.Lock()
	defer s.Unlock()

	_, ok := s.data[id]
	if !ok {
		return ErrNotExists
	}

	e.ID = id
	s.data[id] = e
	return nil
}

func (s *MemoryStorage) Delete(id int) error {
	s.Lock()
	defer s.Unlock()

	_, ok := s.data[id]
	if !ok {
		return ErrNotExists
	}

	delete(s.data, id)
	return nil
}
