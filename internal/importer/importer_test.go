package importer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	gRPC "github.com/opam22/ports/internal/ports/grpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPortsClient struct {
	mock.Mock
}

func (m *MockPortsClient) Store(ctx context.Context, port *gRPC.Port) error {
	args := m.Called(ctx, port)
	return args.Error(0)
}

func (m *MockPortsClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockPortsClient) Get(ctx context.Context) ([]*gRPC.Port, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*gRPC.Port), args.Error(1)
}

func TestService_Run(t *testing.T) {
	mockPortsClient := new(MockPortsClient)

	service := &Service{
		portsClient: mockPortsClient,
		filePath:    "./json/ports.json",
	}

	file, err := os.Open(service.filePath)
	assert.NoError(t, err)
	defer file.Close()

	mockPortsClient.On("Store", mock.Anything, mock.Anything).Return(nil)
	mockPortsClient.On("Close").Return(nil)

	decoder := json.NewDecoder(file)
	_, err = decoder.Token()
	assert.NoError(t, err)

	i := 0
	for decoder.More() {
		// just test 1 object
		if i == 1 {
			break
		}
		key, err := decoder.Token()
		assert.NoError(t, err)

		var port Port
		rawMsg := json.RawMessage{}
		err = decoder.Decode(&rawMsg)
		assert.NoError(t, err)

		err = json.Unmarshal(rawMsg, &port)
		assert.NoError(t, err)

		err = service.portsClient.Store(context.Background(), toProtoPort(
			fmt.Sprintf("%s", key),
			port,
		))
		assert.NoError(t, err)

		i++
	}

	_, err = decoder.Token()
	assert.NoError(t, err)

	err = service.portsClient.Close()
	assert.NoError(t, err)

	mockPortsClient.AssertExpectations(t)
}
