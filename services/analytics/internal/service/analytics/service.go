package analytics

import (
	"github.com/egor200512/URL_shortener/services/analytics/internal/repository"
	"github.com/egor200512/URL_shortener/services/analytics/internal/service"
)

type analyticsService struct {
	repo repository.IAnalyticsRepo
}

func NewAnalyticsService(repo repository.IAnalyticsRepo) service.IAnalyticsService {
	return &analyticsService{
		repo: repo,
	}
}
