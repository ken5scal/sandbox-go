package main

import (
	"context"
	"fmt"
	"github.com/ken5scal/go_todo_app/config"
	"log"
	"net"
)

func run(ctx context.Context) error {
	cfg, err := config.New()
	if err != nil {
		return err
	}
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", cfg.Port, err)
	}
	mux := NewMux() // NewMux()
	s := NewServer(l, mux)

	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with: %v", url)
	return s.Run(ctx)
}

func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminate server: %v", err)
	}
}
