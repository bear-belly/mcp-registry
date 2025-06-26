package main

import (
	"fmt"

	"github.com/bear-belly/mcp-registry/internal/models"
)

func main() {
	// principle of dependency injection - create what i need here and pass it to the dependent object
	// rather than expect them to create what they need

	// everything starts with configuration
	config := models.Config{
		StorageType: "file",
		StoragePath: "./data",
		LogLevel: "INFO",
	}

	// initialise a global logger, based on slog but abstracted to change easily later
	NewLogger(config)

	// create a storage using the factory pattern
	storage, err := NewStorage(config)
	if err != nil {
		logger.Error("Could not start due to error in the storage subsystem", err)
		return
	}
}
