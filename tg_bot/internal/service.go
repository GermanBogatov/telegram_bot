package internal

import (
	"fmt"
	"github.com/GermanBogatov/tg_bot/internal/events"
	"github.com/GermanBogatov/tg_bot/pkg/logging"
	tele "gopkg.in/telebot.v3"
	"strconv"
)

type BotService struct {
	bot    *tele.Bot
	logger *logging.Logger
}

func (bs *BotService) SendMessage(data events.ProccesedEvent) error {
	i, _ := strconv.ParseInt(data.RequestID, 10, 64)
	id, err := bs.bot.ChatByID(i)
	if err != nil {
		bs.logger.Tracef("Bot send ResponseMessage ProcessedEvent: %s", data)
		return fmt.Errorf("failed to get chat by id due to error %v", err)
	}

	message := data.Message
	if data.Err != nil {
		message = fmt.Sprintf("Запрос не обработан, произошла ошибка (%s)", data.Err)
	}

	_, err = bs.bot.Send(id, message)
	if err != nil {
		bs.logger.Tracef("ChatID: %d, Data: %s", id.ID, data)
		return fmt.Errorf("failed to get send due to error %v", err)
	}

	return nil
}
