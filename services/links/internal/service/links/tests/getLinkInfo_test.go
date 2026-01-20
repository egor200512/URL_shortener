package links

import (
	"context"
	"errors"
	"testing"

	s "github.com/egor200512/URL_shortener/services/links/internal/service/links"
	mocksR "github.com/egor200512/URL_shortener/services/links/internal/mocks"
	"github.com/egor200512/URL_shortener/services/links/models"
	mocksC "github.com/egor200512/URL_shortener/shared/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLinksService_GetLinkInfo(t *testing.T) {
	tests := []struct {
		name          string
		shortLink     string
		setupMocks    func(*mocksR.MockILinksRepo)
		expectedLink  *models.Link
		expectedError string
	}{
		{
			name:      "success - link found",
			shortLink: "abc123",
			setupMocks: func(mockRepo *mocksR.MockILinksRepo) {
				mockRepo.EXPECT().GetByShortLink(
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
			name:      "success - link not found",
			shortLink: "missing",
			setupMocks: func(mockRepo *mocksR.MockILinksRepo) {
				mockRepo.EXPECT().GetByShortLink(
					mock.Anything,
					"missing",
				).Return(nil, nil).Once()
			},
			expectedLink:  nil,
			expectedError: "",
		},
		{
			name:      "error - repo failure",
			shortLink: "abc123",
			setupMocks: func(mockRepo *mocksR.MockILinksRepo) {
				mockRepo.EXPECT().GetByShortLink(
					mock.Anything,
					"abc123",
				).Return(nil, errors.New("db error")).Once()
			},
			expectedLink:  nil,
			expectedError: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mocksR.NewMockILinksRepo(t)
			mockJWTConfig := mocksC.NewMockIJwtConf(t)
			tt.setupMocks(mockRepo)

			service := s.NewLinksService(mockRepo, mockJWTConfig)

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
