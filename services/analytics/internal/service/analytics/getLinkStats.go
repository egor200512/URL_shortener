package analytics

import (
	"context"

	desc "github.com/egor200512/URL_shortener/shared/gen/analytics"
)

func (service *analyticsService) GetLinkStats(ctx context.Context, shortLink string) (*desc.GetLinkStatsResponse, error) {
	stats, err := service.repo.GetLinkStats(ctx, shortLink)
	if err != nil {
		return nil, err
	}

	return &desc.GetLinkStatsResponse{
		ShortLink:    stats.ShortLink,
		CreatedCount: stats.CreatedCount,
		FetchedCount: stats.FetchedCount,
		DeletedCount: stats.DeletedCount,
		TotalCount:   stats.TotalCount,
	}, nil
}
