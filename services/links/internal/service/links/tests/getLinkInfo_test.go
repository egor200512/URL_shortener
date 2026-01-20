package links

import (
	"context"
	"errors"
	"testing"

	mocksR "github.com/egor200512/URL_shortener/services/links/internal/mocks"
	s "github.com/egor200512/URL_shortener/services/links/internal/service/links"
	"github.com/egor200512/URL_shortener/services/links/models"
	mocksC "github.com/egor200512/URL_shortener/shared/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
			setupMocks: func(mockRepo *mocksR.MockILinksRepo, mockCache *mocksC.MockICache) {
				mockCache.EXPECT().GetShort(
					mock.Anything,
					"abc123",
				).Return(&models.Link{
					ShortLink:    "abc123",
					OriginalLink: "example.com/test",
				}, nil).Once()
			},
			expectedLink: &models.Link{
				ShortLink:    "abc123",
				OriginalLink: "example.com/test",
			},
			expectedError: "",
		},
		{
			name:      "success - cache miss repo hit",
			shortLink: "abc123",
			setupMocks: func(mockRepo *mocksR.MockILinksRepo, mockCache *mocksC.MockICache) {
				mockCache.EXPECT().GetShort(
					mock.Anything,
					"abc123",
				).Return(nil, nil).Once()
				mockRepo.EXPECT().GetByShortLink(
					mock.Anything,
					"abc123",
				).Return(&models.Link{
					ShortLink:    "abc123",
					OriginalLink: "example.com/test",
				}, nil).Once()
				mockCache.EXPECT().SetShort(
					mock.Anything,
					"abc123",
					mock.MatchedBy(func(link *models.Link) bool {
						return link.ShortLink == "abc123" && link.OriginalLink == "example.com/test"
					}),
				).Return(nil).Once()
			},
			expectedLink: &models.Link{
				ShortLink:    "abc123",
				OriginalLink: "example.com/test",
			},
			expectedError: "",
		},
		{
			name:      "success - not found anywhere",
			shortLink: "missing",
			setupMocks: func(mockRepo *mocksR.MockILinksRepo, mockCache *mocksC.MockICache) {
				mockCache.EXPECT().GetShort(
					mock.Anything,
					"missing",
				).Return(nil, nil).Once()
				mockRepo.EXPECT().GetByShortLink(
					mock.Anything,
					"missing",
				).Return(nil, nil).Once()
			},
			expectedLink:  nil,
			expectedError: "",
		},
		{
			name:      "error - cache failure",
			shortLink: "abc123",
			setupMocks: func(mockRepo *mocksR.MockILinksRepo, mockCache *mocksC.MockICache) {
				mockCache.EXPECT().GetShort(
					mock.Anything,
					"abc123",
				).Return(nil, errors.New("cache error")).Once()
			},
			expectedLink:  nil,
			expectedError: "cache error",
		},
		{
			name:      "error - repo failure",
			shortLink: "abc123",
			setupMocks: func(mockRepo *mocksR.MockILinksRepo, mockCache *mocksC.MockICache) {
				mockCache.EXPECT().GetShort(
					mock.Anything,
					"abc123",
				).Return(nil, nil).Once()
				mockRepo.EXPECT().GetByShortLink(
					mock.Anything,
					"abc123",
				).Return(nil, errors.New("db error")).Once()
			},
			expectedLink:  nil,
			expectedError: "db error",
		},
		{
			name:      "cache set failure after repo hit is ignored",
			shortLink: "abc123",
			setupMocks: func(mockRepo *mocksR.MockILinksRepo, mockCache *mocksC.MockICache) {
				mockCache.EXPECT().GetShort(
					mock.Anything,
					"abc123",
				).Return(nil, nil).Once()
				mockRepo.EXPECT().GetByShortLink(
					mock.Anything,
					"abc123",
				).Return(&models.Link{
					ShortLink:    "abc123",
					OriginalLink: "example.com/test",
				}, nil).Once()
				mockCache.EXPECT().SetShort(
					mock.Anything,
					"abc123",
					mock.AnythingOfType("*models.Link"),
				).Return(errors.New("cache set failed")).Once()
			},
			expectedLink: &models.Link{
				ShortLink:    "abc123",
				OriginalLink: "example.com/test",
			},
			expectedError: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocksR.NewMockILinksRepo(t)
			mockCache := mocksC.NewMockICache(t)
			mockJWTConfig := mocksC.NewMockIJwtConf(t)
			tt.setupMocks(mockRepo, mockCache)

			service := s.NewLinksService(mockRepo, mockCache, mockJWTConfig)

			link, err := service.GetLinkInfo(context.Background(), tt.shortLink)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, link)
				return
			}

			assert.NoError(t, err)
			if tt.expectedLink == nil {
				assert.Nil(t, link)
				return
			}

			assert.NotNil(t, link)
			assert.Equal(t, tt.expectedLink.ShortLink, link.ShortLink)
			assert.Equal(t, tt.expectedLink.OriginalLink, link.OriginalLink)
		})
	}
}
