package main

import (
	"context"
	"log"
	"strings"

	"github.com/opam22/ports/internal/ports"
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

	if config.GetString("ports.port") == "" {
		logger.Fatal("missing ports port")
	}
}

func main() {
	server := ports.NewServer(logger, config)

	ctx, _ := context.WithCancel(context.Background())
	if err := server.Serve(ctx); err != nil {
		log.Println("fatal to serve the server", err)
	}
}
