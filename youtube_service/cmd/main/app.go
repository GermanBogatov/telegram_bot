package main

import (
	"flag"
	"github.com/GermanBogatov/youtube_service/internal"
	"github.com/GermanBogatov/youtube_service/internal/config"
	"github.com/GermanBogatov/youtube_service/pkg/logging"
	"log"
)

var cfgPath string

func init() {
	flag.StringVar(&cfgPath, "config", "configs/dev.yml", "config file path")

}
func main() {
	flag.Parse()
	log.Print("config init")
	cfg := config.GetConfig(cfgPath)

	log.Print("logger init")
	logging.Init(cfg.AppConfig.LogLevel)
	logger := logging.GetLogger()

	logger.Println("Creating Application")
	app, err := internal.NewApp(logger, cfg)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Println("Running Applications")
	app.Run()
}