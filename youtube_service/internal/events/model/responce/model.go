package responce

import "fmt"

type SearchTrack struct {
	Meta Meta `json:"meta"`
	Data Data `json:"data"`
}

type Data struct {
	URL string `json:"url"`
}

type Meta struct {
	RequestID string  `json:"request_id"`
	Error     *string `json:"err,omitempty"`
}

func (m *Meta) String() string {
	return fmt.Sprintf("RequestID: %s, Error: %s", m.RequestID, m.Error)
}
