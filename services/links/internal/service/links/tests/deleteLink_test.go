package links

import (
	"context"
	"errors"
	"testing"

	mocksR "github.com/egor200512/URL_shortener/services/links/internal/mocks"
	service "github.com/egor200512/URL_shortener/services/links/internal/service/links"
	"github.com/egor200512/URL_shortener/services/links/models"
	mocksC "github.com/egor200512/URL_shortener/shared/mocks"
	"github.com/egor200512/URL_shortener/shared/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLinksService_DeleteLink(t *testing.T) {
	tests := []struct {
		name          string
		shortLink     string
		userID        any
		setupMocks    func(*mocksR.MockILinksRepo, *mocksC.MockICache)
		expectedError string
	}{
		{
			name:      "success - link deleted",
			shortLink: "abc123",
			userID:    "user123",
			setupMocks: func(repo *mocksR.MockILinksRepo, cache *mocksC.MockICache) {
				repo.EXPECT().GetByShortLink(mock.Anything, "abc123").
					Return(&models.Link{ShortLink: "abc123"}, nil).Once()
				repo.EXPECT().DeleteLink(mock.Anything, "abc123", "user123").
					Return(nil).Once()
				cache.EXPECT().DelShort(mock.Anything, "abc123").
					Return(nil).Once()
			},
		},
		{
			name:      "error - get by short link fails",
			shortLink: "abc123",
			userID:    "user123",
			setupMocks: func(repo *mocksR.MockILinksRepo, _ *mocksC.MockICache) {
				repo.EXPECT().GetByShortLink(mock.Anything, "abc123").
					Return(nil, errors.New("db error")).Once()
			},
			expectedError: "db error",
		},
		{
			name:      "error - link not found",
			shortLink: "missing",
			userID:    "user123",
			setupMocks: func(repo *mocksR.MockILinksRepo, _ *mocksC.MockICache) {
				repo.EXPECT().GetByShortLink(mock.Anything, "missing").
					Return(nil, nil).Once()
			},
			expectedError: "link for missing doesn't exist",
		},
		{
			name:      "error - missing user id in context",
			shortLink: "abc123",
			userID:    nil,
			setupMocks: func(repo *mocksR.MockILinksRepo, _ *mocksC.MockICache) {
				repo.EXPECT().GetByShortLink(mock.Anything, "abc123").
					Return(&models.Link{ShortLink: "abc123"}, nil).Once()
			},
			expectedError: "failed to get userID from ctx",
		},
		{
			name:      "error - wrong user id type in context",
			shortLink: "abc123",
			userID:    12345,
			setupMocks: func(repo *mocksR.MockILinksRepo, _ *mocksC.MockICache) {
				repo.EXPECT().GetByShortLink(mock.Anything, "abc123").
					Return(&models.Link{ShortLink: "abc123"}, nil).Once()
			},
			expectedError: "failed to get userID from ctx",
		},
		{
			name:      "error - delete link fails",
			shortLink: "abc123",
			userID:    "user123",
			setupMocks: func(repo *mocksR.MockILinksRepo, _ *mocksC.MockICache) {
				repo.EXPECT().GetByShortLink(mock.Anything, "abc123").
					Return(&models.Link{ShortLink: "abc123"}, nil).Once()
				repo.EXPECT().DeleteLink(mock.Anything, "abc123", "user123").
					Return(errors.New("delete failed")).Once()
			},
			expectedError: "delete failed",
		},
		{
			name:      "success - cache delete error is ignored",
			shortLink: "abc123",
			userID:    "user123",
			setupMocks: func(repo *mocksR.MockILinksRepo, cache *mocksC.MockICache) {
				repo.EXPECT().GetByShortLink(mock.Anything, "abc123").
					Return(&models.Link{ShortLink: "abc123"}, nil).Once()
				repo.EXPECT().DeleteLink(mock.Anything, "abc123", "user123").
					Return(nil).Once()
				cache.EXPECT().DelShort(mock.Anything, "abc123").
					Return(errors.New("cache error")).Once()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocksR.NewMockILinksRepo(t)
			cache := mocksC.NewMockICache(t)
			jwtConf := mocksC.NewMockIJwtConf(t)
			tt.setupMocks(repo, cache)

			svc := service.NewLinksService(repo, cache, &fakeProducer{}, jwtConf)

			ctx := context.Background()
			if tt.userID != nil {
				ctx = context.WithValue(ctx, jwt.ClaimsCtxKey, tt.userID)
			}

			err := svc.DeleteLink(ctx, tt.shortLink)
			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				return
			}

			require.NoError(t, err)
		})
	}
}
