package transfer

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/go-openapi/strfmt"

	cars "wheely/test/internal/cars/client"
	carsop "wheely/test/internal/cars/client/operations"
	predict "wheely/test/internal/predict/client"
	predictop "wheely/test/internal/predict/client/operations"
	pmodels "wheely/test/internal/predict/models"

	config "wheely/test/internal/transfer/config"
)

type Transfer interface {
	Predict() int64
}

type TransferServer struct {
	*http.Server

	CarsClient    *cars.CarsService
	PredictClient *predict.PredictService
}

func initHTTPServer(mux *http.ServeMux, cfg *config.Config) *http.Server {
	return &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  cfg.ReadTimout,
		WriteTimeout: cfg.WriteTimeout,
	}
}

func initCarsClient(cfg *config.Config) *cars.CarsService {
	transportCfg := &cars.TransportConfig{
		BasePath: cfg.CarsConfig.BasePath,
		Host:     cfg.CarsConfig.Host,
		Schemes:  cfg.CarsConfig.Schemes,
	}
	return cars.NewHTTPClientWithConfig(strfmt.NewFormats(), transportCfg)
}

func initPredictClient(cfg *config.Config) *predict.PredictService {
	transportCfg := &predict.TransportConfig{
		BasePath: cfg.PredictConfig.BasePath,
		Host:     cfg.PredictConfig.Host,
		Schemes:  cfg.PredictConfig.Schemes,
	}
	return predict.NewHTTPClientWithConfig(strfmt.NewFormats(), transportCfg)
}

type NearestRoute struct {
	Config        *config.Config
	CarsClient    *cars.CarsService
	PredictClient *predict.PredictService
}

func (nr *NearestRoute) getCurrentPos(r io.ReadCloser) (*pmodels.Position, error) {
	defer r.Close()

	resp, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	pos := &pmodels.Position{}
	if err = pos.UnmarshalBinary(resp); err != nil {
		return nil, err
	}

	return pos, nil
}

func (nr *NearestRoute) getCars(pos *pmodels.Position) (*carsop.GetCarsOK, error) {
	params := &carsop.GetCarsParams{
		Lat:   pos.Lat,
		Lng:   pos.Lng,
		Limit: int64(nr.Config.CarsConfig.Limit),
	}
	params.WithTimeout(nr.Config.CarsConfig.Timeout)

	getCarsOK, err := nr.CarsClient.Operations.GetCars(params)
	if err != nil {
		return nil, err
	}

	return getCarsOK, nil
}

func (nr *NearestRoute) getPredicts(pos *pmodels.Position, carsData *carsop.GetCarsOK) (*predictop.PredictOK, error) {
	predictPositions := make([]pmodels.Position, 0)
	for _, car := range carsData.Payload {
		predictPositions = append(predictPositions, pmodels.Position{Lat: car.Lat, Lng: car.Lng})
	}
	// DEBUG
	fmt.Printf("predict positions: %#v\n\n", predictPositions)

	positionList := predictop.PredictBody{
		Target: *pos,
		Source: predictPositions,
	}

	params := &predictop.PredictParams{
		PositionList: positionList,
	}
	params.WithTimeout(nr.Config.PredictConfig.Timeout)

	predictData, err := nr.PredictClient.Operations.Predict(params)
	if err != nil {
		return nil, err
	}

	return predictData, nil
}

func (nr *NearestRoute) getMinTimeTransfer(times []int64) (int64, error) {
	if len(times) == 0 {
		return -1, fmt.Errorf("No predict data")
	}

	var minTime = times[0]
	for i := 1; i < len(times); i++ {
		if times[i] < minTime {
			minTime = times[i]
		}
	}
	return minTime, nil
}

func (nr *NearestRoute) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pos, err := nr.getCurrentPos(r.Body)
	if err != nil {
		log.Fatalf("Error when reading body from POST request: %s\n", err.Error())
	}

	carsData, err := nr.getCars(pos)
	if err != nil {
		log.Fatalf("Error when receiving data from cars client: %s\n", err.Error())
	}
	// DEBUG
	fmt.Printf("%#v\n\n", carsData.Payload)

	predictData, err := nr.getPredicts(pos, carsData)
	if err != nil {
		log.Fatalf("Error when receiving data from predict client: %s\n", err.Error())
	}
	// DEBUG
	fmt.Printf("%#v\n", predictData.Payload)

	minTime, err := nr.getMinTimeTransfer(predictData.Payload)
	if err != nil {
		log.Fatalf("Error when handle predict data: %s\n", err.Error())
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"response": "` + strconv.Itoa(int(minTime)) + `"}`))
}

func NewServer(port string) (*TransferServer, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, err
	}

	carsClient := initCarsClient(cfg)
	predictClient := initPredictClient(cfg)

	routeHandler := &NearestRoute{
		Config:        cfg,
		CarsClient:    carsClient,
		PredictClient: predictClient,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/transfer", routeHandler)
	httpServer := initHTTPServer(mux, cfg)

	return &TransferServer{
		httpServer,
	}
}
