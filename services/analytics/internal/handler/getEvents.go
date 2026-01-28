package handler

import (
	"context"
	"errors"

	desc "github.com/egor200512/URL_shortener/shared/gen/analytics"
)

func (h *AnalyticsHandler) GetEvents(ctx context.Context, req *desc.GetEventsRequest) (*desc.GetEventsResponse, error) {
	if req.Limit <= 0 {
		return nil, errors.New("limit must be positive")
	}
	if req.Offset < 0 {
		return nil, errors.New("offset must be non-negative")
	}

	return h.service.GetEvents(ctx, req)
}
