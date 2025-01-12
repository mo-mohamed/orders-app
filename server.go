package main

import (
	"log"
	"net/http"
	"os"

	"github.com/orders-app/handlers"
	"github.com/orders-app/logger"
	"github.com/orders-app/tracing"
)

func init() {
	env := os.Getenv("ACTIVE_ENV")
	if env == "" {
		env = "dev"
	}

	logger.InitLogger(env)
	defer logger.SyncLogger()

	logger.Log.Info("Initializing Tracer")
	cleanup := tracing.InitTracer()
	defer cleanup()
}

func main() {
	logger.Log.Info("Application Started")

	handler, err := handlers.New()
	if err != nil {
		log.Fatal(err)
	}
	router := handlers.ConfigureHandler(handler)
	logger.Log.Info("Listening on localhost:3000...")
	err = http.ListenAndServe(":3000", router)
	logger.Log.Fatal(err.Error())
}
