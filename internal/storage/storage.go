package storage

import (
	"context"

	"github.com/bear-belly/mcp-registry/internal/models"
)

type Storage interface {
	ListServers(ctx context.Context) ([]models.Server, error)
	CreateServer(ctx context.Context, server models.Server) error
}
