package analytics

import (
	r "github.com/egor200512/URL_shortener/services/analytics/internal/repository"
	s "github.com/egor200512/URL_shortener/services/analytics/internal/service"
)

type analyticsService struct {
	repo r.IAnalyticsRepo
}

func NewAnalyticsService(repo r.IAnalyticsRepo) s.IAnalyticsService {
	return &analyticsService{
		repo: repo,
	}
}
