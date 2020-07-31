package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-openapi/strfmt"

	cars_ops "wheely/test/internal/cars/client/operations"
	predict_ops "wheely/test/internal/predict/client/operations"
	"wheely/test/internal/predict/models"
	"wheely/test/pkg/transfer/clients/cars"
	"wheely/test/pkg/transfer/clients/predict"
	"wheely/test/pkg/transfer/utils"
)

type RouteHandler struct {
	Config *utils.Config

	CarsClient *cars.Client
	// CarsClient cars.Finder

	PredictClient *predict.Client
	// PredictClient predict.Predictor

	Formats strfmt.Registry
}

func NewRouteHandler(cfg *utils.Config) *RouteHandler {
	return &RouteHandler{
		Config:        cfg,
		CarsClient:    cars.NewClient(cfg),
		PredictClient: predict.NewClient(cfg),
		Formats:       strfmt.NewFormats(),
	}
}

func (rh *RouteHandler) GetCars(pos *models.Position) (*cars_ops.GetCarsOK, error) {
	if !rh.CarsClient.Healthy() {
		return nil, fmt.Errorf("Cars service is unavailable")
	}

	carsData, err := rh.CarsClient.GetCars(pos)
	if err != nil {
		return nil, utils.WrapError("Error when receiving data from cars service", err)
	}

	if err = rh.CarsClient.Validate(carsData); err != nil {
		return nil, utils.WrapError("Cars data is invalid", err)
	}

	return carsData, nil
}

func (rh *RouteHandler) GetPredict(
	pos *models.Position,
	carsData *cars_ops.GetCarsOK,
) (*predict_ops.PredictOK, error) {

	if !rh.PredictClient.Healthy() {
		return nil, fmt.Errorf("Predict service is unavailable")
	}

	predictData, err := rh.PredictClient.GetPredict(pos, carsData)
	if err != nil {
		return nil, utils.WrapError("Error when receiving data from predict service", err)
	}

	return predictData, nil
}

func (rh *RouteHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	pos, err := readPos(req, rh.Formats)
	if err != nil {
		utils.HandleError(w, utils.WrapError("Input data error", err), http.StatusInternalServerError)
		return
	}

	carsData, err := rh.GetCars(pos)
	if err != nil {
		utils.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	predictData, err := rh.GetPredict(pos, carsData)
	if err != nil {
		utils.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	minTime, err := utils.Min(predictData.Payload)
	if err != nil && err.Error() == "Empty input array" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"response": "no cars available"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"response": ` + strconv.Itoa(int(minTime)) + `}`))
}
