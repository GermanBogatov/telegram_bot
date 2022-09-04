package events

import (
	"github.com/GermanBogatov/tg_bot/internal"
	"github.com/GermanBogatov/tg_bot/pkg/client/mq"
	"github.com/GermanBogatov/tg_bot/pkg/logging"
)

type worker struct {
	id              int
	client          mq.Consumer
	producer        mq.Producer
	responseQueue   string
	messages        <-chan mq.Message
	logger          *logging.Logger
	processStrategy ProcessEventStrategy
	botService      internal.BotService
}

//TODO попробовать вернуть структуру а не интерфейс
func NewWorker(id int, client mq.Consumer, processStrategy ProcessEventStrategy, producer mq.Producer, messages <-chan mq.Message, logger *logging.Logger) Worker {
	return &worker{id: id, client: client, processStrategy: processStrategy, producer: producer, messages: messages, logger: logger}
}

type Worker interface {
	Process()
}

func (w *worker) Process() {
	for msg := range w.messages {

		processedEvent, err := w.processStrategy.Process(msg.Body)
		if err != nil {
			w.logger.Errorf("[worker #%d]: failed to unmarshal event due to error %v", w.id, err)
			w.logger.Debugf("[worker #%d]: body: %s", w.id, msg.Body)
			w.reject(msg)
			return
		}

		err = w.botService.SendMessage(processedEvent)
		if err != nil {
			w.logger.Errorf("[worker #%d]: failed to send message due to error %v", w.id, err)
			w.logger.Debugf("[worker #%d]: body: %s", w.id, msg.Body)
			w.reject(msg)
			return
		}
		w.ack(msg)
	}

}

func (w *worker) reject(msg mq.Message) {
	if err := w.client.Reject(msg.ID, false); err != nil {
		w.logger.Errorf("[worker #%d]: failed to reject due to error %v", w.id, err)
	}
}

func (w *worker) ack(msg mq.Message) {
	if err := w.client.Ack(msg.ID, false); err != nil {
		w.logger.Errorf("[worker #%d]: failed to ACK due to error %v", w.id, err)
	}
}
