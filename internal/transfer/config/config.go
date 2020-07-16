package config

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
	BasePath string        `envconfig:"CARS_BASE_PATH" default:"/fake-eta"`
	Host     string        `envconfig:"CARS_HOST" default:"dev-api.wheely.com"`
	Schemes  []string      `envconfig:"CARS_SCHEMES" default:"https"`
	Timeout  time.Duration `envconfig:"CARS_TIMEOUT" default:"30s"`
	Limit    int           `envconfig:"CARS_LIMIT" default:10`
}

type PredictConfig struct {
	BasePath string        `envconfig:"PREDICT_BASE_PATH" default:"/fake-eta"`
	Host     string        `envconfig:"PREDICT_HOST" default:"dev-api.wheely.com"`
	Schemes  []string      `envconfig:"PREDICT_SCHEMES" default:"https"`
	Timeout  time.Duration `envconfig:"PREDICT_TIMEOUT" default:"30s"`
}

func New() (*Config, error) {
	// var carsCfg CarsConfig
	// if err := envconfig.Process("cars config", &carsCfg); err != nil {
	// 	return nil, err
	// }
	//
	// var predictCfg PredictConfig
	// if err := envconfig.Process("cars config", &predictCfg); err != nil {
	// 	return nil, err
	// }

	var cfg Config
	if err := envconfig.Process("config", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
