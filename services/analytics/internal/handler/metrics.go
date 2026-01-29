package handler

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *AnalyticsHandler) Metrics(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return h.service.Metrics(ctx, &emptypb.Empty{})
}
