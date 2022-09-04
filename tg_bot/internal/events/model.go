package events

import "fmt"

type ResponseMessage struct {
	Meta ResponseMeta `json:"meta"`
	Data interface{}  `json:"data"`
}

type ResponseMeta struct {
	RequestID string `json:"request_id"`
	Success   bool   `json:"success"`
	Error     string `json:"err"`
}

func (m *ResponseMeta) String() string {
	return fmt.Sprintf("RequestID: %s, Error: %s", m.RequestID, m.Error)
}

type ProccesedEvent struct {
	RequestID string
	Message   string
	Err       error
}

func (m *ProccesedEvent) String() string {
	return fmt.Sprintf("RequestID: %s, Message: %s, Error: %s", m.RequestID, m.Message, m.Err)
}
