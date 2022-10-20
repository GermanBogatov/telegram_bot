package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"sync"
)

type Config struct {
	IsDebug       bool `yaml:"is_debug" env:"ST_BOT_IS_DEBUG" env-default:"false"`
	IsDevelopment bool `yaml:"is_development" env:"ST_BOT_IS_DEVELOPMENT" env-default:"false"`
	Telegram      struct {
		Token string `yaml:"token" env:"ST_BOT_TELEGRAM_TOKEN" env-required:"true"`
	}
	RabbitMQ struct {
		Host     string `yaml:"host" env:"ST_BOT_RABBIT_HOST" env-required:"true"`
		Port     string `yaml:"port" env:"ST_BOT_RABBIT_PORT" env-required:"true"`
		Username string `yaml:"username" env:"ST_BOT_RABBIT_USERNAME" env-required:"true"`
		Password string `yaml:"password" env:"ST_BOT_RABBIT_PASSWORD" env-required:"true"`
		Consumer struct {
			Youtube            string `yaml:"spotify" env:"ST_BOT_RABBIT_CONSUMER_YOUTUBE" env-required:"true"`
			Spotify            string `yaml:"spotify" env:"ST_BOT_RABBIT_CONSUMER_SPOTIFY" env-required:"true"`
			MessagesBufferSize int    `yaml:"messages_buff_size" env:"ST_BOT_RABBIT_CONSUMER_MBS" env-default:"100"`
		} `yaml:"consumer"`
		Producer struct {
			Youtube string `yaml:"spotify" env:"ST_BOT_RABBIT_PRODUCER_YOUTUBE" env-required:"true"`
			Spotify string `yaml:"spotify" env:"ST_BOT_RABBIT_PRODUCER_SPOTIFY" env-required:"true"`
		} `yaml:"producer"`
	}
	AppConfig AppConfig `yaml:"app"`
}

type AppConfig struct {
	EventWorkers struct {
		Youtube int `yaml:"spotify" env:"ST_BOT_EVENT_WORKERS_YT" env-default:"3"`
		Spotify int `yaml:"spotify" env:"ST_BOT_EVENT_WORKERS_SPOT" env-default:"3"`
	} `yaml:"event_workers"`
	LogLevel string `yaml:"log_level" env:"ST_BOT_LOG_LEVEL" env-default:"error"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {

	//TODO update configs!!!!
	os.Setenv("ST_BOT_IS_DEBUG", "true")
	os.Setenv("ST_BOT_IS_DEVELOPMENT", "true")

	os.Setenv("ST_BOT_TELEGRAM_TOKEN", "5462801617:AAGu5112sSwju_IBcdtsUicexWydh_kmXqg")

	os.Setenv("ST_BOT_RABBIT_HOST", "localhost")
	os.Setenv("ST_BOT_RABBIT_PORT", "5672")
	os.Setenv("ST_BOT_RABBIT_USERNAME", "guest")
	os.Setenv("ST_BOT_RABBIT_PASSWORD", "guest")

	os.Setenv("ST_BOT_RABBIT_CONSUMER_YOUTUBE", "yt-s-resp-events")
	os.Setenv("ST_BOT_RABBIT_CONSUMER_SPOTIFY", "spot-s-resp-events")
	os.Setenv("ST_BOT_RABBIT_CONSUMER_MBS", "100")
	os.Setenv("ST_BOT_RABBIT_PRODUCER_YOUTUBE", "yt-s-req-events")
	os.Setenv("ST_BOT_RABBIT_PRODUCER_SPOTIFY", "spot-s-req-events")
	os.Setenv("ST_BOT_EVENT_WORKERS_YT", "3")
	os.Setenv("ST_BOT_EVENT_WORKERS_SPOT", "3")
	os.Setenv("ST_BOT_LOG_LEVEL", "trace")

	once.Do(func() {
		instance = &Config{}

		if err := cleanenv.ReadEnv(instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			log.Print(help)
			log.Fatal(err)
		}
	})
	return instance
}
