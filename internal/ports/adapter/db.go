package adapter

import (
	"context"
	"sync"

	"github.com/opam22/ports/internal/ports/domain/ports"
)

// DB struct using sync.Map to store the ports data using map
// it is conccurent safe so no need to use Mu.Lock Mu.Unlock
type DB struct {
	ports sync.Map
}

func NewDB() *DB {
	return &DB{}
}

// storing data to db
func (db *DB) Store(ctx context.Context, port *ports.Port) error {
	db.ports.Store(port.PortID, port)
	return nil
}

// get data from db
func (db *DB) Get(ctx context.Context) ([]ports.Port, error) {
	var list []ports.Port

	db.ports.Range(func(_, value interface{}) bool {
		if port, ok := value.(*ports.Port); ok {
			list = append(list, *port)
		}
		return true
	})

	return list, nil
}

// find by id
func (db *DB) FindByID(ctx context.Context, portID string) *ports.Port {
	var (
		foundPort *ports.Port
		found     bool
	)

	db.ports.Range(func(key, value interface{}) bool {
		if keyStr, ok := key.(string); ok && keyStr == portID {
			if port, ok := value.(*ports.Port); ok {
				foundPort = port
				found = true
				return false
			}
		}
		return true
	})

	if found {
		return foundPort
	}

	return nil
}
