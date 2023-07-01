package main

import (
	"context"
	"log"

	"github.com/opam22/ports/internal/ports"
)

func main() {
	server := ports.NewServer()

	ctx, _ := context.WithCancel(context.Background())
	if err := server.Serve(ctx); err != nil {
		log.Println("fatal to serve the server", err)
	}
}
