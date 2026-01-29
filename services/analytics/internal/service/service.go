package service

import (
	"context"

	gonats "github.com/nats-io/nats.go"

	desc "github.com/egor200512/URL_shortener/shared/gen/analytics"
	"google.golang.org/protobuf/types/known/emptypb"
)

type IAnalyticsService interface {
	GetEvents(ctx context.Context, req *desc.GetEventsRequest) (*desc.GetEventsResponse, error)
	Health(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error)
	Metrics(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error)

	HandleMessage(ctx context.Context, msg *gonats.Msg) error
}
