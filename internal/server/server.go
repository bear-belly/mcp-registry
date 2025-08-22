package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/bear-belly/mcp-registry/internal/errors"
	"github.com/bear-belly/mcp-registry/internal/middleware"
	"github.com/bear-belly/mcp-registry/internal/models"
	"github.com/bear-belly/mcp-registry/internal/storage"
	"github.com/bear-belly/mcp-registry/internal/templates"
)

type Server struct {
	config        models.Config
	storage       storage.Storage
	mux           *http.ServeMux
	startTime     time.Time
	healthyStatus *bool
}

type Metrics struct {
	Uptime float64 `json:"uptime_seconds"`
}

func New(storage storage.Storage, config models.Config) *Server {
	healthyStatus := true

	return &Server{
		config:        config,
		storage:       storage,
		mux:           http.NewServeMux(),
		startTime:     time.Now(),
		healthyStatus: &healthyStatus,
	}
}

func (s *Server) setupStaticRoutes() {
	// Serve static files
	fs := http.FileServer(http.Dir(filepath.Join("internal", "templates", "static")))
	s.mux.Handle("/static/", http.StripPrefix("/static/", fs))
}

func (s *Server) setupHealthRoutes() {
	// Health check endpoint
	s.mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if !*s.healthyStatus {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "healthy")
	})

	// Uptime endpoint
	s.mux.HandleFunc("/uptime", func(w http.ResponseWriter, r *http.Request) {
		uptime := time.Since(s.startTime).String()
		fmt.Fprintf(w, "Uptime: %s", uptime)
	})
}

func (s *Server) getServerByName(name string) (models.Server, error) {
	ctx := context.Background()
	servers, err := s.storage.ListServers(ctx)
	if err != nil {
		return models.Server{}, err
	}

	for _, server := range servers {
		if server.Name == name {
			return server, nil
		}
	}

	return models.Server{}, errors.NewNotFoundError("Server")
}

func (s *Server) setupHomeRoute() {
	// Home page route
	s.mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only handle exact root path, not other paths
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		ctx := r.Context()

		// Retrieve servers from storage
		servers, err := s.storage.ListServers(ctx)
		if err != nil {
			errors.WriteError(w, errors.NewInternalError("Error retrieving servers", err))
			return
		}

		// Map server data to template data
		data := templates.PageData{
			Title:        "MCP Registry",
			PageTemplate: "index",
			Data:         servers,
		}

		// Render the template
		if err := templates.ExecuteTemplate(ctx, w, "layout.html", data); err != nil {
			errors.WriteError(w, errors.NewInternalError("Error rendering template", err))
		}
	}))

	// Server details route
	s.mux.Handle("/server/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverName := strings.TrimPrefix(r.URL.Path, "/server/")
		if serverName == "" {
			http.NotFound(w, r)
			return
		}

		ctx := r.Context()
		server, err := s.getServerByName(serverName)
		if err != nil {
			if _, ok := err.(*errors.AppError); ok {
				errors.WriteError(w, err)
			} else {
				errors.WriteError(w, errors.NewInternalError("Error retrieving server", err))
			}
			return
		}

		// Convert config to JSON string if it exists
		var configJSON string
		if server.Config != nil {
			configBytes, err := json.MarshalIndent(server.Config, "", "    ")
			if err != nil {
				errors.WriteError(w, errors.NewInternalError("Error formatting config", err))
				return
			}
			configJSON = string(configBytes)
		}

		// Map server data to template data
		data := templates.PageData{
			Title:        server.Name + " - MCP Registry",
			PageTemplate: "server",
			Data:         server,
			ConfigJSON:   configJSON,
		}

		// Render the template
		if err := templates.ExecuteTemplate(ctx, w, "layout.html", data); err != nil {
			errors.WriteError(w, errors.NewInternalError("Error rendering template", err))
		}
	}))
}

func (s *Server) timingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func (s *Server) recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("Panic recovered in %s: %v\n", r.URL.Path, err)
				errors.WriteError(w, errors.NewInternalError("Internal Server Error", fmt.Errorf("%v", err)))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (s *Server) setupApiRoutes() {
	s.mux.Handle("/api/servers/v1", middleware.CorsMiddleware(
		http.HandlerFunc(s.ListServersV1)))

	// TODO: add further CRUD operation endpoints
}

// ListJourneys handles retrieving a list of journeys
func (s *Server) ListServersV1(w http.ResponseWriter, r *http.Request) {
	servers, err := s.storage.ListServers(r.Context())
	if err != nil {
		errors.WriteError(w, errors.NewDatabaseError("Failed to retrieve journeys", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(servers)
}

func (s *Server) SetupRoutes() {
	s.setupStaticRoutes()
	s.setupHealthRoutes()
	s.setupApiRoutes()
	s.setupHomeRoute()
}

func (s *Server) Handler() http.Handler {
	return s.recoveryMiddleware(s.timingMiddleware(s.mux))
}

func (s *Server) SetHealthStatus(healthy bool) {
	*s.healthyStatus = healthy
}
