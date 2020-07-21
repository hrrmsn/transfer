package cars

import (
	"strings"
	"time"
	"wheely/test/internal/cars/client"
	"wheely/test/internal/cars/client/operations"
	"wheely/test/internal/predict/models"
	"wheely/test/pkg/transfer/utils"

	"github.com/go-openapi/strfmt"
)

type Client struct {
	*client.CarsService

	Formats strfmt.Registry
	Timeout time.Duration
}

func NewClient(cfg *utils.Config) *Client {
	transportCfg := &client.TransportConfig{
		BasePath: cfg.CarsConfig.BasePath,
		Host:     cfg.CarsConfig.Host,
		Schemes:  cfg.CarsConfig.Schemes,
	}
	return &Client{
		CarsService: client.NewHTTPClientWithConfig(strfmt.NewFormats(), transportCfg),
		Formats:     strfmt.NewFormats(),
		Timeout:     cfg.CarsConfig.Timeout,
	}
}

func (c *Client) GetCars(cfg *utils.Config, pos *models.Position) (*operations.GetCarsOK, error) {
	params := &operations.GetCarsParams{
		Lat:   pos.Lat,
		Lng:   pos.Lng,
		Limit: int64(cfg.CarsConfig.Limit),
	}
	params.WithTimeout(cfg.CarsConfig.Timeout)

	getCarsOK, err := c.Operations.GetCars(params)
	if err != nil {
		return nil, err
	}

	return getCarsOK, nil
}

func (c *Client) Validate(carsData *operations.GetCarsOK) error {
	var err error
	for _, car := range carsData.Payload {
		if err = car.Validate(c.Formats); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) health() (*operations.HealthOK, error) {
	params := &operations.HealthParams{}
	params.WithTimeout(c.Timeout)

	healthData, err := c.Operations.Health(params)
	if err != nil {
		return nil, err
	}

	return healthData, nil
}

func (c *Client) Healthy() bool {
	healthData, err := c.health()
	if err != nil {
		return false
	}
	return strings.Contains(healthData.Error(), "200") || strings.Contains(healthData.Error(), "healthOK")
}

func (c *Client) Unhealthy() bool {
	return !c.Healthy()
}
