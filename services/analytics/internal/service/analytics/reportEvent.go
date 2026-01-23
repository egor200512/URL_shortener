package analytics

import (
	"context"

	desc "github.com/egor200512/URL_shortener/shared/gen/analytics"
)

func (service *analyticsService) ReportEvent(ctx context.Context, req *desc.LinkEvent) (*desc.ReportEventResponse, error) {
	return &desc.ReportEventResponse{}, nil
}
