package internal

import (
	"errors"
	"fmt"
	"github.com/GermanBogatov/youtube_service/internal/config"
	"github.com/GermanBogatov/youtube_service/internal/events"
	youtube2 "github.com/GermanBogatov/youtube_service/internal/youtube"
	"github.com/GermanBogatov/youtube_service/pkg/client/mq/rabbitmq"
	"github.com/GermanBogatov/youtube_service/pkg/client/youtube"
	"github.com/GermanBogatov/youtube_service/pkg/logging"
	"github.com/GermanBogatov/youtube_service/pkg/shutdown"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"syscall"
	"time"
)

type app struct {
	cfg            *config.Config
	logger         *logging.Logger
	httpServer     *http.Server
	youtubeService youtube2.Service
	router         *httprouter.Router
}

func NewApp(logger *logging.Logger, cfg *config.Config) (App, error) {

	router := httprouter.New()
	logger.Println("init Youtube client")
	youtubeClient := youtube.NewClient(cfg.Youtube.APIURL, cfg.Youtube.AccessToken, &http.Client{})
	youtubeService := youtube2.NewService(youtubeClient, logger)

	return &app{
		cfg:            cfg,
		logger:         logger,
		youtubeService: youtubeService,
		router:         router,
	}, nil
}

type App interface {
	Run()
}

func (a *app) Run() {
	a.startConsume()
	a.startHTTP()
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
	consumer.DeclareQueue(a.cfg.RabbitMQ.Consumer.Queue, true, false, false, nil)
	if err != nil {
		a.logger.Fatal(err)
	}

	messages, err := consumer.Consume(a.cfg.RabbitMQ.Consumer.Queue)
	if err != nil {
		a.logger.Fatal(err)
	}

	for i := 0; i < a.cfg.AppConfig.EventWorkers; i++ {
		worker := events.NewWorker(i, consumer, a.cfg.RabbitMQ.Producer.Queue, producer, messages, a.logger, a.youtubeService)

		go worker.Process()
		a.logger.Infof("EVent Worker #%d statred", i)
	}
}

func (a *app) startHTTP() {
	a.logger.Info("start HTTP")

	var listener net.Listener

	if a.cfg.Listen.Type == "sock" {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			a.logger.Fatal(err)
		}
		socketPath := path.Join(appDir, "app.sock")
		a.logger.Infof("socket path: %s", socketPath)

		a.logger.Info("create and listen unix socket")
		listener, err = net.Listen("unix", socketPath)
		if err != nil {
			a.logger.Fatal(err)
		}
	} else {
		a.logger.Infof("bind application to host: %s and port: %s", a.cfg.Listen.BindIP, a.cfg.Listen.Port)
		var err error
		listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", a.cfg.Listen.BindIP, a.cfg.Listen.Port))
		if err != nil {
			a.logger.Fatal(err)
		}
	}

	c := cors.New(cors.Options{
		AllowedMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodOptions, http.MethodDelete},
		AllowedOrigins:     []string{"*"},
		AllowCredentials:   true,
		AllowedHeaders:     []string{"Access-Token", "Refresh-Token", "Authorization", "Location", "Charset", "Access-Control-Allow-Origin", "Content-Type", "content-type", "Origin", "Accept", "Content-Length", "Accept-Encoding", "X-CSRF-Token"},
		OptionsPassthrough: true,
		ExposedHeaders:     []string{"Access-Token", "Refresh-Token", "Location", "Authorization", "Content-Disposition"},
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})

	handler := c.Handler(a.router)

	a.httpServer = &http.Server{
		Handler:      handler,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go shutdown.Graceful([]os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM},
		a.httpServer)

	a.logger.Println("application completely initialized and started")

	if err := a.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			a.logger.Warn("server shutdown")
		default:
			a.logger.Fatal(err)
		}
	}
}
