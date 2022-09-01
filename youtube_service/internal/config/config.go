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
	Youtube struct {
		APIURL      string `yaml:"api_url"`
		AccessToken string `yaml:"access_token"`
	}
	AppConfig AppConfig `yaml:"app" env-required:"true"`
}

type AppConfig struct {
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
