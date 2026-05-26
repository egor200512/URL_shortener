package service

import (
	"context"

	desc "github.com/egor200512/URL_shortener/shared/gen/analytics"
)

type IAnalyticsService interface {
	GetLinkStats(ctx context.Context, shortLink string) (*desc.GetLinkStatsResponse, error)
}
