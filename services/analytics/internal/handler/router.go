package handler

import (
	"github.com/egor200512/URL_shortener/services/analytics/internal/service"
	desc "github.com/egor200512/URL_shortener/shared/gen/analytics"
)

type AnalyticsHandler struct {
	desc.UnimplementedAnalyticsServiceServer
	analyticsService service.IAnalyticsService
}

func NewAnalyticsRouter(analyticsService service.IAnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}
