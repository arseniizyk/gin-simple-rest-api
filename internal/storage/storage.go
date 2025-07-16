package storage

import (
	"errors"
	"sync"
)

var (
	ErrNotExists = errors.New("employee with this id does not exist")
)

type Employee struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Sex    string  `json:"sex"`
	Age    int     `json:"age"`
	Salary float64 `json:"salary"`
}

type Storage interface {
	Insert(e *Employee)
	Get(id int) (*Employee, error)
	Update(id int, e Employee) error
	Delete(id int) error
}

type MemoryStorage struct {
	counter int
	data    map[int]Employee
	sync.RWMutex
}

func NewMemoryStorage() Storage {
	return &MemoryStorage{
		counter: 1,
		data:    make(map[int]Employee),
	}
}

func (s *MemoryStorage) Insert(e *Employee) {
	s.Lock()
	defer s.Unlock()

	e.ID = s.counter
	s.data[e.ID] = *e

	s.counter++
}

func (s *MemoryStorage) Get(id int) (*Employee, error) {
	s.RLock()
	defer s.RUnlock()

	e, exists := s.data[id]
	if !exists {
		return nil, ErrNotExists
	}

	return &e, nil
}

func (s *MemoryStorage) Update(id int, e Employee) error {
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
