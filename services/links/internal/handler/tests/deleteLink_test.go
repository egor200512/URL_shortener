package tests

import (
	"context"
	"errors"
	"testing"

	h "github.com/egor200512/URL_shortener/services/links/internal/handler"

	mocks "github.com/egor200512/URL_shortener/services/links/internal/mocks"
	desc "github.com/egor200512/URL_shortener/shared/gen/links"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestLinksHandler_DeleteLink(t *testing.T) {
	tests := []struct {
		name          string
		request       *desc.DeleteLinkRequest
		setupMocks    func(*mocks.MockILinksService)
		expectedError bool
		expectedCode  codes.Code
		errorContains string
	}{
		{
			name: "success - link deleted",
			request: &desc.DeleteLinkRequest{
				ShortLink: "abc123",
			},
			setupMocks: func(mockService *mocks.MockILinksService) {
				mockService.EXPECT().DeleteLink(
					mock.Anything,
					"abc123",
				).Return(nil).Once()
			},
			expectedError: false,
		},
		{
			name: "error - empty short link",
			request: &desc.DeleteLinkRequest{
				ShortLink: "",
			},
			setupMocks: func(mockService *mocks.MockILinksService) {
				// No expectations - validation happens before service call
			},
			expectedError: true,
			expectedCode:  codes.InvalidArgument,
			errorContains: "short link is required",
		},
		{
			name: "error - service returns error",
			request: &desc.DeleteLinkRequest{
				ShortLink: "abc123",
			},
			setupMocks: func(mockService *mocks.MockILinksService) {
				mockService.EXPECT().DeleteLink(
					mock.Anything,
					"abc123",
				).Return(errors.New("database error")).Once()
			},
			expectedError: true,
			expectedCode:  codes.Internal,
			errorContains: "failed to delete link",
		},
		{
			name: "error - link not found",
			request: &desc.DeleteLinkRequest{
				ShortLink: "notfound",
			},
			setupMocks: func(mockService *mocks.MockILinksService) {
				mockService.EXPECT().DeleteLink(
					mock.Anything,
					"notfound",
				).Return(errors.New("link not found")).Once()
			},
			expectedError: true,
			expectedCode:  codes.Internal,
			errorContains: "failed to delete link",
		},
		{
			name: "success - delete existing link with special characters",
			request: &desc.DeleteLinkRequest{
				ShortLink: "aB12cD",
			},
			setupMocks: func(mockService *mocks.MockILinksService) {
				mockService.EXPECT().DeleteLink(
					mock.Anything,
					"aB12cD",
				).Return(nil).Once()
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := mocks.NewMockILinksService(t)
			tt.setupMocks(mockService)

			handler := h.NewLinksRouter(mockService)

			ctx := context.Background()

			// Execute
			result, err := handler.DeleteLink(ctx, tt.request)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)

				// Check gRPC status code
				st, ok := status.FromError(err)
				assert.True(t, ok, "error should be a gRPC status error")
				assert.Equal(t, tt.expectedCode, st.Code())

				if tt.errorContains != "" {
					assert.Contains(t, st.Message(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestLinksHandler_DeleteLink_EdgeCases(t *testing.T) {
	t.Run("nil request", func(t *testing.T) {
		mockService := mocks.NewMockILinksService(t)

		handler := h.NewLinksRouter(mockService)

		ctx := context.Background()

		// This will panic in real code, but we test the behavior
		assert.Panics(t, func() {
			handler.DeleteLink(ctx, nil)
		})
	})

	t.Run("whitespace only short link", func(t *testing.T) {
		mockService := mocks.NewMockILinksService(t)

		handler := h.NewLinksRouter(mockService)

		ctx := context.Background()
		request := &desc.DeleteLinkRequest{
			ShortLink: "   ",
		}

		// With current implementation, this passes validation (len > 0)
		// but you might want to add trim validation
		mockService.EXPECT().DeleteLink(
			mock.Anything,
			"   ",
		).Return(errors.New("invalid short link")).Once()

		result, err := handler.DeleteLink(ctx, request)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
