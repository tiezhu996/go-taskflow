package worker

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"taskflow/internal/model"
	"taskflow/internal/service"
	"taskflow/internal/store"
)

func TestPoolProcessesAll(t *testing.T) {
	st := store.New()
	svc := service.New(st)
	for i := 0; i < 200; i++ {
		if _, err := svc.Submit(fmt.Sprintf("t%d", i), "p"); err != nil {
			t.Fatal(err)
		}
	}
	ids := st.IDs()
	pool := New(st, svc, 8)
	pool.Run(context.Background(), ids)
	if got := st.Processed(); got != 200 {
		t.Fatalf("processed=%d want 200", got)
	}
	for _, id := range ids {
		tk, _ := st.Get(id)
		if tk.Status != model.StatusDone {
			t.Fatalf("task %s status=%q want done", id, tk.Status)
		}
	}
}

type countingProcessor struct {
	mu     sync.Mutex
	calls  int
	once   func()
	onceDo sync.Once
}

func (p *countingProcessor) Process(ctx context.Context, t *model.Task) (string, error) {
	p.onceDo.Do(func() {
		if p.once != nil {
			p.once()
		}
	})
	p.mu.Lock()
	p.calls++
	p.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return "", err
	}
	return t.Payload + "-done", nil
}

func TestPoolStopsOnCancel(t *testing.T) {
	st := store.New()
	for i := 0; i < 100; i++ {
		if err := st.Put(&model.Task{ID: fmt.Sprintf("t%d", i), Payload: "p", Status: model.StatusPending}); err != nil {
			t.Fatal(err)
		}
	}
	proc := &countingProcessor{}
	pool := New(st, proc, 4)
	ctx, cancel := context.WithCancel(context.Background())
	proc.once = func() { cancel() }
	pool.Run(ctx, st.IDs())
	if proc.calls >= 100 {
		t.Fatalf("pool kept processing after cancel: calls=%d", proc.calls)
	}
}
