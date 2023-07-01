package importer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/opam22/ports/internal/portsclient"
)

type Service struct {
	portsClient portsclient.Client
}

func NewService() (*Service, error) {
	client, err := portsclient.NewClient()
	if err != nil {
		return &Service{}, fmt.Errorf("faiiled to initializing importer service %w", err)
	}

	return &Service{
		portsClient: *client,
	}, nil
}

func (s *Service) Run(ctx context.Context) (err error) {
	file, err := os.Open("internal/importer/json/ports.json")
	if err != nil {
		fmt.Println("error opening file:", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	_, err = decoder.Token()
	if err != nil {
		fmt.Println("error decoding JSON:", err)
		return
	}

	for decoder.More() {
		key, err := decoder.Token()
		if err != nil {
			fmt.Println("error decoding JSON:", err)
			continue
		}

		port := Port{}
		rawMsg := json.RawMessage{}
		err = decoder.Decode(&rawMsg)
		if err != nil {
			fmt.Println("error decoding JSON:", err)
		}

		err = json.Unmarshal(rawMsg, &port)
		if err != nil {
			fmt.Println("error decoding JSON:", err)
		}

		if err := s.portsClient.Store(ctx, toProtoPort(
			fmt.Sprintf("%s", key),
			port,
		)); err != nil {

		}
	}

	_, err = decoder.Token()
	if err != nil {
		fmt.Println("error decoding JSON:", err)
		return
	}

	if err := s.portsClient.Close(); err != nil {
		err = fmt.Errorf("failed to close client connection: %w", err)
	}

	return nil
}
