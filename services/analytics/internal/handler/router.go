package handler

import (
	"github.com/egor200512/URL_shortener/services/analytics/internal/service"
	desc "github.com/egor200512/URL_shortener/shared/gen/analytics"
)

type AnalyticsHandler struct {
	desc.UnimplementedAnalyticsServiceServer
	svc service.IAnalyticsService
}

func NewAnalyticsRouter(svc service.IAnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		svc: svc,
	}
}
