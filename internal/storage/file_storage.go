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

func (fs *FileStorage) ListServers(ctx context.Context) ([]models.Server, error) {
	entries, err := os.ReadDir(fs.StoragePath)
	if err != nil {
		return nil, err
	}

	var servers []models.Server

	for _, fsEntry := range entries {
		filename := fs.StoragePath + "/" + fsEntry.Name()
		content, err := os.ReadFile(filename)
		if err != nil {
			return nil, err
		}

		var server models.Server
		json.Unmarshal(content, &server)
		servers = append(servers, server)
	}

	return servers, nil
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
