package repository

import (
	"fmt"
	"log"

	"github.com/royhairul/live-studio-api/internal/clients/shopee"
	"github.com/royhairul/live-studio-api/internal/clients/shopee/params"
)

type ShopeeLiveRepositoryImpl struct {
	client *shopee.ShopeeClient
}

func NewShopeeLiveRepository(client *shopee.ShopeeClient) ShopeeLiveRepository {
	return &ShopeeLiveRepositoryImpl{client}
}

// GetShopeeLiveRealTime implements ShopeeLiveRepository.
func (s *ShopeeLiveRepositoryImpl) GetShopeeLiveRealTime(cookie string) ([]params.ShopeeLiveReportItemRT, error) {
	endpoint := "/supply/api/lm/sellercenter/realtime/sessionList"
	query := map[string]string{
		"page":     "1",
		"pageSize": "10",
		"name":     "",
		"orderBy":  "",
		"sort":     "desc",
	}

	req, err := s.client.NewRequest("GET", endpoint, query, nil, cookie)
	if err != nil {
		return nil, fmt.Errorf("failed to create request to %s: %w", endpoint, err)
	}
	log.Println("RESULT", req)

	var result params.ShopeeApiResponse[params.ShopeeLiveRealTimeResponse]

	if err := s.client.DoRequest(req, &result); err != nil {
		return nil, fmt.Errorf("failed to do request to %s: %w", endpoint, err)
	}

	if result.Error != 0 {
		return nil, fmt.Errorf(result.ErrorMsg)
	}

	return result.Data.List, nil
}
