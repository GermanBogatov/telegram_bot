package model

import "fmt"

type ResponseMessage struct {
	Meta ResponseMeta `json:"meta"`
	Data interface{}  `json:"data"`
}

type ResponseMeta struct {
	RequestID string  `json:"request_id"`
	Error     *string `json:"err,omitempty"`
}

func (m *ResponseMeta) String() string {
	return fmt.Sprintf("RequestID: %s, Error: %s", m.RequestID, m.Error)
}

type ProcesedEvent struct {
	RequestID string
	Message   string
	Err       error
}

func (m *ProcesedEvent) String() string {
	return fmt.Sprintf("RequestID: %s, Message: %s, Error: %s", m.RequestID, m.Message, m.Err)
}