package service

import (
	"context"

	desc "github.com/egor200512/URL_shortener/shared/gen/analytics"
)

type IAnalyticsService interface {
	RecordEvent(ctx context.Context, req *desc.RecordEventRequest) error
	GetEvents(ctx context.Context, limit, offset int32) (*desc.GetEventsResponse, error)
	GetLinkStats(ctx context.Context, shortLink string) (*desc.GetLinkStatsResponse, error)
}
