package events

import (
	"context"
	"encoding/json"
	"github.com/GermanBogatov/youtube_service/internal/youtube"
	"github.com/GermanBogatov/youtube_service/pkg/client/mq"
	"github.com/GermanBogatov/youtube_service/pkg/logging"
	"strconv"
)

type worker struct {
	id            int
	client        mq.Consumer
	producer      mq.Producer
	responseQueue string
	messages      <-chan mq.Message
	logger        *logging.Logger
	service       youtube.Service
}

//TODO попробовать вернуть структуру а не интерфейс
func NewWorker(id int, client mq.Consumer, producer mq.Producer, messages <-chan mq.Message, logger *logging.Logger, service youtube.Service) Worker {
	return &worker{id: id, client: client, producer: producer, messages: messages, logger: logger, service: service}
}

type Worker interface {
	Process()
}

func (w *worker) Process() {
	for msg := range w.messages {
		event := SearchTrack{}
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			w.logger.Errorf("[worker #%d]: failed to unmarshal event due to error %v", w.id, err)
			w.logger.Errorf("[worker #%d]: body: %s", w.id, msg.Body)

			w.reject(msg)
			continue
		}
		respData := map[string]string{
			"request_id": event.RequestID,
		}
		name, err := w.service.FindTrackByName(context.TODO(), event.Name)
		if err != nil {
			respData["err"] = err.Error()
		} else {
			respData["name"] = name

		}

		respData["success"] = strconv.FormatBool(err == nil)

		w.sendResponse(respData)
	}
}

func (w *worker) sendResponse(d map[string]string) {
	b, err := json.Marshal(d)
	if err != nil {
		w.logger.Errorf("[worker #%d]: failed to reject due to error %v", w.id, err)
		return
	}

	err = w.producer.Publish(w.responseQueue, b)
	if err != nil {
		w.logger.Errorf("[worker #%d]: failed to response due to error %v", w.id, err)
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
