package events

import "github.com/GermanBogatov/tg_bot/internal/events/model"

type ProcessEventStrategy interface {
	Process(eventBody []byte) (model.ProcessedEvent, error)
}
