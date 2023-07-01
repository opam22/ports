package importer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/opam22/ports/internal/portsclient"
	"github.com/sirupsen/logrus"
)

type Service struct {
	logger      *logrus.Logger
	portsClient portsclient.Client
}

func NewService(logger *logrus.Logger) (*Service, error) {
	client, err := portsclient.NewClient(logger)
	if err != nil {
		return &Service{}, fmt.Errorf("faiiled to initializing importer service %w", err)
	}

	return &Service{
		logger:      logger,
		portsClient: *client,
	}, nil
}

func (s *Service) Run(ctx context.Context) (err error) {
	file, err := os.Open("internal/importer/json/ports.json")
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
