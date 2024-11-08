package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arthben/BackendGolang/api-gateway/api/handlers"
	"github.com/arthben/BackendGolang/api-gateway/internal/config"
	"github.com/arthben/BackendGolang/api-gateway/internal/database"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

// @title Indego & Open Weather API Documentation
// @version 1.0
// @BasePath /

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	initLogger()

	// init database
	dbPool, err := database.NewPool(cfg)
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}

	defer dbPool.Close()

	// init request handler
	handler, err := handlers.NewHandler(dbPool, cfg).BuildHandler()
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}

	// timeout := cfg.App.WaitTimeOut
	server := http.Server{
		Addr:         "0.0.0.0:" + cfg.App.Port,
		Handler:      handler,
		ReadTimeout:  time.Duration(cfg.App.WaitTimeOut) * time.Second,
		WriteTimeout: time.Duration(cfg.App.WaitTimeOut) * time.Second,
	}
	defer server.Close()

	osSignal := make(chan os.Signal, 2)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info().Str("SERVER START", "address : 0.0.0.0:"+cfg.App.Port)
		log.Warn().Str("SERVER CLOSE", "error : "+server.ListenAndServe().Error())
	}()

	// wait until server closed
	select {
	case <-osSignal:
		log.Warn().Msg("Server Closed Interrupted by OS")
	}

	log.Info().Msg("SERVER STOP")
}

func initLogger() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = zerolog.New(&lumberjack.Logger{
		Filename:   "log/api_gateway.log",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     30,
		Compress:   true,
	})
	log.Logger = log.With().Str("service", "api-gateway").Logger()
	log.Logger = log.With().Timestamp().Logger()
}
