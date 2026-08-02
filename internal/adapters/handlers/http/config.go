package http

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type ServerConfig struct {
	Port int `env:"PORT" envDefault:"8080"`
}

func LoadServerConfig() (*ServerConfig, error) {
	cfg := &ServerConfig{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to load http server config: %w", err)
	}
	if cfg.Port < 1 || cfg.Port > 65535 {
		return nil, fmt.Errorf("invalid PORT %d: must be between 1 and 65535", cfg.Port)
	}
	return cfg, nil
}
