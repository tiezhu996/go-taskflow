package worker

import (
	"context"
	"sync"

	"taskflow/internal/model"
	"taskflow/internal/store"
)

type Processor interface {
	Process(ctx context.Context, t *model.Task) (string, error)
}

type Pool struct {
	store   *store.Store
	proc    Processor
	workers int
}

func New(s *store.Store, p Processor, workers int) *Pool {
	if workers <= 0 {
		workers = 1
	}
	return &Pool{store: s, proc: p, workers: workers}
}

func (p *Pool) Run(ctx context.Context, ids []string) {
	var wg sync.WaitGroup
	ch := make(chan string, len(ids))

	go func() {
		defer close(ch)
		for _, id := range ids {
			select {
			case <-ctx.Done():
				return
			case ch <- id:
			}
		}
	}()

	for i := 0; i < p.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for id := range ch {
				select {
				case <-ctx.Done():
					return
				default:
				}
				t, err := p.store.Get(id)
				if err != nil {
					continue
				}
				out, err := p.proc.Process(ctx, t)
				if err != nil {
					_ = p.store.UpdateStatus(t.ID, model.StatusFailed, err.Error())
					continue
				}
				_ = p.store.UpdateStatus(t.ID, model.StatusDone, out)
				p.store.IncrementProcessed()
			}
		}()
	}

	wg.Wait()
}
