package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	HTTPPort string `env:"HTTP_PORT" envDefault:"8080"`

	PostgresDSN string `env:"POSTGRES_DSN,required"`

	KafkaBrokers string `env:"KAFKA_BROKERS" envDefault:"localhost:9092"`
	KafkaTopic   string `env:"KAFKA_TOPIC" envDefault:"events"`
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("error parsing config:", err)
	}

	return cfg, nil
}
