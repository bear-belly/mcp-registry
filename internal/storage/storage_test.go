package storage

import (
	"context"
	"errors"
	"testing"

	"github.com/bear-belly/mcp-registry/internal/models"
)

// MockStorage implements the Storage interface for testing.
type MockStorage struct {
	Servers []models.Server
	Err     error
}

func (m *MockStorage) ListServers(ctx context.Context) ([]models.Server, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Servers, nil
}

func (m *MockStorage) CreateServer(ctx context.Context, server models.Server) error {
	if m.Err != nil {
		return m.Err
	}
	m.Servers = append(m.Servers, server)
	return nil
}

func TestListServers_Success(t *testing.T) {
	mock := &MockStorage{
		Servers: []models.Server{
			{Name: "Server1"},
			{Name: "Server2"},
		},
	}
	servers, err := mock.ListServers(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(servers) != 2 {
		t.Errorf("expected 2 servers, got %d", len(servers))
	}
}

func TestListServers_Error(t *testing.T) {
	mock := &MockStorage{
		Err: errors.New("db error"),
	}
	_, err := mock.ListServers(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCreateServer_Success(t *testing.T) {
	mock := &MockStorage{}
	server := models.Server{Name: "Server3"}
	err := mock.CreateServer(context.Background(), server)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(mock.Servers) != 1 {
		t.Errorf("expected 1 server, got %d", len(mock.Servers))
	}
}

func TestCreateServer_Error(t *testing.T) {
	mock := &MockStorage{
		Err: errors.New("insert error"),
	}
	server := models.Server{Name: "Server4"}
	err := mock.CreateServer(context.Background(), server)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
