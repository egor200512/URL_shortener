package tests

import (
	"context"
	"errors"
	"testing"

	handler "github.com/egor200512/URL_shortener/services/links/internal/handler"
	mocks "github.com/egor200512/URL_shortener/services/links/internal/mocks"
	desc "github.com/egor200512/URL_shortener/shared/gen/links"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestLinksHandler_GetUserLinks(t *testing.T) {
	tests := []struct {
		name        string
		req         *desc.GetUserLinksRequest
		mockService func(*mocks.MockILinksService)
		expCount    int
		expTotal    int32
		expErr      bool
		expErrCode  codes.Code
	}{
		{
			name: "success - returns links",
			req:  &desc.GetUserLinksRequest{Limit: 10, Offset: 0},
			mockService: func(m *mocks.MockILinksService) {
				m.EXPECT().
					GetUserLinks(mock.Anything, int32(10), int32(0)).
					Return([]string{"example.com/path"}, int32(1), nil).
					Once()
			},
			expCount: 1,
			expTotal: 1,
			expErr:   false,
		},
		{
			name: "error - invalid limit",
			req:  &desc.GetUserLinksRequest{Limit: 0, Offset: 0},
			mockService: func(m *mocks.MockILinksService) {
				// validation happens before service call
			},
			expErr:     true,
			expErrCode: codes.InvalidArgument,
		},
		{
			name: "error - invalid offset",
			req:  &desc.GetUserLinksRequest{Limit: 10, Offset: -1},
			mockService: func(m *mocks.MockILinksService) {
				// validation happens before service call
			},
			expErr:     true,
			expErrCode: codes.InvalidArgument,
		},
		{
			name: "error - service returns error",
			req:  &desc.GetUserLinksRequest{Limit: 10, Offset: 0},
			mockService: func(m *mocks.MockILinksService) {
				m.EXPECT().
					GetUserLinks(mock.Anything, int32(10), int32(0)).
					Return(nil, int32(0), errors.New("db error")).
					Once()
			},
			expErr:     true,
			expErrCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svcMock := mocks.NewMockILinksService(t)
			tt.mockService(svcMock)

			h := handler.NewLinksRouter(svcMock)

			resp, err := h.GetUserLinks(context.Background(), tt.req)

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
			assert.Len(t, resp.Links, tt.expCount)
			assert.Equal(t, tt.expTotal, resp.TotalCount)
			assert.Equal(t, "example.com/path", resp.Links[0])
		})
	}
}
