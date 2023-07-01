package ports

import (
	"context"
	"fmt"
	"net"

	"github.com/opam22/ports/internal/ports/adapter"
	"github.com/opam22/ports/internal/ports/domain/ports"
	gRPC "github.com/opam22/ports/internal/ports/grpc"
	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	port = ":50001"
)

type Server struct {
	logger *logrus.Logger
	gRPC.UnimplementedPortServiceServer
	service ports.Service
}

func NewServer(logger *logrus.Logger) *Server {
	return &Server{
		logger: logger,
		service: ports.NewService(
			logger,
			adapter.NewDB(),
		),
	}
}

func (s *Server) Serve(ctx context.Context) error {
	grpcServer := grpc.NewServer()
	gRPC.RegisterPortServiceServer(grpcServer, s)
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("server fail to listening: %w", err)
	}

	s.logger.Info("server running on ", port)
	err = grpcServer.Serve(listener)
	if err != nil {
		return fmt.Errorf("server failed to serve: %w", err)
	}

	return nil
}

func (s *Server) Store(ctx context.Context, req *gRPC.StoreRequest) (*emptypb.Empty, error) {
	if err := s.service.Store(ctx, toPortDomain(req.Port)); err != nil {
		s.logger.Errorf("error when storing port %+v", err)
		return &emptypb.Empty{}, fmt.Errorf("fail to store port")
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) Get(ctx context.Context, _ *emptypb.Empty) (*gRPC.GetResponse, error) {
	ports, err := s.service.Get(ctx)
	if err != nil {
		return nil, err
	}

	return &gRPC.GetResponse{
		Ports: toProtoPort(ports),
	}, nil
}

func toProtoPort(ports []ports.Port) []*gRPC.Port {
	proto := []*gRPC.Port{}
	for _, p := range ports {
		port := &gRPC.Port{
			PortId:      p.PortID,
			Name:        p.Name,
			City:        p.City,
			Country:     p.Country,
			Alias:       p.Alias,
			Regions:     p.Regions,
			Coordinates: p.Coordinates,
			Province:    p.Province,
			Timezone:    p.Timezone,
			Unlocs:      p.Unlocs,
			Code:        p.Code,
		}

		proto = append(proto, port)
	}

	return proto
}

func toPortDomain(p *gRPC.Port) *ports.Port {
	return &ports.Port{
		PortID:      p.PortId,
		Name:        p.Name,
		City:        p.City,
		Country:     p.Country,
		Alias:       p.Alias,
		Regions:     p.Regions,
		Coordinates: p.Coordinates,
		Province:    p.Province,
		Timezone:    p.Timezone,
		Unlocs:      p.Unlocs,
		Code:        p.Code,
	}
}
