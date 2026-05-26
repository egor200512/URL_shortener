package handler

import (
	"context"
	"fmt"

	desc "github.com/egor200512/URL_shortener/shared/gen/analytics"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (handler *AnalyticsHandler) GetLinkStats(ctx context.Context, in *desc.GetLinkStatsRequest) (*desc.GetLinkStatsResponse, error) {
	if in.ShortLink == "" {
		return nil, status.Error(codes.InvalidArgument, "short link is required")
	}

	stats, err := handler.analyticsService.GetLinkStats(ctx, in.ShortLink)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get link stats: %s", err.Error()))
	}

	return stats, nil
}
