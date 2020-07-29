package cars

import (
	"strings"
	"time"

	"github.com/go-openapi/strfmt"

	"wheely/test/internal/cars/client"
	cars_ops "wheely/test/internal/cars/client/operations"
	"wheely/test/internal/predict/models"
	"wheely/test/pkg/transfer/utils"
)

// type Finder interface {
// 	GetCars(*models.Position) (*cars_ops.GetCarsOK, error)
// 	Validate(*cars_ops.GetCarsOK) error
// 	Healthy() bool
// }

type Client struct {
	*client.CarsService

	Formats   strfmt.Registry
	Timeout   time.Duration
	CarsLimit int
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
		CarsLimit:   cfg.CarsConfig.Limit,
	}
}

func (c *Client) GetCars(pos *models.Position) (*cars_ops.GetCarsOK, error) {
	params := &cars_ops.GetCarsParams{
		Lat:   pos.Lat,
		Lng:   pos.Lng,
		Limit: int64(c.CarsLimit),
	}
	params.WithTimeout(c.Timeout)

	getCarsOK, err := c.Operations.GetCars(params)
	if err != nil {
		return nil, err
	}

	return getCarsOK, nil
}

func (c *Client) Validate(carsData *cars_ops.GetCarsOK) error {
	var err error
	for _, car := range carsData.Payload {
		if err = car.Validate(c.Formats); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) health() (*cars_ops.HealthOK, error) {
	params := &cars_ops.HealthParams{}
	params.WithTimeout(c.Timeout)

	healthData, err := c.Operations.Health(params)
	if err != nil {
		return nil, err
	}

	return healthData, nil
}

func (c *Client) Healthy() bool {
	healthData, err := c.health()

	if err != nil && !strings.Contains(err.Error(), "unknown error") {
		return false
	}

	return strings.Contains(healthData.Error(), "200") || strings.Contains(healthData.Error(), "healthOK")
}
