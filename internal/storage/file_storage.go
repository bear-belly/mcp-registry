package storage

import (
	"context"
	"encoding/json"
	"os"

	"github.com/bear-belly/mcp-registry/internal/models"
)

type FileStorage struct {
	StoragePath string
}

func NewFileStorage(path string) *FileStorage {
	fs := &FileStorage{
		StoragePath: path,
	}

	// TODO: initialise storage path

	return fs
}

func (fs *FileStorage) CreateServer(ctx context.Context, server models.Server) error {
	data, err := json.Marshal(server)
	if err != nil {
		return err
	}

	err = os.WriteFile(server.Name+".json", data, 0644)
	if err != nil {
		return err
	}

	return nil
}
