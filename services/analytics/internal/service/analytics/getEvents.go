package analytics

import (
	"context"
	"strings"

	desc "github.com/egor200512/URL_shortener/shared/gen/analytics"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (service *analyticsService) GetEvents(ctx context.Context, req *desc.GetEventsRequest) (*desc.GetEventsResponse, error) {
	items, err := service.repo.GetEvents(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	events := make([]*desc.LinkEvent, 0, len(items))
	for _, item := range items {
		events = append(events, &desc.LinkEvent{
			Type:         toProtoEventType(item.EventType),
			UserId:       item.UserID,
			ShortLink:    item.ShortLink,
			OriginalLink: item.OriginalLink,
			OccurredAt:   timestamppb.New(item.ExecutedAt),
		})
	}

	return &desc.GetEventsResponse{
		Events: events,
	}, nil
}

func toProtoEventType(eventType string) desc.LinkEventType {
	switch strings.ToLower(eventType) {
	case "created":
		return desc.LinkEventType_LINK_EVENT_TYPE_CREATED
	case "fetched":
		return desc.LinkEventType_LINK_EVENT_TYPE_FETCHED
	case "deleted":
		return desc.LinkEventType_LINK_EVENT_TYPE_DELETED
	default:
		return desc.LinkEventType_LINK_EVENT_TYPE_UNSPECIFIED
	}
}
