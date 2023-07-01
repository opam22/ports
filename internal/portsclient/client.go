package portsclient

import (
	"context"
	"fmt"
	"log"

	gRPC "github.com/opam22/ports/internal/ports/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Client struct {
	client gRPC.PortServiceClient
	conn   *grpc.ClientConn
}

func NewClient() (*Client, error) {
	connection, err := grpc.Dial(":50001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return &Client{}, fmt.Errorf("error tryng to dial on %v: %w", 50001, err)
	}

	return &Client{
		client: gRPC.NewPortServiceClient(connection),
		conn:   connection,
	}, nil
}

func (c *Client) Store(ctx context.Context, port *gRPC.Port) error {
	log.Println("store client reached")
	_, err := c.client.Store(ctx, &gRPC.StoreRequest{
		Port: port,
	})
	return err
}

func (c *Client) Get(ctx context.Context) ([]*gRPC.Port, error) {
	response, err := c.client.Get(ctx, &emptypb.Empty{})
	return response.GetPorts(), err
}

func (c *Client) Close() error {
	return c.conn.Close()
}
