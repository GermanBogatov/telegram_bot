package internal

import (
	"encoding/json"
	"fmt"
	"github.com/GermanBogatov/tg_bot/internal/config"
	"github.com/GermanBogatov/tg_bot/internal/events"
	"github.com/GermanBogatov/tg_bot/internal/events/youtube"
	"github.com/GermanBogatov/tg_bot/internal/service/bot"
	"github.com/GermanBogatov/tg_bot/pkg/client/mq"
	"github.com/GermanBogatov/tg_bot/pkg/client/mq/rabbitmq"
	"github.com/GermanBogatov/tg_bot/pkg/logging"
	tele "gopkg.in/telebot.v3"
	"net/http"
	"time"
)

type app struct {
	cfg                    *config.Config
	logger                 *logging.Logger
	httpServer             *http.Server
	producer               mq.Producer
	youtubeProcessStrategy events.ProcessEventStrategy
	bot                    *tele.Bot
}

func NewApp(logger *logging.Logger, cfg *config.Config) (App, error) {

	return &app{
		cfg:                    cfg,
		logger:                 logger,
		youtubeProcessStrategy: youtube.NewYoutubeProcessEventStrategy(logger),
	}, nil
}

type App interface {
	Run()
}

func (a *app) Run() {

	bot, err := a.createBot()
	if err != nil {
		return
	}
	a.bot = bot
	a.startConsume()
	a.bot.Start()
}

func (a *app) startConsume() {
	test := new(chan mq.Message)
	fmt.Println("test", test)
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
	err = consumer.DeclareQueue(a.cfg.RabbitMQ.Consumer.YouTubeQueue, true, false, false, nil)
	if err != nil {
		a.logger.Fatal(err)
	}
	youtubeMessages, err := consumer.Consume(a.cfg.RabbitMQ.Consumer.YouTubeQueue)
	if err != nil {
		a.logger.Fatal(err)
	}

	botservice := bot.Service{
		Bot:    a.bot,
		Logger: a.logger,
	}

	for i := 0; i < a.cfg.AppConfig.EventWorkers.YoutubeWorkers; i++ {
		worker := events.NewWorker(i, consumer, a.youtubeProcessStrategy, botservice, producer, youtubeMessages, a.logger)

		go worker.Process()
		a.logger.Infof("EVent Worker #%d statred", i)
	}

	a.logger.Println("start producer")
	a.producer = producer
}

func (a *app) createBot() (abot *tele.Bot, botErr error) {
	a.logger.Info("Init bot token")
	pref := tele.Settings{
		Token:   a.cfg.Telegram.Token,
		Poller:  &tele.LongPoller{Timeout: 60 * time.Second},
		Verbose: false,
		OnError: a.OnBotError,
	}

	a.logger.Info("Create NewBot")
	abot, botErr = tele.NewBot(pref)
	if botErr != nil {
		a.logger.Fatal(botErr)
		return
	}

	abot.Handle("/help", func(c tele.Context) error {
		return c.Send(fmt.Sprintf("/yt - find youtube track!"))
	})

	abot.Handle("/yt", func(c tele.Context) error {
		trackname := c.Message().Payload
		request := youtube.SearchTrackRequest{
			RequestID: fmt.Sprintf("%d", c.Sender().ID),
			Name:      trackname,
		}

		marshal, _ := json.Marshal(request)
		err := a.producer.Publish(a.cfg.RabbitMQ.Producer.YouTubeQueue, marshal)
		if err != nil {
			return c.Send(fmt.Sprintf("ошибка: %s", err.Error()))
		}
		return c.Send(fmt.Sprintf("Заявка принята"))

	})

	return

}

func (a *app) OnBotError(err error, ctx tele.Context) {
	a.logger.Error(err)
}
