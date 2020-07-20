package transfer

import (
	"net/http"

	"wheely/test/pkg/transfer/client/cars"
	"wheely/test/pkg/transfer/client/predict"
	"wheely/test/pkg/transfer/handlers"
	"wheely/test/pkg/transfer/utils"
)

/*
type Transfer interface {
	Predict() int64
}
*/

func NewServer(cfg *utils.Config) *http.Server {
	routeHandler := &handlers.Route{
		Config:        cfg,
		CarsClient:    cars.NewClient(cfg),
		PredictClient: predict.NewClient(cfg),
	}

	mux := http.NewServeMux()
	mux.Handle("/transfer", routeHandler)

	return &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  cfg.ReadTimout,
		WriteTimeout: cfg.WriteTimeout,
	}
}
