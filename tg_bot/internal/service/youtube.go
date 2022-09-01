package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/GermanBogatov/tg_bot/pkg/client/youtube"
	"github.com/GermanBogatov/tg_bot/pkg/logging"
)

type Service struct {
	Client youtube.Client
	logger *logging.Logger
}

func NewYoutubeService(client youtube.Client, logger *logging.Logger) *Service {
	return &Service{Client: client, logger: logger}
}

type YoutubeService interface {
	FindTrackByName(ctx context.Context, trackName string) (string, error)
}

func (s *Service) FindTrackByName(ctx context.Context, trackName string) (string, error) {
	response, err := s.Client.SearchTrack(ctx, trackName)
	if err != nil {
		return "", err
	}

	var responseData map[string]interface{}

	if err = json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		return "", err
	}

	//s.logger.Info(responseData)
	//return "LINK", nil
	a := responseData["items"].([]interface{})
	b := a[0].(map[string]interface{})["id"].(map[string]interface{})["videoId"].(string)
	return fmt.Sprintf("https://music.youtube.com/watch?v=%s", b), nil
}
