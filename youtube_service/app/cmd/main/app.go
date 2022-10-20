package main

import (
	"github.com/GermanBogatov/youtube_service/internal"
	"github.com/GermanBogatov/youtube_service/internal/config"
	"github.com/GermanBogatov/youtube_service/pkg/logging"
	"log"
)

func main() {
	log.Print("config initializing")
	cfg := config.GetConfig()

	log.Print("logger initializing")
	logging.Init(cfg.AppConfig.LogLevel)
	logger := logging.GetLogger()

	logger.Println("Creating Application")
	app, err := internal.NewApp(logger, cfg)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Println("Running Application")
	app.Run()
}
