package predict

import (
	"strings"
	"time"

	"github.com/go-openapi/strfmt"

	cars_ops "wheely/test/internal/cars/client/operations"
	"wheely/test/internal/predict/client"
	predict_ops "wheely/test/internal/predict/client/operations"
	"wheely/test/internal/predict/models"
	"wheely/test/pkg/transfer/utils"
)

// type Predictor interface {
// 	GetPredict(*models.Position, *cars_ops.GetCarsOK) (*predict_ops.PredictOK, error)
// 	Healthy() bool
// }

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

func (c *Client) GetPredict(pos *models.Position, carsData *cars_ops.GetCarsOK) (*predict_ops.PredictOK, error) {
	predictPositions, err := carsToPositions(carsData, c.Formats)
	if err != nil {
		return nil, err
	}

	positionList := predict_ops.PredictBody{
		Target: *pos,
		Source: predictPositions,
	}

	params := &predict_ops.PredictParams{
		PositionList: positionList,
	}
	params.WithTimeout(c.Timeout)

	predictData, err := c.Operations.Predict(params)
	if err != nil {
		return nil, err
	}

	return predictData, nil
}

func (c *Client) health() (*predict_ops.HealthOK, error) {
	params := &predict_ops.HealthParams{}
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
