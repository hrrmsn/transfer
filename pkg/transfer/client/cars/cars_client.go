package cars

import (
	"wheely/test/internal/cars/client"
	"wheely/test/internal/cars/client/operations"
	"wheely/test/internal/predict/models"
	"wheely/test/pkg/transfer/utils"

	"github.com/go-openapi/strfmt"
)

type Client struct {
	*client.CarsService
}

func NewClient(cfg *utils.Config) *Client {
	transportCfg := &client.TransportConfig{
		BasePath: cfg.CarsConfig.BasePath,
		Host:     cfg.CarsConfig.Host,
		Schemes:  cfg.CarsConfig.Schemes,
	}
	return &Client{
		CarsService: client.NewHTTPClientWithConfig(strfmt.NewFormats(), transportCfg),
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
