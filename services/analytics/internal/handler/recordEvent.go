package handler

import (
	"context"
	"fmt"

	desc "github.com/egor200512/URL_shortener/shared/gen/analytics"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (handler *AnalyticsHandler) RecordEvent(ctx context.Context, in *desc.RecordEventRequest) (*emptypb.Empty, error) {
	if in.EventType == "" {
		return nil, status.Error(codes.InvalidArgument, "event type is required")
	}
	if in.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}
	if in.ShortLink == "" {
		return nil, status.Error(codes.InvalidArgument, "short link is required")
	}
	if in.OriginalLink == "" {
		return nil, status.Error(codes.InvalidArgument, "original link is required")
	}

	if err := handler.analyticsService.RecordEvent(ctx, in); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to record event: %s", err.Error()))
	}

	return &emptypb.Empty{}, nil
}
