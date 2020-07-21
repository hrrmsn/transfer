package predict

import (
	"time"
	cars_ops "wheely/test/internal/cars/client/operations"
	"wheely/test/internal/predict/client"
	predict_ops "wheely/test/internal/predict/client/operations"
	"wheely/test/internal/predict/models"
	"wheely/test/pkg/transfer/utils"

	"github.com/go-openapi/strfmt"
)

type Client struct {
	*client.PredictService

	Formats strfmt.Registry
	Timeout time.Duration
}

func NewClient(cfg *utils.Config) *Client {
	transportCfg := &client.TransportConfig{
		BasePath: cfg.PredictConfig.BasePath,
		Host:     cfg.PredictConfig.Host,
		Schemes:  cfg.PredictConfig.Schemes,
	}
	return &Client{
		PredictService: client.NewHTTPClientWithConfig(strfmt.NewFormats(), transportCfg),
		Formats:        strfmt.NewFormats(),
		Timeout:        cfg.PredictConfig.Timeout,
	}
}

func (c *Client) GetPredict(
	cfg *utils.Config,
	pos *models.Position,
	carsData *cars_ops.GetCarsOK,
) (*predict_ops.PredictOK, error) {

	predictPositions := make([]models.Position, len(carsData.Payload))
	for i := 0; i < len(predictPositions); i++ {
		var car = carsData.Payload[i]
		predictPositions[i] = models.Position{Lat: car.Lat, Lng: car.Lng}
	}

	positionList := predict_ops.PredictBody{
		Target: *pos,
		Source: predictPositions,
	}

	params := &predict_ops.PredictParams{
		PositionList: positionList,
	}
	params.WithTimeout(cfg.PredictConfig.Timeout)

	predictData, err := c.Operations.Predict(params)
	if err != nil {
		return nil, err
	}

	return predictData, nil
}

func (c *Client) Health() (*predict_ops.HealthOK, error) {
	params := &predict_ops.HealthParams{}
	params.WithTimeout(c.Timeout)

	healthData, err := c.Operations.Health(params)
	if err != nil {
		return nil, err
	}

	return healthData, nil
}
