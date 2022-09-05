package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/GermanBogatov/youtube_service/pkg/client/youtube"
	"github.com/GermanBogatov/youtube_service/pkg/logging"
)

type service struct {
	Client youtube.Client
	logger *logging.Logger
}

func NewService(client youtube.Client, logger *logging.Logger) *service {
	return &service{Client: client, logger: logger}
}

type Service interface {
	FindTrackByName(ctx context.Context, trackName string) (string, error)
}

func (s *service) FindTrackByName(ctx context.Context, trackName string) (string, error) {
	response, err := s.Client.SearchTrack(ctx, trackName)
	if err != nil {
		return "", err
	}

	var responseData map[string]interface{}

	if err = json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		return "", err
	}

	if response.StatusCode != 200 {
		s.logger.Error(responseData["error"].(map[string]interface{})["message"].(string))
		return "", fmt.Errorf("Fail")
	}
	//s.logger.Info(responseData)
	//return "LINK", nil
	if a, ok := responseData["items"].([]interface{}); ok {
		b := a[0].(map[string]interface{})["id"].(map[string]interface{})["videoId"].(string)
		return fmt.Sprintf("https://music.youtube.com/watch?v=%s", b), nil
	} else {
		return "", fmt.Errorf("yotube request failed due to error %v", responseData)
	}
}
