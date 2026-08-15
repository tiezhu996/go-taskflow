package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("TASKFLOW_WORKERS", "")
	c := Load()
	if c.Workers != 2 {
		t.Fatalf("Workers=%d want 2", c.Workers)
	}
}

func TestLoadEnv(t *testing.T) {
	t.Setenv("TASKFLOW_WORKERS", "7")
	c := Load()
	if c.Workers != 7 {
		t.Fatalf("Workers=%d want 7", c.Workers)
	}
}
