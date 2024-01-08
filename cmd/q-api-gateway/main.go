package main

import (
	"github.com/MeM0rd/q-api-gateway/internal/cron"
	"github.com/MeM0rd/q-api-gateway/internal/handlers/auth"
	"github.com/MeM0rd/q-api-gateway/internal/handlers/profile"
	"github.com/MeM0rd/q-api-gateway/internal/handlers/quote"
	"github.com/MeM0rd/q-api-gateway/pkg/client/postgres"
	logger "github.com/MeM0rd/q-api-gateway/pkg/logger"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"os"
	"time"
)

func init() {
	godotenv.Load(".env")

	cron.Start()

	postgres.Open()
}

func main() {
	r := httprouter.New()
	l := logger.New()

	defer cron.Stop()
	defer postgres.Close()
	defer l.Info("Main func closed")

	registerRoutes(r, l)

	start(r, l)
}

func start(r *httprouter.Router, logger *logger.Logger) {
	listener, err := net.Listen("tcp", os.Getenv("PORT"))
	if err != nil {
		logger.Fatalf("Erorr net.Listen: %v", err)
	}

	server := http.Server{
		Handler:      r,
		WriteTimeout: 15 * time.Second,
	}

	logger.Info("Server starting")
	err = server.Serve(listener)
	if err != nil {
		logger.Fatalf("Error server.Serve: %v", err)
	}
}

func registerRoutes(r *httprouter.Router, l *logger.Logger) {
	authHandler := auth.NewHandler(l)
	authHandler.Route(r)

	profileHandler := profile.NewHandler()
	profileHandler.Route(r)

	quoteHandler := quote.NewHandler(l)
	quoteHandler.Route(r)
}
