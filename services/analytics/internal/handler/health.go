package handler

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *AnalyticsHandler) Health(ctx context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
