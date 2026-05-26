package handler

import (
	"context"
	"fmt"

	desc "github.com/egor200512/URL_shortener/shared/gen/analytics"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (handler *AnalyticsHandler) GetEvents(ctx context.Context, in *desc.GetEventsRequest) (*desc.GetEventsResponse, error) {
	if in.Limit <= 0 {
		return nil, status.Error(codes.InvalidArgument, "limit must be positive")
	}
	if in.Offset < 0 {
		return nil, status.Error(codes.InvalidArgument, "offset must be non-negative")
	}

	events, err := handler.analyticsService.GetEvents(ctx, in.Limit, in.Offset)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get events: %s", err.Error()))
	}

	return events, nil
}
