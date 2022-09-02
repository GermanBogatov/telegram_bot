package internal

import (
	"github.com/GermanBogatov/youtube_service/internal/config"
	"github.com/GermanBogatov/youtube_service/internal/events"
	youtube2 "github.com/GermanBogatov/youtube_service/internal/youtube"
	"github.com/GermanBogatov/youtube_service/pkg/client/mq/rabbitmq"
	"github.com/GermanBogatov/youtube_service/pkg/client/youtube"
	"github.com/GermanBogatov/youtube_service/pkg/logging"
	"net/http"
)

type app struct {
	cfg            *config.Config
	logger         *logging.Logger
	httpServer     *http.Server
	youtubeService youtube2.Service
}

func NewApp(logger *logging.Logger, cfg *config.Config) (App, error) {

	logger.Println("init Youtube client")
	youtubeClient := youtube.NewClient(cfg.Youtube.APIURL, cfg.Youtube.AccessToken, &http.Client{})
	youtubeService := youtube2.NewService(youtubeClient, logger)

	return &app{
		cfg:            cfg,
		logger:         logger,
		youtubeService: youtubeService,
	}, nil
}

type App interface {
	Run()
}

func (a *app) Run() {
	a.startConsume()
}

func (a *app) startConsume() {
	a.logger.Info("Start consumer")
	consumer, err := rabbitmq.NewRabbitMQConsumer(rabbitmq.ConsumerConfig{
		BaseConfig: rabbitmq.BaseConfig{
			Host:     a.cfg.RabbitMQ.Host,
			Port:     a.cfg.RabbitMQ.Port,
			Username: a.cfg.RabbitMQ.Username,
			Password: a.cfg.RabbitMQ.Password,
		},
		PrefetchCount: a.cfg.RabbitMQ.Consumer.MessagesBufferSize,
	})
	if err != nil {
		a.logger.Fatal(err)
	}

	a.logger.Info("Start producer")
	producer, err := rabbitmq.NewRabbitMQProducer(rabbitmq.ProducerConfig{
		BaseConfig: rabbitmq.BaseConfig{
			Host:     a.cfg.RabbitMQ.Host,
			Port:     a.cfg.RabbitMQ.Port,
			Username: a.cfg.RabbitMQ.Username,
			Password: a.cfg.RabbitMQ.Password,
		},
	})
	if err != nil {
		a.logger.Fatal(err)
	}

	messages, err := consumer.Consume(a.cfg.RabbitMQ.Consumer.Queue)
	if err != nil {
		a.logger.Fatal(err)
	}

	for i := 0; i < a.cfg.AppConfig.EventWorkers; i++ {
		worker := events.NewWorker(i, consumer, producer, messages, a.logger, a.youtubeService)

		go worker.Process()
		a.logger.Infof("EVent Worker #%d statred", i)
	}
}
