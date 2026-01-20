package tests

import (
	"context"
	"errors"
	"testing"

	handler "github.com/egor200512/URL_shortener/services/links/internal/handler"
	mocks "github.com/egor200512/URL_shortener/services/links/internal/mocks"
	"github.com/egor200512/URL_shortener/services/links/models"
	desc "github.com/egor200512/URL_shortener/shared/gen/links"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestLinksHandler_GetOriginalLink(t *testing.T) {
	tests := []struct {
		name        string
		req         *desc.GetOriginalLinkRequest
		mockService func(*mocks.MockILinksService)
		expResp     *desc.GetOriginalLinkResponse
		expErr      bool
		expErrCode  codes.Code
	}{
		{
			name: "success - returns original link",
			req:  &desc.GetOriginalLinkRequest{ShortLink: "abc123"},
			mockService: func(m *mocks.MockILinksService) {
				m.EXPECT().
					GetLinkInfo(mock.Anything, "abc123").
					Return(&models.Link{
						ShortLink:    "abc123",
						OriginalLink: "example.com/path",
					}, nil).
					Once()
			},
			expResp: &desc.GetOriginalLinkResponse{OriginalUrl: "example.com/path"},
			expErr:  false,
		},
		{
			name: "error - empty short link",
			req:  &desc.GetOriginalLinkRequest{ShortLink: ""},
			mockService: func(m *mocks.MockILinksService) {
				// validation happens before service call
			},
			expErr:     true,
			expErrCode: codes.InvalidArgument,
		},
		{
			name: "error - service returns error",
			req:  &desc.GetOriginalLinkRequest{ShortLink: "abc123"},
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
			req:  &desc.GetOriginalLinkRequest{ShortLink: "missing"},
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

			resp, err := h.GetOriginalLink(context.Background(), tt.req)

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
			assert.Equal(t, tt.expResp.OriginalUrl, resp.OriginalUrl)
		})
	}
}
