package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug" env:"YTS-IsDebug" env-default:"false" env-required:"true"`
	Listen  struct {
		Type   string `yaml:"type" env:"YTS-ListenType" env-default:"port"`
		BindIP string `yaml:"bind_ip" env:"YTS-BindIP" env-default:"localhost"`
		Port   string `yaml:"port" env:"YTS-Port" env-default:"8080"`
	} `yaml:"listen" env-required:"true"`
	Youtube struct {
		APIURL      string `yaml:"api_url" env:"ST-BOT-YoutubeAPIURL" env-required:"true"`
		AccessToken string `yaml:"access_token" env:"ST-BOT-YoutubeAccessToken" env-required:"true"`
	} `yaml:"youtube"`
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
			Queue string `yaml:"queue" env:"YTS_Rabbit_PRODUCERQUEUE" env-required:"true"`
		} `yaml:"producer" env-required:"true"`
	} `yaml:"rabbitMQ"`
	AppConfig AppConfig `yaml:"app" env-required:"true"`
}

type AppConfig struct {
	EventWorkers int    `yaml:"event_workers" env:"ST-BOT-EventWorkers" env-default:"3" env-required:"true"`
	LogLevel     string `yaml:"log_level" env:"ST-BOT-Loglevel" env-default:"error" env-required:"true"`
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
