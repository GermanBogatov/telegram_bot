package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug"`
	Listen  struct {
		Type   string `yaml:"type" env-default:"port"`
		BindIP string `yaml:"bind_ip" env-default:"localhost"`
		Port   string `yaml:"port" env-default:"8080"`
	} `yaml:"listen" env-required:"true"`
	Telegram struct {
		Token string `yaml:"token" env:"ST-BOT-TelegramToken" env-required:"true"`
	}
	RabbitMQ struct {
		Host     string `yaml:"host" env:"YTS_RABBIT_HOST" env-required:"true"`
		Port     string `yaml:"port" env:"YTS_RABBIT_PORT" env-required:"true"`
		Username string `yaml:"username" env:"YTS_RABBIT_USERNAME" env-required:"true"`
		Password string `yaml:"password" env:"YTS_RABBIT_PASSWORD" env-required:"true"`
		Consumer struct {
			Queue              string `yaml:"queue" env:"YTS_RABBIT_CONSUMER_QUEUE" env-required:"true"`
			MessagesBufferSize int    `yaml:"messagesbuffersize" env:"YTS_RABBIT_CONSUMER_MBS" env-default:"100"`
		} `yaml:"consumer" env-required:"true"`
		Producer struct {
			YouTubeQueue string `yaml:"ytqueue" env:"YTS_Rabbit_PRODUCERQUEUE" env-required:"true"`
		} `yaml:"producer" env-required:"true"`
	} `yaml:"rabbitMQ"`
	AppConfig AppConfig `yaml:"app" env-required:"true"`
}

type AppConfig struct {
	EventWorkers struct {
		YoutubeWorkers int `yaml:"youtubeWorkers" env:"ST-BOT-YoutubeEventWorkers" env-default:"3" env-required:"true"`
	} `yaml:"event_workers"`
	LogLevel string `yaml:"log_level" env:"ST-BOT-loglevel" env-default:"error" env-required:"true"`
}

var instance *Config
var once sync.Once

func GetConfig(path string) *Config {

	once.Do(func() {
		log.Printf("read application config in path %s", path)
		instance = &Config{}

		if err := cleanenv.ReadConfig(path, instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			log.Print(help)
			log.Fatal(err)
		}
	})
	return instance
}
