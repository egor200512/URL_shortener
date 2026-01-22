package links

// import (
// 	"context"
// 	"errors"
// 	"net/url"
// 	"testing"

// 	s "github.com/egor200512/URL_shortener/services/links/internal/service/links"

// 	mocksR "github.com/egor200512/URL_shortener/services/links/internal/mocks"
// 	"github.com/egor200512/URL_shortener/services/links/models"
// 	mocksC "github.com/egor200512/URL_shortener/shared/mocks"
// 	"github.com/egor200512/URL_shortener/shared/pkg/jwt"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// func TestLinksService_CreateLink(t *testing.T) {
// 	tests := []struct {
// 		name          string
// 		url           string
// 		userID        string
// 		setupMocks    func(*mocksR.MockILinksRepo, *mocksC.MockICache)
// 		expectedError string
// 		checkResult   func(t *testing.T, shortLink string)
// 	}{
// 		{
// 			name:   "success - new link created",
// 			url:    "https://example.com/test",
// 			userID: "user123",
// 			setupMocks: func(mockRepo *mocksR.MockILinksRepo, mockCache *mocksC.MockICache) {
// 				// GetByOriginalLink returns nil (link doesn't exist)
// 				mockRepo.EXPECT().GetByOriginalLink(
// 					mock.Anything,
// 					"example.com/test",
// 				).Return(nil, nil).Once()

// 				// GetByShortLink returns nil (generated short link is unique)
// 				mockRepo.EXPECT().GetByShortLink(
// 					mock.Anything,
// 					mock.AnythingOfType("string"),
// 				).Return(nil, nil).Once()

// 				// InsertLink succeeds
// 				mockRepo.EXPECT().InsertLink(
// 					mock.Anything,
// 					mock.MatchedBy(func(req *models.CreateLinkReq) bool {
// 						return req.UserID == "user123" &&
// 							req.OriginalLink == "example.com/test" &&
// 							req.OriginalLinkHost == "example.com" &&
// 							len(req.ShortLink) == 6
// 					}),
// 				).Return(&models.Link{
// 					ShortLink:    "abc123",
// 					OriginalLink: "example.com/test",
// 				}, nil).Once()
// 				mockCache.EXPECT().SetShort(
// 					mock.Anything,
// 					mock.AnythingOfType("string"),
// 					mock.AnythingOfType("*models.Link"),
// 				).Return(nil).Once()
// 			},
// 			expectedError: "",
// 			checkResult: func(t *testing.T, shortLink string) {
// 				assert.Len(t, shortLink, 6)
// 				assert.Regexp(t, "^[a-zA-Z0-9]{6}$", shortLink)
// 			},
// 		},
// 		{
// 			name:   "error - original link already exists",
// 			url:    "https://example.com/existing",
// 			userID: "user123",
// 			setupMocks: func(mockRepo *mocksR.MockILinksRepo, mockCache *mocksC.MockICache) {
// 				existingLink := &models.Link{
// 					ShortLink:    "abc123",
// 					OriginalLink: "example.com/existing",
// 				}
// 				mockRepo.EXPECT().GetByOriginalLink(
// 					mock.Anything,
// 					"example.com/existing",
// 				).Return(existingLink, nil).Once()
// 				mockCache.EXPECT().SetShort(mock.Anything, mock.Anything, mock.Anything).Maybe().Return(nil)
// 			},
// 			expectedError: "example.com/existing alredy exists",
// 			checkResult:   nil,
// 		},
// 		{
// 			name:   "error - GetByOriginalLink fails",
// 			url:    "https://example.com/test",
// 			userID: "user123",
// 			setupMocks: func(mockRepo *mocksR.MockILinksRepo, mockCache *mocksC.MockICache) {
// 				mockRepo.EXPECT().GetByOriginalLink(
// 					mock.Anything,
// 					"example.com/test",
// 				).Return(nil, errors.New("database error")).Once()
// 				mockCache.EXPECT().SetShort(mock.Anything, mock.Anything, mock.Anything).Maybe().Return(nil)
// 			},
// 			expectedError: "database error",
// 			checkResult:   nil,
// 		},
// 		{
// 			name:   "error - GetByShortLink fails during collision check",
// 			url:    "https://example.com/test",
// 			userID: "user123",
// 			setupMocks: func(mockRepo *mocksR.MockILinksRepo, mockCache *mocksC.MockICache) {
// 				mockRepo.EXPECT().GetByOriginalLink(
// 					mock.Anything,
// 					"example.com/test",
// 				).Return(nil, nil).Once()

// 				mockRepo.EXPECT().GetByShortLink(
// 					mock.Anything,
// 					mock.AnythingOfType("string"),
// 				).Return(nil, errors.New("database error")).Once()
// 				mockCache.EXPECT().SetShort(mock.Anything, mock.Anything, mock.Anything).Maybe().Return(nil)
// 			},
// 			expectedError: "database error",
// 			checkResult:   nil,
// 		},
// 		{
// 			name:   "success - handles short link collision",
// 			url:    "https://example.com/test",
// 			userID: "user123",
// 			setupMocks: func(mockRepo *mocksR.MockILinksRepo, mockCache *mocksC.MockICache) {
// 				mockRepo.EXPECT().GetByOriginalLink(
// 					mock.Anything,
// 					"example.com/test",
// 				).Return(nil, nil).Once()

// 				// First generated short link exists (collision)
// 				existingLink := &models.Link{
// 					ShortLink: "abc123",
// 				}
// 				mockRepo.EXPECT().GetByShortLink(
// 					mock.Anything,
// 					mock.AnythingOfType("string"),
// 				).Return(existingLink, nil).Once()

// 				// Second generated short link is unique
// 				mockRepo.EXPECT().GetByShortLink(
// 					mock.Anything,
// 					mock.AnythingOfType("string"),
// 				).Return(nil, nil).Once()

// 				mockRepo.EXPECT().InsertLink(
// 					mock.Anything,
// 					mock.MatchedBy(func(req *models.CreateLinkReq) bool {
// 						return req.UserID == "user123"
// 					}),
// 				).Return(&models.Link{ShortLink: "def456"}, nil).Once()
// 				mockCache.EXPECT().SetShort(
// 					mock.Anything,
// 					mock.AnythingOfType("string"),
// 					mock.AnythingOfType("*models.Link"),
// 				).Return(nil).Once()
// 			},
// 			expectedError: "",
// 			checkResult: func(t *testing.T, shortLink string) {
// 				assert.Len(t, shortLink, 6)
// 			},
// 		},
// 		{
// 			name:   "error - InsertLink fails",
// 			url:    "https://example.com/test",
// 			userID: "user123",
// 			setupMocks: func(mockRepo *mocksR.MockILinksRepo, mockCache *mocksC.MockICache) {
// 				mockRepo.EXPECT().GetByOriginalLink(
// 					mock.Anything,
// 					"example.com/test",
// 				).Return(nil, nil).Once()

// 				mockRepo.EXPECT().GetByShortLink(
// 					mock.Anything,
// 					mock.AnythingOfType("string"),
// 				).Return(nil, nil).Once()

// 				mockRepo.EXPECT().InsertLink(
// 					mock.Anything,
// 					mock.AnythingOfType("*models.CreateLinkReq"),
// 				).Return(nil, errors.New("insert failed")).Once()
// 				mockCache.EXPECT().SetShort(mock.Anything, mock.Anything, mock.Anything).Maybe().Return(nil)
// 			},
// 			expectedError: "insert failed",
// 			checkResult:   nil,
// 		},
// 		{
// 			name:   "error - missing userID in context",
// 			url:    "https://example.com/test",
// 			userID: "", // Empty userID to trigger context error
// 			setupMocks: func(mockRepo *mocksR.MockILinksRepo, mockCache *mocksC.MockICache) {
// 				mockRepo.EXPECT().GetByOriginalLink(
// 					mock.Anything,
// 					"example.com/test",
// 				).Return(nil, nil).Once()

// 				mockRepo.EXPECT().GetByShortLink(
// 					mock.Anything,
// 					mock.AnythingOfType("string"),
// 				).Return(nil, nil).Once()
// 				mockCache.EXPECT().SetShort(mock.Anything, mock.Anything, mock.Anything).Maybe().Return(nil)
// 			},
// 			expectedError: "failed to get userID from ctx",
// 			checkResult:   nil,
// 		},
// 		{
// 			name:   "success - URL with query parameters (ignored)",
// 			url:    "https://example.com/test?param=value",
// 			userID: "user123",
// 			setupMocks: func(mockRepo *mocksR.MockILinksRepo, mockCache *mocksC.MockICache) {
// 				mockRepo.EXPECT().GetByOriginalLink(
// 					mock.Anything,
// 					"example.com/test",
// 				).Return(nil, nil).Once()

// 				mockRepo.EXPECT().GetByShortLink(
// 					mock.Anything,
// 					mock.AnythingOfType("string"),
// 				).Return(nil, nil).Once()

// 				mockRepo.EXPECT().InsertLink(
// 					mock.Anything,
// 					mock.MatchedBy(func(req *models.CreateLinkReq) bool {
// 						return req.OriginalLink == "example.com/test"
// 					}),
// 				).Return(&models.Link{ShortLink: "ghi789", OriginalLink: "example.com/test"}, nil).Once()
// 				mockCache.EXPECT().SetShort(
// 					mock.Anything,
// 					mock.AnythingOfType("string"),
// 					mock.AnythingOfType("*models.Link"),
// 				).Return(nil).Once()
// 			},
// 			expectedError: "",
// 			checkResult: func(t *testing.T, shortLink string) {
// 				assert.Len(t, shortLink, 6)
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Setup
// 			mockRepo := mocksR.NewMockILinksRepo(t)
// 			mockCache := mocksC.NewMockICache(t)
// 			mockJWTConfig := mocksC.NewMockIJwtConf(t)
// 			tt.setupMocks(mockRepo, mockCache)

// 			service := s.NewLinksService(mockRepo, mockCache, mockJWTConfig)

// 			// Create context with userID
// 			ctx := context.Background()
// 			if tt.userID != "" {
// 				ctx = context.WithValue(ctx, jwt.ClaimsCtxKey, tt.userID)
// 			}

// 			// Parse URL
// 			parsedURL, err := url.Parse(tt.url)
// 			assert.NoError(t, err)

// 			// Execute
// 			shortLink, err := service.CreateLink(ctx, parsedURL)

// 			// Assert
// 			if tt.expectedError != "" {
// 				assert.Error(t, err)
// 				assert.Contains(t, err.Error(), tt.expectedError)
// 				assert.Empty(t, shortLink)
// 			} else {
// 				assert.NoError(t, err)
// 				assert.NotEmpty(t, shortLink)
// 				if tt.checkResult != nil {
// 					tt.checkResult(t, shortLink)
// 				}
// 			}
// 		})
// 	}
// }

// func TestLinksService_CreateLink_ContextValues(t *testing.T) {
// 	t.Run("userID wrong type in context", func(t *testing.T) {
// 		mockRepo := mocksR.NewMockILinksRepo(t)
// 		mockJWTConfig := mocksC.NewMockIJwtConf(t)

// 		mockRepo.EXPECT().GetByOriginalLink(
// 			mock.Anything,
// 			"example.com/test",
// 		).Return(nil, nil).Once()

// 		mockRepo.EXPECT().GetByShortLink(
// 			mock.Anything,
// 			mock.AnythingOfType("string"),
// 		).Return(nil, nil).Once()

// 		mockCache := mocksC.NewMockICache(t)
// 		service := s.NewLinksService(mockRepo, mockCache, mockJWTConfig)

// 		// Set context value with wrong type (int instead of string)
// 		ctx := context.WithValue(context.Background(), jwt.ClaimsCtxKey, 12345)

// 		parsedURL, err := url.Parse("https://example.com/test")
// 		assert.NoError(t, err)

// 		shortLink, err := service.CreateLink(ctx, parsedURL)

// 		assert.Error(t, err)
// 		assert.Contains(t, err.Error(), "failed to get userID from ctx")
// 		assert.Empty(t, shortLink)
// 	})
// }
