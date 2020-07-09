package delivery

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-openapi/strfmt"

	cars "wheely/test/internal/cars/client"
	carsOperations "wheely/test/internal/cars/client/operations"
	predict "wheely/test/internal/predict/client"
	predictOperations "wheely/test/internal/predict/client/operations"
	predictModels "wheely/test/internal/predict/models"
)

const (
	carsLimit = 10
	timeout   = 2 * time.Second
)

func getCurrentPos(w http.ResponseWriter, r *http.Request) (*predictModels.Position, error) {
	defer r.Body.Close()

	resp, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "can't read from request body"}`))
		return nil, err
	}

	pos := &predictModels.Position{}
	if err = pos.UnmarshalBinary(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "can't unmarshal request body"}`))
		return nil, err
	}

	return pos, nil
}

// call cars API
func getCars(pos *predictModels.Position) (*carsOperations.GetCarsOK, error) {
	carsTransportCfg := &cars.TransportConfig{
		BasePath: "/fake-eta",
		Host:     "dev-api.wheely.com",
		Schemes:  []string{"https"},
	}

	formats := strfmt.NewFormats()
	carsClient := cars.NewHTTPClientWithConfig(formats, carsTransportCfg)

	params := &carsOperations.GetCarsParams{
		Lat:   pos.Lat,
		Lng:   pos.Lng,
		Limit: carsLimit,
	}
	params.WithTimeout(timeout)

	getCarsOK, err := carsClient.Operations.GetCars(params)
	if err != nil {
		return nil, err
	}

	return getCarsOK, nil
}

func getPredicts(pos *predictModels.Position, carsData *carsOperations.GetCarsOK) (
	*predictOperations.PredictOK, error) {

	formats := strfmt.NewFormats()
	transportCfg := &predict.TransportConfig{
		BasePath: "/fake-eta",
		Host:     "dev-api.wheely.com",
		Schemes:  []string{"https"},
	}
	predictClient := predict.NewHTTPClientWithConfig(formats, transportCfg)

	predictPositions := make([]predictModels.Position, 0)
	for _, car := range carsData.Payload {
		predictPositions = append(predictPositions, predictModels.Position{Lat: car.Lat, Lng: car.Lng})
	}
	// DEBUG
	fmt.Printf("predict positions: %#v\n\n", predictPositions)

	positionList := predictOperations.PredictBody{
		Target: *pos,
		Source: predictPositions,
	}

	params := &predictOperations.PredictParams{
		PositionList: positionList,
	}
	params.WithTimeout(timeout)

	predictData, err := predictClient.Operations.Predict(params)
	if err != nil {
		return nil, err
	}

	return predictData, nil
}

func getMinTimeDelivery(times []int64) (int, error) {
	if len(times) == 0 {
		return -1, fmt.Errorf("No predict data")
	}

	var minTime = times[0]
	for i := 1; i < len(times); i++ {
		if times[i] < minTime {
			minTime = times[i]
		}
	}
	return int(minTime), nil
}

func findNearestRoute(w http.ResponseWriter, r *http.Request) {
	pos, err := getCurrentPos(w, r)
	if err != nil {
		log.Fatalf("Error when reading body from POST request: %s\n", err.Error())
	}

	carsData, err := getCars(pos)
	if err != nil {
		log.Fatalf("Error when receiving data from cars client: %s\n", err.Error())
	}
	// DEBUG
	fmt.Printf("%#v\n\n", carsData.Payload)

	predictData, err := getPredicts(pos, carsData)
	if err != nil {
		log.Fatalf("Error when receiving data from predict client: %s\n", err.Error())
	}
	// DEBUG
	fmt.Printf("%#v\n", predictData.Payload)

	minTime, err := getMinTimeDelivery(predictData.Payload)
	if err != nil {
		log.Fatalf("Error when handle predict data: %s\n", err.Error())
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"response": "` + strconv.Itoa(minTime) + `"}`))
}

func NewServer(port string) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/delivery", findNearestRoute)

	return &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
