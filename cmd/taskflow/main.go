package main

import (
	"fmt"

	"taskflow/internal/config"
	"taskflow/internal/service"
	"taskflow/internal/store"
)

func main() {
	cfg := config.Load()
	st := store.New()
	svc := service.New(st)
	_ = cfg
	_ = svc
	fmt.Println("taskflow ready")
}
