package storage

import (
	"fmt"

	"github.com/bear-belly/mcp-registry/internal/models"
)

func NewStorage(config models.Config) (Storage, error) {
	switch config.StorageType {
	case "file":
		return NewFileStorage(config.StoragePath), nil
	case "psql":
		return nil, fmt.Errorf("Not implemented yet")
	}

	return nil, fmt.Errorf("Unknown storage subsystem")
}
