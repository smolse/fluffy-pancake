package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// ServerConfig represents the configuration of the HTTP server.
type ServerConfig struct {
	Port                    int           `default:"8080"`
	GracefulShutdownTimeout time.Duration `default:"5s"`
}

// DataStoreConfig represents the configuration of the data store.
type DataStoreConfig struct {
	Type string `default:"syncmap"`
}

// GinConfig represents the configuration of the Gin framework.
type GinConfig struct {
	Mode string `default:"release"`
}

// Config represents the configuration of the application.
type Config struct {
	Server    ServerConfig
	DataStore DataStoreConfig
	Gin       GinConfig
}

// LoadConfig loads the whole application configuration from the environment.
func LoadConfig() (Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	return cfg, err
}
