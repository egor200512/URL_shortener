package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/egor200512/URL_shortener/services/analytics/internal/handler"
	mocks "github.com/egor200512/URL_shortener/services/analytics/internal/mocks"
	desc "github.com/egor200512/URL_shortener/shared/gen/analytics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandlerGetEvents(t *testing.T) {
	tests := []struct {
		name    string
		req     *desc.GetEventsRequest
		mockSvc func(*mocks.MockIAnalyticsService)
		wantErr bool
	}{
		{
			name: "success",
			req: &desc.GetEventsRequest{
				Limit:  10,
				Offset: 0,
			},
			mockSvc: func(m *mocks.MockIAnalyticsService) {
				m.EXPECT().
					GetEvents(mock.Anything, mock.AnythingOfType("*analytics.GetEventsRequest")).
					Return(&desc.GetEventsResponse{}, nil).
					Once()
			},
		},
		{
			name: "invalid limit",
			req: &desc.GetEventsRequest{
				Limit:  0,
				Offset: 0,
			},
			wantErr: true,
		},
		{
			name: "invalid offset",
			req: &desc.GetEventsRequest{
				Limit:  1,
				Offset: -1,
			},
			wantErr: true,
		},
		{
			name: "service error",
			req: &desc.GetEventsRequest{
				Limit:  5,
				Offset: 0,
			},
			mockSvc: func(m *mocks.MockIAnalyticsService) {
				m.EXPECT().
					GetEvents(mock.Anything, mock.AnythingOfType("*analytics.GetEventsRequest")).
					Return(&desc.GetEventsResponse{}, errors.New("boom")).
					Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := mocks.NewMockIAnalyticsService(t)
			if tt.mockSvc != nil {
				tt.mockSvc(svc)
			}

			h := handler.NewAnalyticsRouter(svc)

			_, err := h.GetEvents(context.Background(), tt.req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
