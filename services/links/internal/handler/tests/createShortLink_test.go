package tests

import (
	"context"
	"net/url"
	"testing"

	handler "github.com/egor200512/URL_shortener/services/links/internal/handler"
	mocks "github.com/egor200512/URL_shortener/services/links/internal/mocks"
	desc "github.com/egor200512/URL_shortener/shared/gen/links"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestLinksHandler_CreateLink(t *testing.T) {
	tests := []struct {
		name         string
		req          *desc.CreateLinkRequest
		mockService  func(*mocks.MockILinksService)
		expShortLink string
		expErr       bool
		expErrCode   codes.Code
	}{
		{
			name: "Success - valid URL",
			req:  &desc.CreateLinkRequest{OriginalLink: "https://example.com/valid"},
			mockService: func(m *mocks.MockILinksService) {
				u, _ := url.Parse("https://example.com/valid")
				m.EXPECT().
					CreateLink(mock.Anything, u).
					Return("abc123", nil).
					Once()
			},
			expShortLink: "abc123",
			expErr:       false,
		},
		{
			name: "Invalid URL format",
			req:  &desc.CreateLinkRequest{OriginalLink: "invalid-url"},
			mockService: func(m *mocks.MockILinksService) {
				// НЕ вызывается, т.к. падает на url.ParseRequestURI
			},
			expShortLink: "",
			expErr:       true,
			expErrCode:   codes.InvalidArgument,
		},
		{
			name: "Service error",
			req:  &desc.CreateLinkRequest{OriginalLink: "https://example.com/error"},
			mockService: func(m *mocks.MockILinksService) {
				u, _ := url.Parse("https://example.com/error")
				m.EXPECT().
					CreateLink(mock.Anything, u).
					Return("", assert.AnError).
					Once()
			},
			expShortLink: "",
			expErr:       true,
			expErrCode:   codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svcMock := mocks.NewMockILinksService(t)
			tt.mockService(svcMock)

			h := handler.NewLinksRouter(svcMock)

			resp, err := h.CreateLink(context.Background(), tt.req)

			if tt.expErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
				st, ok := status.FromError(err)
				if ok {
					assert.Equal(t, tt.expErrCode, st.Code())
				}
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, tt.expShortLink, resp.ShortLink)
		})
	}
}
