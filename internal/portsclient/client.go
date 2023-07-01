package portsclient

import (
	"context"
	"fmt"

	gRPC "github.com/opam22/ports/internal/ports/grpc"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Client struct {
	logger *logrus.Logger
	client gRPC.PortServiceClient
	conn   *grpc.ClientConn
}

func NewClient(logger *logrus.Logger) (*Client, error) {
	connection, err := grpc.Dial(":50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &Client{}, fmt.Errorf("error tryng to dial on %v: %w", 50001, err)
	}

	return &Client{
		logger: logger,
		client: gRPC.NewPortServiceClient(connection),
		conn:   connection,
	}, nil
}

func (c *Client) Store(ctx context.Context, port *gRPC.Port) error {
	c.logger.Info("postclient Store called")
	_, err := c.client.Store(ctx, &gRPC.StoreRequest{
		Port: port,
	})
	return err
}

func (c *Client) Get(ctx context.Context) ([]*gRPC.Port, error) {
	c.logger.Info("postclient Get called")
	response, err := c.client.Get(ctx, &emptypb.Empty{})
	return response.GetPorts(), err
}

func (c *Client) Close() error {
	return c.conn.Close()
}
