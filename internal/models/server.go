package models

type Server struct {
	Name      string `json:"name"`
	Transport string `json:"transport"`
	URL       string `json:"url"`
}
