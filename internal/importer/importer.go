package importer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	gRPC "github.com/opam22/ports/internal/ports/grpc"
	"github.com/opam22/ports/internal/portsclient"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type PortsClient interface {
	PortsClientWriter
	PortsClientReader
}

type PortsClientWriter interface {
	Store(ctx context.Context, port *gRPC.Port) error
	Close() error
}

type PortsClientReader interface {
	Get(ctx context.Context) ([]*gRPC.Port, error)
}

type Service struct {
	logger      *logrus.Logger
	portsClient PortsClient
	config      *viper.Viper
	filePath    string
}

func NewService(logger *logrus.Logger, config *viper.Viper) (*Service, error) {
	client, err := portsclient.NewClient(logger, config)
	if err != nil {
		return &Service{}, fmt.Errorf("faiiled to initializing importer service %w", err)
	}

	return &Service{
		logger:      logger,
		portsClient: client,
		filePath:    config.GetString("importer.filePath"),
	}, nil
}

func (s *Service) Run(ctx context.Context) (err error) {
	// graceful shutdown
	go func() {
		<-ctx.Done()
		if err := s.portsClient.Close(); err != nil {
			s.logger.Error("error closing client connection: %w", err)
		}
		return
	}()

	file, err := os.Open(s.filePath)
	if err != nil {
		s.logger.Error("error opening file: %w", err)
		return errors.New("fail to read file")
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	_, err = decoder.Token()
	if err != nil {
		s.logger.Error("error decoding JSON: %w", err)
		return errors.New("fail to read json")
	}

	for decoder.More() {
		key, err := decoder.Token()
		if err != nil {
			s.logger.Error("error decoding json: %w", err)
			continue
		}

		port := Port{}
		rawMsg := json.RawMessage{}
		err = decoder.Decode(&rawMsg)
		if err != nil {
			s.logger.Error("error decoding json: %w", err)
			continue
		}

		err = json.Unmarshal(rawMsg, &port)
		if err != nil {
			s.logger.Error("error decoding json: %w", err)
			continue
		}

		if err := s.portsClient.Store(ctx, toProtoPort(
			fmt.Sprintf("%s", key),
			port,
		)); err != nil {
			s.logger.Error(err)
			continue
		}
	}

	_, err = decoder.Token()
	if err != nil {
		s.logger.Error("error decoding json: %w", err)
		return
	}

	if err := s.portsClient.Close(); err != nil {
		s.logger.Error("error closing client connection: %w", err)
		return err
	}

	return nil
}
