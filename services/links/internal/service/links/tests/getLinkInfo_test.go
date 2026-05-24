package links

import (
	"context"
	"errors"
	"testing"

	mocksR "github.com/egor200512/URL_shortener/services/links/internal/mocks"
	service "github.com/egor200512/URL_shortener/services/links/internal/service/links"
	"github.com/egor200512/URL_shortener/services/links/models"
	mocksC "github.com/egor200512/URL_shortener/shared/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLinksService_GetLinkInfo(t *testing.T) {
	tests := []struct {
		name          string
		shortLink     string
		setupMocks    func(*mocksR.MockILinksRepo, *mocksC.MockICache)
		expectedLink  *models.Link
		expectedError string
	}{
		{
			name:      "success - cache hit",
			shortLink: "abc123",
			setupMocks: func(_ *mocksR.MockILinksRepo, cache *mocksC.MockICache) {
				cache.EXPECT().GetShort(mock.Anything, "abc123").
					Return(&models.Link{ShortLink: "abc123", OriginalLink: "example.com/test"}, nil).Once()
			},
			expectedLink: &models.Link{ShortLink: "abc123", OriginalLink: "example.com/test"},
		},
		{
			name:      "success - cache miss repo hit",
			shortLink: "abc123",
			setupMocks: func(repo *mocksR.MockILinksRepo, cache *mocksC.MockICache) {
				cache.EXPECT().GetShort(mock.Anything, "abc123").
					Return(nil, nil).Once()
				repo.EXPECT().GetByShortLink(mock.Anything, "abc123").
					Return(&models.Link{ShortLink: "abc123", OriginalLink: "example.com/test"}, nil).Once()
				cache.EXPECT().SetShort(mock.Anything, "abc123", mock.MatchedBy(func(link *models.Link) bool {
					return link.ShortLink == "abc123" && link.OriginalLink == "example.com/test"
				})).Return(nil).Once()
			},
			expectedLink: &models.Link{ShortLink: "abc123", OriginalLink: "example.com/test"},
		},
		{
			name:      "success - not found anywhere",
			shortLink: "missing",
			setupMocks: func(repo *mocksR.MockILinksRepo, cache *mocksC.MockICache) {
				cache.EXPECT().GetShort(mock.Anything, "missing").
					Return(nil, nil).Once()
				repo.EXPECT().GetByShortLink(mock.Anything, "missing").
					Return(nil, nil).Once()
			},
			expectedLink: nil,
		},
		{
			name:      "cache get failure is ignored",
			shortLink: "abc123",
			setupMocks: func(repo *mocksR.MockILinksRepo, cache *mocksC.MockICache) {
				cache.EXPECT().GetShort(mock.Anything, "abc123").
					Return(nil, errors.New("cache error")).Once()
				repo.EXPECT().GetByShortLink(mock.Anything, "abc123").
					Return(&models.Link{ShortLink: "abc123", OriginalLink: "example.com/test"}, nil).Once()
				cache.EXPECT().SetShort(mock.Anything, "abc123", mock.AnythingOfType("*models.Link")).
					Return(nil).Once()
			},
			expectedLink: &models.Link{ShortLink: "abc123", OriginalLink: "example.com/test"},
		},
		{
			name:      "error - repo failure",
			shortLink: "abc123",
			setupMocks: func(repo *mocksR.MockILinksRepo, cache *mocksC.MockICache) {
				cache.EXPECT().GetShort(mock.Anything, "abc123").
					Return(nil, nil).Once()
				repo.EXPECT().GetByShortLink(mock.Anything, "abc123").
					Return(nil, errors.New("db error")).Once()
			},
			expectedError: "db error",
		},
		{
			name:      "cache set failure after repo hit is ignored",
			shortLink: "abc123",
			setupMocks: func(repo *mocksR.MockILinksRepo, cache *mocksC.MockICache) {
				cache.EXPECT().GetShort(mock.Anything, "abc123").
					Return(nil, nil).Once()
				repo.EXPECT().GetByShortLink(mock.Anything, "abc123").
					Return(&models.Link{ShortLink: "abc123", OriginalLink: "example.com/test"}, nil).Once()
				cache.EXPECT().SetShort(mock.Anything, "abc123", mock.AnythingOfType("*models.Link")).
					Return(errors.New("cache set failed")).Once()
			},
			expectedLink: &models.Link{ShortLink: "abc123", OriginalLink: "example.com/test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocksR.NewMockILinksRepo(t)
			cache := mocksC.NewMockICache(t)
			jwtConf := mocksC.NewMockIJwtConf(t)
			tt.setupMocks(repo, cache)

			svc := service.NewLinksService(repo, cache, &fakeProducer{}, jwtConf)

			link, err := svc.GetLinkInfo(context.Background(), tt.shortLink)
			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, link)
				return
			}

			require.NoError(t, err)
			if tt.expectedLink == nil {
				assert.Nil(t, link)
				return
			}

			require.NotNil(t, link)
			assert.Equal(t, tt.expectedLink.ShortLink, link.ShortLink)
			assert.Equal(t, tt.expectedLink.OriginalLink, link.OriginalLink)
		})
	}
}
