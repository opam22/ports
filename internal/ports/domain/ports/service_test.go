package ports

import (
	"context"
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type MockRepository struct {
	ports map[string]*Port
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		ports: make(map[string]*Port),
	}
}

func (r *MockRepository) Store(ctx context.Context, port *Port) error {
	if port == nil {
		return fmt.Errorf("port is nil")
	}

	r.ports[port.PortID] = port
	return nil
}

func (r *MockRepository) Get(ctx context.Context) ([]Port, error) {
	ports := make([]Port, 0, len(r.ports))
	for _, p := range r.ports {
		ports = append(ports, *p)
	}
	return ports, nil
}

func (r *MockRepository) FindByID(ctx context.Context, portID string) *Port {
	return r.ports[portID]
}

func TestService_Store(t *testing.T) {
	logger := logrus.New()
	mockRepo := NewMockRepository()

	service := NewService(logger, mockRepo)

	port := &Port{
		PortID: "AEAJM",
		Name:   "Ajman",
	}
	err := service.Store(context.Background(), port)
	assert.NoError(t, err)

	tests := []struct {
		name        string
		port        *Port
		shouldError bool
		errMsg      string
	}{
		{
			name:        "no port given",
			port:        nil,
			shouldError: true,
			errMsg:      "port is nil",
		},
		{
			name: "missing id and name",
			port: &Port{
				City: "Jakarta",
			},
			shouldError: true,
			errMsg:      "port id and port name is required",
		},
		{
			name: "success",
			port: &Port{
				PortID: "AEAJM",
				Name:   "Ajman",
			},
			shouldError: false,
		},
	}

	for _, tt := range tests {
		if err = service.Store(context.Background(), tt.port); err != nil {
			if !tt.shouldError {
				t.Fatalf("expecting success but got error: %+v", err)
			}

			assert.EqualError(t, err, tt.errMsg)
		}
	}
}
