package transfer

import (
	"net/http"

	"wheely/test/pkg/transfer/handlers"
	"wheely/test/pkg/transfer/utils"
)

func NewServer(cfg *utils.Config) *http.Server {
	routeHandler := handlers.NewRouteHandler(cfg)

	mux := http.NewServeMux()
	mux.Handle("/transfer", routeHandler)

	return &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  cfg.ReadTimout,
		WriteTimeout: cfg.WriteTimeout,
	}
}
