package http_interface

import (
	"github.com/OmgAbear/gosolve/internal/config"
	"github.com/OmgAbear/gosolve/internal/infrastructure"
	"github.com/gorilla/mux"
	"net/http"
)

// RegisterRoutes sets the handlers for a specific route
func RegisterRoutes(router *mux.Router) {
	cfg, _ := config.GetConfig()
	logger := config.GetLogger()

	repoFactory := func() NumbersRepo {
		return infrastructure.NewNumbersRepo(cfg)
	}

	numbersHandler := NewNumbersHandler(repoFactory, logger)

	router.HandleFunc("/endpoint/{value}", numbersHandler.get).Methods(http.MethodGet)
}