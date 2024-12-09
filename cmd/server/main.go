package main

import (
	"fmt"
	"github.com/OmgAbear/gosolve/internal/config"
	"github.com/rs/cors"

	httpInterface "github.com/OmgAbear/gosolve/internal/http_interface"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func main() {
	router := mux.NewRouter()
	router.Use(httpInterface.LoggingMiddleware)
	httpInterface.RegisterRoutes(router)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handlerWithCORS := c.Handler(router)

	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}
	logger := config.GetLogger()

	// start the http server without any recovery, graceful shutdown etc
	if err := http.ListenAndServe(":"+cfg.Server.Port, handlerWithCORS); err != nil {
		logger.Error(err.Error())
	}

}
