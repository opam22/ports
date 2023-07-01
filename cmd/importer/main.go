package main

import (
	"context"
	"log"
	"strings"

	"github.com/opam22/ports/internal/importer"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	logger *logrus.Logger
	config *viper.Viper
)

func init() {
	logger = logrus.New()
	logger.Formatter = &logrus.JSONFormatter{PrettyPrint: true}

	config = viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	config.SetConfigFile("./env/config")
	config.SetConfigType("ini")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("error loading configuration: %v", err)
	}

	if config.GetString("importer.serverPort") == "" {
		logger.Fatal("missing importer serverPort")
	}

}

func main() {
	importer, err := importer.NewService(logger, config)
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
