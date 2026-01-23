package analytics

import (
	"context"

	desc "github.com/egor200512/URL_shortener/shared/gen/analytics"
)

func (service *analyticsService) GetCounters(ctx context.Context, req *desc.GetCountersRequest) (*desc.GetCountersResponse, error) {
	return &desc.GetCountersResponse{}, nil
}
