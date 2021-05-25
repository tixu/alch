package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/tixu/alch/config"
	"github.com/tixu/alch/server"
)

func main() {

	logger := logrus.New()
	// Set JSON as default for the moment
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Load configuration from file
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}

	// Configure logger
	err = config.ConfigureLogger(logger, cfg.Log)
	if err != nil {
		logger.Fatal(err)
		os.Exit(1)
	}
	logger.Info("Configuration successfully loaded and logger configured")
	server.NewServer(cfg, logger).Run()
}
