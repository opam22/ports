package adapter

import (
	"context"
	"testing"

	"github.com/opam22/ports/internal/ports/domain/ports"
	"github.com/stretchr/testify/assert"
)

func TestDB_StoreAndGet(t *testing.T) {
	db := NewDB()

	port := &ports.Port{
		PortID: "ZWUTA",
		Name:   "Mutare",
	}

	err := db.Store(context.Background(), port)
	assert.NoError(t, err)

	list, err := db.Get(context.Background())
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, *port, list[0])

	foundPort := db.FindByID(context.Background(), "ZWUTA")
	assert.NotNil(t, foundPort)
	assert.Equal(t, *port, *foundPort)
}

func TestDB_FindByIDNotFound(t *testing.T) {
	db := NewDB()

	port := db.FindByID(context.Background(), "ZWUTA")
	assert.Nil(t, port)
}
