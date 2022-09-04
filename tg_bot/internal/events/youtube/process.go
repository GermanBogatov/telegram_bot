package youtube

import (
	"encoding/json"
	"fmt"
	"github.com/GermanBogatov/tg_bot/internal/events"
	"github.com/GermanBogatov/tg_bot/pkg/logging"
)

type yt struct {
	logger *logging.Logger
}

func NewYoutubeProcessEventStrategy(logger *logging.Logger) events.ProcessEventStrategy {
	return &yt{
		logger: logger,
	}
}
func (p *yt) Process(eventBody []byte) (response events.ProccesedEvent, err error) {

	event := SearchTrackResponse{}
	if err = json.Unmarshal(eventBody, &event); err != nil {
		return response, fmt.Errorf("failed to unmarshal event due to error %v", err)

	}

	var eventErr error
	if event.Meta.Error != "" {
		eventErr = fmt.Errorf(event.Meta.Error)
	}
	return events.ProccesedEvent{
		RequestID: event.Meta.RequestID,
		Message:   event.Data.URL,
		Err:       eventErr,
	}, nil
}
