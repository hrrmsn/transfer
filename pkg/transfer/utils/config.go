package utils

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port         string        `envconfig:"PORT" default:"8081"`
	ReadTimout   time.Duration `envconfig:"READ_TIMEOUT" default:"10s"`
	WriteTimeout time.Duration `envconfig:"READ_TIMEOUT" default:"10s"`

	CarsConfig
	PredictConfig
}

type CarsConfig struct {
	BasePath string        `envconfig:"CARS_BASE_PATH" default:"/path"`
	Host     string        `envconfig:"CARS_HOST" default:"example.com"`
	Schemes  []string      `envconfig:"CARS_SCHEMES" default:"https"`
	Timeout  time.Duration `envconfig:"CARS_TIMEOUT" default:"30s"`
	Limit    int           `envconfig:"CARS_LIMIT" default:"10"`
}

type PredictConfig struct {
	BasePath string        `envconfig:"PREDICT_BASE_PATH" default:"/path"`
	Host     string        `envconfig:"PREDICT_HOST" default:"example.com"`
	Schemes  []string      `envconfig:"PREDICT_SCHEMES" default:"https"`
	Timeout  time.Duration `envconfig:"PREDICT_TIMEOUT" default:"30s"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("config", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
