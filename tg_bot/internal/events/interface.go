package events

type ProcessEventStrategy interface {
	Process(eventBody []byte) (ProccesedEvent, error)
}
