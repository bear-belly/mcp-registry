package main

import (
	"net/http"

	"github.com/bear-belly/mcp-registry/internal/logger"
	"github.com/bear-belly/mcp-registry/internal/models"
	"github.com/bear-belly/mcp-registry/internal/server"
	"github.com/bear-belly/mcp-registry/internal/storage"
	"github.com/bear-belly/mcp-registry/internal/templates"
)

func main() {
	// principle of dependency injection - create what i need here and pass it to the dependent object
	// rather than expect them to create what they need

	// everything starts with configuration
	config := models.Config{
		StorageType:  "file",
		StoragePath:  "./data",
		TemplatePath: "./internal/templates",
		LogLevel:     "INFO",
	}

	// initialise a global logger, based on slog but abstracted to change easily later
	logger.NewLogger(config)

	// create a storage using the factory pattern
	logger.Info("Configuring storage...")
	storage, err := storage.NewStorage(config)
	if err != nil {
		logger.Error("Could not start due to error in the storage subsystem", err)
		return
	}

	// initialise page templates
	logger.Info("Configuring templater...")
	err = templates.InitTemplates(config)
	if err != nil {
		logger.Error("Could not configure templater", err)
		return
	}

	// create and configure HTTP server
	server := server.New(storage, config)
	server.SetupRoutes()

	logger.Info("Starting server on :8080")
	if err := http.ListenAndServe(":8080", server.Handler()); err != nil {
		logger.Error("Server failed to start: ", err)
	}
}
