package storage

import (
	"context"

	"github.com/bear-belly/mcp-registry/internal/models"
)

type Storage interface {
	CreateServer(ctx context.Context, server models.Server) error
}
