package ports

import (
	"context"
	"fmt"
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
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{
		repository: repository,
	}
}

func (s Service) Store(ctx context.Context, port *Port) error {
	if port == nil {
		return fmt.Errorf("nil port given")
	}

	if port.PortID == "" || port.Name == "" {
		return fmt.Errorf("port id and port name is required")
	}

	if err := s.repository.Store(ctx, port); err != nil {
		return fmt.Errorf("error storing port err: %w", err)
	}

	return nil
}

func (s Service) Get(ctx context.Context) ([]Port, error) {
	ports, err := s.repository.Get(ctx)
	return ports, err
}
