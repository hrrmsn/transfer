package handlers

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"wheely/test/pkg/transfer/client/cars"
	"wheely/test/pkg/transfer/client/predict"

	"wheely/test/internal/predict/models"
	"wheely/test/pkg/transfer/utils"
)

type Route struct {
	Config        *utils.Config
	CarsClient    *cars.Client
	PredictClient *predict.Client
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

func (r *Route) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	pos, err := readBody(req.Body)
	if err != nil {
		log.Fatalf("Error when reading body from POST request: %s\n", err.Error())
	}

	carsData, err := r.CarsClient.GetCars(r.Config, pos)
	if err != nil {
		log.Fatalf("Error when receiving data from cars client: %s\n", err.Error())
	}

	predictData, err := r.PredictClient.GetPredict(r.Config, pos, carsData)
	if err != nil {
		log.Fatalf("Error when receiving data from predict client: %s\n", err.Error())
	}

	minTime, err := utils.Min(predictData.Payload)
	if err != nil {
		log.Fatalf("Error when handle predict data: %s\n", err.Error())
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"response": "` + strconv.Itoa(int(minTime)) + `"}`))
}
