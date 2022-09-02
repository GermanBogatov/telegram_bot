package internal

import (
	"encoding/json"
	"fmt"
	"github.com/GermanBogatov/tg_bot/internal/config"
	"github.com/GermanBogatov/tg_bot/internal/events"
	"github.com/GermanBogatov/tg_bot/pkg/client/mq"
	"github.com/GermanBogatov/tg_bot/pkg/client/mq/rabbitmq"
	"github.com/GermanBogatov/tg_bot/pkg/logging"
	tele "gopkg.in/telebot.v3"
	"net/http"
	"time"
)

type app struct {
	cfg        *config.Config
	logger     *logging.Logger
	httpServer *http.Server
	bot        *tele.Bot
	producer   mq.Producer
}

func NewApp(logger *logging.Logger, cfg *config.Config) (App, error) {

	return &app{
		cfg:    cfg,
		logger: logger,
	}, nil
}

type App interface {
	Run()
}

func (a *app) Run() {
	a.startBot()
	a.startConsume()
	a.bot.Start()

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
		worker := events.NewWorker(i, consumer, producer, messages, a.logger, a.bot)

		go worker.Process()
		a.logger.Infof("EVent Worker #%d statred", i)
	}
	a.producer = producer
}

func (a *app) startBot() {
	pref := tele.Settings{
		Token:  a.cfg.Telegram.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	var botErr error
	a.bot, botErr = tele.NewBot(pref)
	if botErr != nil {
		a.logger.Fatal(botErr)
		return
	}

	a.bot.Handle("/yt", func(c tele.Context) error {
		trackname := c.Message().Payload
		request := events.SearchTrackRequest{
			RequestID: fmt.Sprintf("%d", c.Sender().ID),
			Name:      trackname,
		}

		marshal, _ := json.Marshal(request)
		err := a.producer.Publish(a.cfg.RabbitMQ.Producer.Queue, marshal)
		if err != nil {
			return c.Send(fmt.Sprintf("ошибка: %s", err.Error()))
		}
		return c.Send(fmt.Sprintf("Заявка принята"))

	})

}
