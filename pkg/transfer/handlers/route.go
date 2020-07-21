package handlers

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"wheely/test/pkg/transfer/clients/cars"
	"wheely/test/pkg/transfer/clients/predict"

	"wheely/test/internal/predict/models"
	"wheely/test/pkg/transfer/utils"

	"github.com/go-openapi/strfmt"
)

type Route struct {
	Config        *utils.Config
	CarsClient    *cars.Client
	PredictClient *predict.Client
	Formats       strfmt.Registry
}

func NewRoute(cfg *utils.Config) *Route {
	return &Route{
		Config:        cfg,
		CarsClient:    cars.NewClient(cfg),
		PredictClient: predict.NewClient(cfg),
		Formats:       strfmt.NewFormats(),
	}
}

func readBody(body io.ReadCloser) (*models.Position, error) {
	defer body.Close()

	resp, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	pos := &models.Position{}
	if err = pos.UnmarshalBinary(resp); err != nil {
		return nil, err
	}

	return pos, nil
}

func readPos(request *http.Request, formats strfmt.Registry) (*models.Position, error) {
	pos, err := readBody(request.Body)
	if err != nil {
		return nil, fmt.Errorf("Error when reading body of POST request: %s", err.Error())
	}

	if err = pos.Validate(formats); err != nil {
		return nil, fmt.Errorf("Predict position validate error: %s", err.Error())
	}

	return pos, nil
}

func (r *Route) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	pos, err := readPos(req, r.Formats)
	if err != nil {
		utils.HandleError("Input data error", http.StatusInternalServerError, err, w)
		return
	}

	if r.CarsClient.Unhealthy() {
		// take data from cache
	}

	carsData, err := r.CarsClient.GetCars(r.Config, pos)
	if err != nil {
		utils.HandleError("Error when receiving data from cars service", http.StatusInternalServerError, err, w)
		return
	}

	predictData, err := r.PredictClient.GetPredict(r.Config, pos, carsData)
	if err != nil {
		utils.HandleError("Error when receiving data from predict service", http.StatusInternalServerError, err, w)
		return
	}

	minTime, err := utils.Min(predictData.Payload)
	if err != nil && err.Error() == "Empty input array" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"response": "no cars available"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"response": "` + strconv.Itoa(int(minTime)) + `"}`))
}
