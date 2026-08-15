package store

import (
	"errors"
	"sync"

	"taskflow/internal/model"
)

var (
	ErrNotFound      = errors.New("task not found")
	ErrAlreadyExists = errors.New("task already exists")
)

type Store struct {
	mu        sync.RWMutex
	tasks     map[string]*model.Task
	order     []string
	processed int
}

func New() *Store {
	return &Store{
		tasks: make(map[string]*model.Task),
		order: []string{},
	}
}

func (s *Store) Put(t *model.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tasks[t.ID]; ok {
		return ErrAlreadyExists
	}
	s.tasks[t.ID] = t
	s.order = append(s.order, t.ID)
	return nil
}

func (s *Store) Get(id string) (*model.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tasks[id]
	if !ok {
		return nil, ErrNotFound
	}
	return t, nil
}

func (s *Store) List() []*model.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*model.Task, 0, len(s.order))
	for _, id := range s.order {
		out = append(out, s.tasks[id])
	}
	return out
}

func (s *Store) IDs() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]string, len(s.order))
	copy(out, s.order)
	return out
}

func (s *Store) UpdateStatus(id string, status model.Status, result string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.tasks[id]
	if !ok {
		return ErrNotFound
	}
	t.Status = status
	t.Result = result
	return nil
}

func (s *Store) IncrementProcessed() {
	s.processed++
}

func (s *Store) Processed() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.processed
}
