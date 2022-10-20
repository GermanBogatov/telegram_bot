package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"sync"
)

type Config struct {
	IsDebug       bool `env:"YTS_IS_DEBUG" env-default:"false"`
	IsDevelopment bool `env:"YTS_IS_DEVELOPMENT" env-default:"false"`
	YouTube       struct {
		APIURL          string `env:"YTS_YT_API_URL" env-required:"true"`
		RefreshTokenURL string `env:"YTS_YT_RefreshTokenURL" env-required:"true"`
		APIKey          string `env:"YTS_YT_APIKey" env-required:"true"`
		ClientID        string `env:"YTS_YT_CLIENT_ID" env-required:"true"`
		ClientSecret    string `env:"YTS_YT_CLIENT_SECRET" env-required:"true"`
		AccessToken     string `env:"YTS_YT_ACCESS_TOKEN" env-required:"true"`
		RefreshToken    string `env:"YTS_YT_REFRESH_TOKEN" env-required:"true"`
		AuthRedirectUri string `env:"YTS_YT_AUTH_REDIRECT_URI" env-required:"true"`
		AuthSuccessUri  string `env:"YTS_YT_AUTH_SUCCESS_URI" env-required:"true"`
		AccountsUri     string `env:"YTS_YT_ACCOUNTS_URI" env-required:"true"`
	}
	RabbitMQ struct {
		Host     string `env:"YTS_RABBIT_HOST" env-required:"true"`
		Port     string `env:"YTS_RABBIT_PORT" env-required:"true"`
		Username string `env:"YTS_RABBIT_USERNAME" env-required:"true"`
		Password string `env:"YTS_RABBIT_PASSWORD" env-required:"true"`
		Consumer struct {
			Queue              string `env:"YTS_RABBIT_CONSUMER_QUEUE" env-required:"true"`
			MessagesBufferSize int    `env:"YTS_RABBIT_CONSUMER_MBS" env-default:"100"`
		}
		Producer struct {
			Queue string `env:"YTS_Rabbit_PRODUCERQUEUE" env-required:"true"`
		}
	}
	AppConfig AppConfig
}

type AppConfig struct {
	EventWorkers int    `env:"YTS_EVENT_WORKERS" env-default:"3"`
	LogLevel     string `env:"YTS_LOG_LEVEL" env-default:"error"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {

	//TODO update configs!!!!
	os.Setenv("YTS_IS_DEBUG", "true")
	os.Setenv("YTS_IS_DEVELOPMENT", "true")

	os.Setenv("YTS_YT_API_URL", "https://www.googleapis.com/spotify/v3")
	os.Setenv("YTS_YT_RefreshTokenURL", "https://oauth2.googleapis.com/token")
	os.Setenv("YTS_YT_APIKey", "AIzaSyDojIFy7LSpuBvtsa3jcGbc-rMFY_7oHZ8")
	os.Setenv("YTS_YT_CLIENT_ID", "635230768534-kvce6o4aphfqbc1q52gfvj6kjmu3hrrb.apps.googleusercontent.com")
	os.Setenv("YTS_YT_CLIENT_SECRET", "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")
	os.Setenv("YTS_YT_ACCESS_TOKEN", "ya29.a0Aa4xrXO53r29VQn3jvVRotV1fMN1A7tfEUWToTWoMSVdWexVRGOfcieNirjScQypDZa0HUtqp5Fitc16hK0mw028EykK6-bu3_FyB76FBjmVBW0ihmGf3doBB6TsYbystlG4of0HOqWKtwH5yqQfRhSuWJBKaCgYKATASARMSFQEjDvL9JQzxEy2H6_3zan2DYi3Xzw0163")
	os.Setenv("YTS_YT_REFRESH_TOKEN", "1//0cqzUrCOxMXWvCgYIARAAGAwSNwF-L9IrGf2nItKWLHyvw9BofK-6Fhvb9rdTbsY0fxzw4j3f7dNKA1MSQb1ykkhiNDaIu7zfNHc")
	os.Setenv("YTS_YT_AUTH_REDIRECT_URI", "https://vk.com/boqatov")
	os.Setenv("YTS_YT_AUTH_SUCCESS_URI", "")
	os.Setenv("YTS_YT_ACCOUNTS_URI", "")

	os.Setenv("YTS_RABBIT_HOST", "localhost")
	os.Setenv("YTS_RABBIT_PORT", "5672")
	os.Setenv("YTS_RABBIT_USERNAME", "guest")
	os.Setenv("YTS_RABBIT_PASSWORD", "guest")

	os.Setenv("YTS_RABBIT_CONSUMER_QUEUE", "yt-s-req-events")
	os.Setenv("YTS_RABBIT_CONSUMER_MBS", "100")
	os.Setenv("YTS_Rabbit_PRODUCERQUEUE", "yt-s-resp-events")
	os.Setenv("YTS_EVENT_WORKERS", "3")
	os.Setenv("YTS_LOG_LEVEL", "trace")
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
