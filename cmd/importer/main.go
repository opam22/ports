package main

import (
	"context"
	"log"

	"github.com/opam22/ports/internal/importer"
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
	importer, err := importer.NewService(logger)
	if err != nil {
		log.Println(err)
	}

	ctx, _ := context.WithCancel(context.Background())
	err = importer.Run(ctx)
	if err != nil {
		log.Println(err)
	}

	log.Println("importer done")
}
