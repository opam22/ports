package main

import (
	"context"
	"log"

	"github.com/opam22/ports/internal/ports"
	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Logger
)

func init() {
	logger = logrus.New()
	logger.Formatter = &logrus.JSONFormatter{PrettyPrint: true}
}

func main() {
	server := ports.NewServer(logger)

	ctx, _ := context.WithCancel(context.Background())
	if err := server.Serve(ctx); err != nil {
		log.Println("fatal to serve the server", err)
	}
}
