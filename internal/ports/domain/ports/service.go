package ports

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

type Repository interface {
	RepositoryWriter
	RepositoryReader
}

type RepositoryWriter interface {
	Store(context.Context, *Port) error
}

type RepositoryReader interface {
	Get(context.Context) ([]Port, error)
	FindByID(context.Context, string) *Port
}

type Service struct {
	logger     *logrus.Logger
	repository Repository
}

func NewService(logger *logrus.Logger, repository Repository) Service {
	return Service{
		logger:     logger,
		repository: repository,
	}
}

func (s Service) Store(ctx context.Context, port *Port) error {
	if port == nil {
		return fmt.Errorf("port is nil")
	}

	if port.PortID == "" || port.Name == "" {
		return fmt.Errorf("port id and port name is required")
	}

	if err := s.repository.Store(ctx, port); err != nil {
		return fmt.Errorf("error storing port err: %w", err)
	}

	s.logger.Infof("successfully storing port: %+v", port)

	return nil
}

func (s Service) Get(ctx context.Context) ([]Port, error) {
	s.logger.Info("service Get called ")

	ports, err := s.repository.Get(ctx)
	return ports, err
}
