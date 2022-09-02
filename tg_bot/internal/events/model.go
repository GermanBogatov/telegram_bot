package events

type SearchTrackRequest struct {
	RequestID string `json:"request_id"`
	Name      string `json:"name"`
}

type SearchTrackResponse struct {
	RequestID string `json:"request_id"`
	Name      string `json:"name"`
	Success   string `json:"success"`
	Error     string `json:"err"`
}
