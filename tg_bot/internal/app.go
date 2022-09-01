package internal

import (
	"context"
	"fmt"
	"github.com/GermanBogatov/tgbot/internal/config"
	"github.com/GermanBogatov/tgbot/internal/service"
	"github.com/GermanBogatov/tgbot/pkg/client/youtube"
	"github.com/GermanBogatov/tgbot/pkg/logging"
	tele "gopkg.in/telebot.v3"
	"net/http"
	"time"
)

type app struct {
	cfg            *config.Config
	logger         *logging.Logger
	httpServer     *http.Server
	youtubeService service.YoutubeService
}

func NewApp(logger *logging.Logger, cfg *config.Config) (App, error) {
	logger.Println("init Youtube client")
	youtubeClient := youtube.NewClient(cfg.Youtube.APIURL, cfg.Youtube.AccessToken, &http.Client{})
	youtubeService := service.NewYoutubeService(youtubeClient, logger)

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
	a.startBot()
}

func (a *app) startBot() {
	pref := tele.Settings{
		Token:  a.cfg.Telegram.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		a.logger.Fatal(err)
		return
	}

	b.Handle("/yt", func(c tele.Context) error {
		trackname := c.Message().Payload
		name, err := a.youtubeService.FindTrackByName(context.Background(), trackname)
		if err != nil {
			return c.Send("Твой трек не найден")
		}
		return c.Send(fmt.Sprintf("This is your track: %s", name))
	})

	b.Start()
}
