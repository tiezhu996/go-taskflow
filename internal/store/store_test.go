package store

import (
	"testing"

	"taskflow/internal/model"
)

func TestPutGetList(t *testing.T) {
	s := New()
	if err := s.Put(&model.Task{ID: "1", Status: model.StatusPending}); err != nil {
		t.Fatal(err)
	}
	if err := s.Put(&model.Task{ID: "2", Status: model.StatusPending}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Get("1"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Get("missing"); err != ErrNotFound {
		t.Fatalf("err=%v want ErrNotFound", err)
	}
	got := s.List()
	if len(got) != 2 {
		t.Fatalf("List len=%d want 2", len(got))
	}
	if err := s.Put(&model.Task{ID: "1"}); err != ErrAlreadyExists {
		t.Fatalf("dup err=%v want ErrAlreadyExists", err)
	}
}

func TestIDsFreshCopy(t *testing.T) {
	s := New()
	for _, id := range []string{"b", "a", "c"} {
		if err := s.Put(&model.Task{ID: id, Status: model.StatusPending}); err != nil {
			t.Fatal(err)
		}
	}
	first := s.IDs()
	// mutate what the caller got back; internal order must be untouched
	for i := range first {
		first[i] = "x"
	}
	after := s.IDs()
	want := []string{"b", "a", "c"}
	for i := range want {
		if after[i] != want[i] {
			t.Fatalf("IDs() after=%v want %v", after, want)
		}
	}
}

func TestProcessedCounter(t *testing.T) {
	s := New()
	for i := 0; i < 10; i++ {
		s.IncrementProcessed()
	}
	if s.Processed() != 10 {
		t.Fatalf("Processed=%d want 10", s.Processed())
	}
}
