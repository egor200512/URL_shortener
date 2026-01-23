package service

import (
	"context"

	desc "github.com/egor200512/URL_shortener/shared/gen/analytics"
	"google.golang.org/protobuf/types/known/emptypb"
)

type IAnalyticsService interface {
	ReportEvent(ctx context.Context, req *desc.LinkEvent) (*desc.ReportEventResponse, error)
	GetCounters(ctx context.Context, req *desc.GetCountersRequest) (*desc.GetCountersResponse, error)
	Health(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error)
}
