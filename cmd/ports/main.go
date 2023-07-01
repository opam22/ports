package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

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

	ctx, cancel := context.WithCancel(context.Background())
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go func(cancel context.CancelFunc) {
		s := <-signalCh
		logger.Info("gracefully shutdown ports service: signal: ", s)
		cancel()
	}(cancel)

	if err := server.Serve(ctx); err != nil {
		logger.Fatal("fatal to serve the server", err)
	}
}
