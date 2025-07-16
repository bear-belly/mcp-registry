package models

type Config struct {
	StorageType  string `json:"storage_type"`
	StoragePath  string `json:"storage_path"`
	TemplatePath string `json:"template_path"`
	LogLevel     string `json:"log_level"`
}
