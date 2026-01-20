package tests

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	handler "github.com/egor200512/URL_shortener/services/links/internal/handler"
	mocks "github.com/egor200512/URL_shortener/services/links/internal/mocks"
	"github.com/egor200512/URL_shortener/services/links/models"
	desc "github.com/egor200512/URL_shortener/shared/gen/links"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestLinksHandler_GetLinkInfo(t *testing.T) {
	now := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	linkID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name        string
		req         *desc.GetLinkInfoRequest
		mockService func(*mocks.MockILinksService)
		expResp     *desc.LinkInfo
		expErr      bool
		expErrCode  codes.Code
	}{
		{
			name: "success - returns link info",
			req:  &desc.GetLinkInfoRequest{ShortLink: "abc123"},
			mockService: func(m *mocks.MockILinksService) {
				m.EXPECT().
					GetLinkInfo(mock.Anything, "abc123").
					Return(&models.Link{
						ID:               linkID,
						UserID:           userID,
						ShortLink:        "abc123",
						OriginalLinkHost: "example.com",
						OriginalLink:     "example.com/path",
						CreatedAt:        sql.NullTime{Time: now, Valid: true},
					}, nil).
					Once()
			},
			expResp: &desc.LinkInfo{
				Id:               linkID.String(),
				UserId:           userID.String(),
				ShortLink:        "abc123",
				OriginalLinkHost: "example.com",
				OriginalLink:     "example.com/path",
			},
			expErr: false,
		},
		{
			name: "error - empty short link",
			req:  &desc.GetLinkInfoRequest{ShortLink: ""},
			mockService: func(m *mocks.MockILinksService) {
				// validation happens before service call
			},
			expErr:     true,
			expErrCode: codes.InvalidArgument,
		},
		{
			name: "error - service returns error",
			req:  &desc.GetLinkInfoRequest{ShortLink: "abc123"},
			mockService: func(m *mocks.MockILinksService) {
				m.EXPECT().
					GetLinkInfo(mock.Anything, "abc123").
					Return(nil, errors.New("db error")).
					Once()
			},
			expErr:     true,
			expErrCode: codes.Internal,
		},
		{
			name: "error - link not found",
			req:  &desc.GetLinkInfoRequest{ShortLink: "missing"},
			mockService: func(m *mocks.MockILinksService) {
				m.EXPECT().
					GetLinkInfo(mock.Anything, "missing").
					Return(nil, nil).
					Once()
			},
			expErr:     true,
			expErrCode: codes.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svcMock := mocks.NewMockILinksService(t)
			tt.mockService(svcMock)

			h := handler.NewLinksRouter(svcMock)

			resp, err := h.GetLinkInfo(context.Background(), tt.req)

			if tt.expErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
				st, ok := status.FromError(err)
				assert.True(t, ok, "error should be a gRPC status error")
				assert.Equal(t, tt.expErrCode, st.Code())
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, tt.expResp.Id, resp.Id)
			assert.Equal(t, tt.expResp.UserId, resp.UserId)
			assert.Equal(t, tt.expResp.ShortLink, resp.ShortLink)
			assert.Equal(t, tt.expResp.OriginalLinkHost, resp.OriginalLinkHost)
			assert.Equal(t, tt.expResp.OriginalLink, resp.OriginalLink)
			assert.True(t, resp.CreatedAt.AsTime().Equal(now))
		})
	}
}
