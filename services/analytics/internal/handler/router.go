package handler

import (
	"github.com/egor200512/URL_shortener/services/analytics/internal/service"
	desc "github.com/egor200512/URL_shortener/shared/gen/analytics"
)

type AnalyticsHandler struct {
	desc.UnimplementedAnalyticsServiceServer
	service service.IAnalyticsService
}

func NewAnalyticsRouter(service service.IAnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		service: service,
	}
}
