package handler

import (
	"context"
	"fmt"
	"log/slog"

	desc "github.com/egor200512/URL_shortener/shared/gen/analytics"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (handler *AnalyticsHandler) GetLinkStats(ctx context.Context, in *desc.GetLinkStatsRequest) (*desc.GetLinkStatsResponse, error) {
	slog.Info("get link stats endpoint called", "short_link", in.ShortLink)

	if in.ShortLink == "" {
		slog.Warn("get link stats validation failed", "error", "short link is required")
		return nil, status.Error(codes.InvalidArgument, "short link is required")
	}

	stats, err := handler.analyticsService.GetLinkStats(ctx, in.ShortLink)
	if err != nil {
		slog.Error("get link stats failed", "short_link", in.ShortLink, "error", err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get link stats: %s", err.Error()))
	}

	slog.Info("get link stats completed", "short_link", in.ShortLink)
	return stats, nil
}
