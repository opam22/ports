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

func main() {
	importer, err := importer.NewService()
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
