package links

import (
	"context"
	"errors"
	"testing"

	mocksR "github.com/egor200512/URL_shortener/services/links/internal/mocks"
	s "github.com/egor200512/URL_shortener/services/links/internal/service/links"
	"github.com/egor200512/URL_shortener/services/links/models"
	mocksC "github.com/egor200512/URL_shortener/shared/mocks"
	"github.com/egor200512/URL_shortener/shared/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLinksService_DeleteLink(t *testing.T) {
	tests := []struct {
		name          string
		shortLink     string
		userID        string
		setupMocks    func(*mocksR.MockILinksRepo, *mocksC.MockICache)
		expectedError string
	}{
		{
			name:      "success - link deleted",
			shortLink: "abc123",
			userID:    "user123",
			setupMocks: func(mockRepo *mocksR.MockILinksRepo, mockCache *mocksC.MockICache) {
				mockRepo.EXPECT().GetByShortLink(
					mock.Anything,
					"abc123",
				).Return(&models.Link{
					ShortLink: "abc123",
				}, nil).Once()

				mockRepo.EXPECT().DeleteLink(
					mock.Anything,
					"abc123",
					"user123",
				).Return(nil).Once()

				mockCache.EXPECT().DelShort(
					mock.Anything,
					"abc123",
				).Return(nil).Once()
			},
			expectedError: "",
		},
		{
			name:      "error - GetByShortLink fails",
			shortLink: "abc123",
			userID:    "user123",
			setupMocks: func(mockRepo *mocksR.MockILinksRepo, mockCache *mocksC.MockICache) {
				mockRepo.EXPECT().GetByShortLink(
					mock.Anything,
					"abc123",
				).Return(nil, errors.New("db error")).Once()
				mockCache.EXPECT().DelShort(mock.Anything, mock.Anything).Maybe().Return(nil)
			},
			expectedError: "db error",
		},
		{
			name:      "error - link not found",
			shortLink: "missing",
			userID:    "user123",
			setupMocks: func(mockRepo *mocksR.MockILinksRepo, mockCache *mocksC.MockICache) {
				mockRepo.EXPECT().GetByShortLink(
					mock.Anything,
					"missing",
				).Return(nil, nil).Once()
				mockCache.EXPECT().DelShort(mock.Anything, mock.Anything).Maybe().Return(nil)
			},
			expectedError: "link for missing doesn't exist",
		},
		{
			name:      "error - missing userID in context",
			shortLink: "abc123",
			userID:    "",
			setupMocks: func(mockRepo *mocksR.MockILinksRepo, mockCache *mocksC.MockICache) {
				mockRepo.EXPECT().GetByShortLink(
					mock.Anything,
					"abc123",
				).Return(&models.Link{ShortLink: "abc123"}, nil).Once()
				mockCache.EXPECT().DelShort(mock.Anything, mock.Anything).Maybe().Return(nil)
			},
			expectedError: "failed to get userID from ctx",
		},
		{
			name:      "error - DeleteLink fails",
			shortLink: "abc123",
			userID:    "user123",
			setupMocks: func(mockRepo *mocksR.MockILinksRepo, mockCache *mocksC.MockICache) {
				mockRepo.EXPECT().GetByShortLink(
					mock.Anything,
					"abc123",
				).Return(&models.Link{ShortLink: "abc123"}, nil).Once()

				mockRepo.EXPECT().DeleteLink(
					mock.Anything,
					"abc123",
					"user123",
				).Return(errors.New("delete failed")).Once()
				mockCache.EXPECT().DelShort(mock.Anything, mock.Anything).Maybe().Return(nil)
			},
			expectedError: "delete failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocksR.NewMockILinksRepo(t)
			mockCache := mocksC.NewMockICache(t)
			mockJWTConfig := mocksC.NewMockIJwtConf(t)
			tt.setupMocks(mockRepo, mockCache)

			service := s.NewLinksService(mockRepo, mockCache, mockJWTConfig)

			ctx := context.Background()
			if tt.userID != "" {
				ctx = context.WithValue(ctx, jwt.ClaimsCtxKey, tt.userID)
			}

			err := service.DeleteLink(ctx, tt.shortLink)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestLinksService_DeleteLink_ContextValues(t *testing.T) {
	t.Run("userID wrong type in context", func(t *testing.T) {
		mockRepo := mocksR.NewMockILinksRepo(t)
		mockCache := mocksC.NewMockICache(t)
		mockJWTConfig := mocksC.NewMockIJwtConf(t)

		mockRepo.EXPECT().GetByShortLink(
			mock.Anything,
			"abc123",
		).Return(&models.Link{ShortLink: "abc123"}, nil).Once()

		service := s.NewLinksService(mockRepo, mockCache, mockJWTConfig)

		ctx := context.WithValue(context.Background(), jwt.ClaimsCtxKey, 12345)

		err := service.DeleteLink(ctx, "abc123")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get userID from ctx")
	})
}
