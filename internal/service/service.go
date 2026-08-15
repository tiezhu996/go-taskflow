package service

import (
	"context"
	"errors"
	"fmt"

	"taskflow/internal/model"
	"taskflow/internal/store"
)

var ErrEmptyPayload = errors.New("empty payload")

type Service struct {
	store *store.Store
}

func New(s *store.Store) *Service {
	return &Service{store: s}
}

func (svc *Service) Submit(id, payload string) (*model.Task, error) {
	if payload == "" {
		return nil, ErrEmptyPayload
	}
	t := &model.Task{ID: id, Payload: payload, Status: model.StatusPending}
	if err := svc.store.Put(t); err != nil {
		return nil, fmt.Errorf("submit task %s: %w", id, err)
	}
	return t, nil
}

func (svc *Service) Get(id string) (*model.Task, error) {
	t, err := svc.store.Get(id)
	if err != nil {
		return nil, fmt.Errorf("get task %s: %w", id, err)
	}
	return t, nil
}

func (svc *Service) ListIDs() []string {
	return svc.store.IDs()
}

// Process implements worker.Processor so a Service can be used as the task
// processor inside a worker pool.
func (svc *Service) Process(ctx context.Context, t *model.Task) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	return t.Payload + "-processed", nil
}
