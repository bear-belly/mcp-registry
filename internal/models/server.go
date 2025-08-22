package models

import "time"

type Server struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Transport   string                 `json:"transport"`
	Status      string                 `json:"status"`
	CreatedAt   time.Time              `json:"createdAt"`
	URL         string                 `json:"url"`
	Config      map[string]interface{} `json:"config,omitempty"`
}
