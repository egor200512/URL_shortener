package analytics

import (
	"context"

	desc "github.com/egor200512/URL_shortener/shared/gen/analytics"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (service *analyticsService) GetEvents(ctx context.Context, limit, offset int32) (*desc.GetEventsResponse, error) {
	events, total, err := service.repo.GetEvents(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	resp := &desc.GetEventsResponse{
		Events:     make([]*desc.LinkEvent, 0, len(events)),
		TotalCount: total,
	}

	for _, event := range events {
		resp.Events = append(resp.Events, &desc.LinkEvent{
			Id:           event.ID,
			EventType:    event.EventType,
			UserId:       event.UserID,
			ShortLink:    event.ShortLink,
			OriginalLink: event.OriginalLink,
			ExecutedAt:   timestamppb.New(event.ExecutedAt),
		})
	}

	return resp, nil
}
