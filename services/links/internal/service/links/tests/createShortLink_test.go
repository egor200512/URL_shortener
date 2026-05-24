package links

import (
	"context"
	"errors"
	"net/url"
	"testing"

	mocksR "github.com/egor200512/URL_shortener/services/links/internal/mocks"
	service "github.com/egor200512/URL_shortener/services/links/internal/service/links"
	"github.com/egor200512/URL_shortener/services/links/models"
	mocksC "github.com/egor200512/URL_shortener/shared/mocks"
	"github.com/egor200512/URL_shortener/shared/pkg/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLinksService_CreateLink(t *testing.T) {
	userID := uuid.New().String()

	tests := []struct {
		name          string
		rawURL        string
		userID        any
		setupMocks    func(*mocksR.MockILinksRepo, *mocksC.MockICache)
		expectedError string
	}{
		{
			name:   "success - new link created",
			rawURL: "https://example.com/test",
			userID: userID,
			setupMocks: func(repo *mocksR.MockILinksRepo, cache *mocksC.MockICache) {
				repo.EXPECT().GetByOriginalLink(mock.Anything, "example.com/test").
					Return(nil, nil).Once()
				repo.EXPECT().GetByShortLink(mock.Anything, mock.AnythingOfType("string")).
					Return(nil, nil).Once()
				repo.EXPECT().InsertLink(mock.Anything, mock.MatchedBy(func(req *models.CreateLinkReq) bool {
					return req.UserID == userID &&
						req.OriginalLink == "example.com/test" &&
						req.OriginalLinkHost == "example.com" &&
						len(req.ShortLink) == 6
				})).Return(&models.Link{
					UserID:           uuid.MustParse(userID),
					ShortLink:        "abc123",
					OriginalLinkHost: "example.com",
					OriginalLink:     "example.com/test",
				}, nil).Once()
				cache.EXPECT().SetShort(mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*models.Link")).
					Return(nil).Once()
			},
		},
		{
			name:   "error - original link already exists",
			rawURL: "https://example.com/existing",
			userID: userID,
			setupMocks: func(repo *mocksR.MockILinksRepo, _ *mocksC.MockICache) {
				repo.EXPECT().GetByOriginalLink(mock.Anything, "example.com/existing").
					Return(&models.Link{ShortLink: "abc123", OriginalLink: "example.com/existing"}, nil).Once()
			},
			expectedError: "example.com/existing alredy exists",
		},
		{
			name:   "error - get by original link fails",
			rawURL: "https://example.com/test",
			userID: userID,
			setupMocks: func(repo *mocksR.MockILinksRepo, _ *mocksC.MockICache) {
				repo.EXPECT().GetByOriginalLink(mock.Anything, "example.com/test").
					Return(nil, errors.New("database error")).Once()
			},
			expectedError: "database error",
		},
		{
			name:   "error - get by short link fails",
			rawURL: "https://example.com/test",
			userID: userID,
			setupMocks: func(repo *mocksR.MockILinksRepo, _ *mocksC.MockICache) {
				repo.EXPECT().GetByOriginalLink(mock.Anything, "example.com/test").
					Return(nil, nil).Once()
				repo.EXPECT().GetByShortLink(mock.Anything, mock.AnythingOfType("string")).
					Return(nil, errors.New("database error")).Once()
			},
			expectedError: "database error",
		},
		{
			name:   "error - missing user id in context",
			rawURL: "https://example.com/test",
			userID: nil,
			setupMocks: func(repo *mocksR.MockILinksRepo, _ *mocksC.MockICache) {
				repo.EXPECT().GetByOriginalLink(mock.Anything, "example.com/test").
					Return(nil, nil).Once()
				repo.EXPECT().GetByShortLink(mock.Anything, mock.AnythingOfType("string")).
					Return(nil, nil).Once()
			},
			expectedError: "failed to get userID from ctx",
		},
		{
			name:   "error - wrong user id type in context",
			rawURL: "https://example.com/test",
			userID: 12345,
			setupMocks: func(repo *mocksR.MockILinksRepo, _ *mocksC.MockICache) {
				repo.EXPECT().GetByOriginalLink(mock.Anything, "example.com/test").
					Return(nil, nil).Once()
				repo.EXPECT().GetByShortLink(mock.Anything, mock.AnythingOfType("string")).
					Return(nil, nil).Once()
			},
			expectedError: "failed to get userID from ctx",
		},
		{
			name:   "error - insert link fails",
			rawURL: "https://example.com/test",
			userID: userID,
			setupMocks: func(repo *mocksR.MockILinksRepo, _ *mocksC.MockICache) {
				repo.EXPECT().GetByOriginalLink(mock.Anything, "example.com/test").
					Return(nil, nil).Once()
				repo.EXPECT().GetByShortLink(mock.Anything, mock.AnythingOfType("string")).
					Return(nil, nil).Once()
				repo.EXPECT().InsertLink(mock.Anything, mock.AnythingOfType("*models.CreateLinkReq")).
					Return(nil, errors.New("insert failed")).Once()
			},
			expectedError: "insert failed",
		},
		{
			name:   "success - url query is ignored",
			rawURL: "https://example.com/test?param=value",
			userID: userID,
			setupMocks: func(repo *mocksR.MockILinksRepo, cache *mocksC.MockICache) {
				repo.EXPECT().GetByOriginalLink(mock.Anything, "example.com/test").
					Return(nil, nil).Once()
				repo.EXPECT().GetByShortLink(mock.Anything, mock.AnythingOfType("string")).
					Return(nil, nil).Once()
				repo.EXPECT().InsertLink(mock.Anything, mock.MatchedBy(func(req *models.CreateLinkReq) bool {
					return req.OriginalLink == "example.com/test"
				})).Return(&models.Link{
					UserID:       uuid.MustParse(userID),
					ShortLink:    "ghi789",
					OriginalLink: "example.com/test",
				}, nil).Once()
				cache.EXPECT().SetShort(mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*models.Link")).
					Return(nil).Once()
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

			parsedURL, err := url.Parse(tt.rawURL)
			require.NoError(t, err)

			shortLink, err := svc.CreateLink(ctx, parsedURL)
			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Empty(t, shortLink)
				return
			}

			require.NoError(t, err)
			assert.Len(t, shortLink, 6)
			assert.Regexp(t, "^[a-zA-Z0-9]{6}$", shortLink)
		})
	}
}
