package service

import (
	"context"
	"errors"
	"sort"
	"testing"

	"taskflow/internal/model"
	"taskflow/internal/store"
)

func TestSubmitGet(t *testing.T) {
	s := store.New()
	svc := New(s)
	tk, err := svc.Submit("1", "hello")
	if err != nil {
		t.Fatal(err)
	}
	if tk.Status != model.StatusPending {
		t.Fatalf("status=%q want pending", tk.Status)
	}
	if _, err := svc.Submit("2", ""); !errors.Is(err, ErrEmptyPayload) {
		t.Fatalf("empty payload err=%v", err)
	}
	got, err := svc.Get("1")
	if err != nil {
		t.Fatal(err)
	}
	if got.Payload != "hello" {
		t.Fatalf("payload=%q", got.Payload)
	}
}

func TestListIDsIsolated(t *testing.T) {
	s := store.New()
	svc := New(s)
	for _, id := range []string{"b", "a", "c"} {
		if _, err := svc.Submit(id, "p"); err != nil {
			t.Fatal(err)
		}
	}
	ids := svc.ListIDs()
	sort.Strings(ids)
	after := svc.ListIDs()
	want := []string{"b", "a", "c"}
	for i := range want {
		if after[i] != want[i] {
			t.Fatalf("ListIDs after=%v want %v (returned slice aliases internal state)", after, want)
		}
	}
}

func TestGetNotFoundWraps(t *testing.T) {
	s := store.New()
	svc := New(s)
	_, err := svc.Get("nope")
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("errors.Is(err, ErrNotFound)=false, err=%v", err)
	}
}

func TestProcessContext(t *testing.T) {
	s := store.New()
	svc := New(s)
	out, err := svc.Process(context.Background(), &model.Task{Payload: "x"})
	if err != nil || out != "x-processed" {
		t.Fatalf("out=%q err=%v", out, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := svc.Process(ctx, &model.Task{Payload: "x"}); err == nil {
		t.Fatal("expected error for cancelled context")
	}
}
