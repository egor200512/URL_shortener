package analytics

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

// Metrics conforms to the proto but real Prometheus scraping should hit the /metrics HTTP endpoint directly.
func (service *analyticsService) Metrics(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
