package handler

import (
	"context"
	"fmt"
	"log/slog"

	desc "github.com/egor200512/URL_shortener/shared/gen/analytics"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (handler *AnalyticsHandler) RecordEvent(ctx context.Context, in *desc.RecordEventRequest) (*emptypb.Empty, error) {
	slog.Info("record event endpoint called", "event_type", in.EventType, "user_id", in.UserId, "short_link", in.ShortLink)

	if in.EventType == "" {
		slog.Warn("record event validation failed", "error", "event type is required")
		return nil, status.Error(codes.InvalidArgument, "event type is required")
	}
	if in.UserId == "" {
		slog.Warn("record event validation failed", "event_type", in.EventType, "error", "user id is required")
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}
	if in.ShortLink == "" {
		slog.Warn("record event validation failed", "event_type", in.EventType, "user_id", in.UserId, "error", "short link is required")
		return nil, status.Error(codes.InvalidArgument, "short link is required")
	}
	if in.OriginalLink == "" {
		slog.Warn("record event validation failed", "event_type", in.EventType, "user_id", in.UserId, "short_link", in.ShortLink, "error", "original link is required")
		return nil, status.Error(codes.InvalidArgument, "original link is required")
	}

	if err := handler.analyticsService.RecordEvent(ctx, in); err != nil {
		slog.Error("record event failed", "event_type", in.EventType, "user_id", in.UserId, "short_link", in.ShortLink, "error", err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to record event: %s", err.Error()))
	}

	slog.Info("record event completed", "event_type", in.EventType, "user_id", in.UserId, "short_link", in.ShortLink)
	return &emptypb.Empty{}, nil
}
