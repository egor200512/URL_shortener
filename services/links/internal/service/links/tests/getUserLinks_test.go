package links

// import (
// 	"context"
// 	"errors"
// 	"testing"

// 	mocksR "github.com/egor200512/URL_shortener/services/links/internal/mocks"
// 	s "github.com/egor200512/URL_shortener/services/links/internal/service/links"
// 	mocksC "github.com/egor200512/URL_shortener/shared/mocks"
// 	"github.com/egor200512/URL_shortener/shared/pkg/jwt"
// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// func TestLinksService_GetUserLinks(t *testing.T) {
// 	userID := uuid.New().String()

// 	tests := []struct {
// 		name          string
// 		userID        string
// 		limit         int32
// 		offset        int32
// 		setupMocks    func(*mocksR.MockILinksRepo, *mocksC.MockICache)
// 		expectedCount int
// 		expectedTotal int32
// 		expectedError string
// 	}{
// 		{
// 			name:   "success - returns links",
// 			userID: userID,
// 			limit:  10,
// 			offset: 0,
// 			setupMocks: func(mockRepo *mocksR.MockILinksRepo, _ *mocksC.MockICache) {
// 				mockRepo.EXPECT().GetUserLinks(
// 					mock.Anything,
// 					userID,
// 					int32(10),
// 					int32(0),
// 				).Return([]string{"example.com/path"}, int32(1), nil).Once()
// 			},
// 			expectedCount: 1,
// 			expectedTotal: 1,
// 		},
// 		{
// 			name:          "error - missing userID",
// 			userID:        "",
// 			limit:         10,
// 			offset:        0,
// 			setupMocks:    func(*mocksR.MockILinksRepo, *mocksC.MockICache) {},
// 			expectedError: "failed to get userID from ctx",
// 		},
// 		{
// 			name:   "error - repo failure",
// 			userID: userID,
// 			limit:  5,
// 			offset: 5,
// 			setupMocks: func(mockRepo *mocksR.MockILinksRepo, _ *mocksC.MockICache) {
// 				mockRepo.EXPECT().GetUserLinks(
// 					mock.Anything,
// 					userID,
// 					int32(5),
// 					int32(5),
// 				).Return(nil, int32(0), errors.New("db error")).Once()
// 			},
// 			expectedError: "db error",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mockRepo := mocksR.NewMockILinksRepo(t)
// 			mockCache := mocksC.NewMockICache(t)
// 			mockJWTConfig := mocksC.NewMockIJwtConf(t)
// 			tt.setupMocks(mockRepo, mockCache)

// 			service := s.NewLinksService(mockRepo, mockCache, mockJWTConfig)

// 			ctx := context.Background()
// 			if tt.userID != "" {
// 				ctx = context.WithValue(ctx, jwt.ClaimsCtxKey, tt.userID)
// 			}

// 			links, total, err := service.GetUserLinks(ctx, tt.limit, tt.offset)

// 			if tt.expectedError != "" {
// 				assert.Error(t, err)
// 				assert.Contains(t, err.Error(), tt.expectedError)
// 				assert.Nil(t, links)
// 				assert.Equal(t, int32(0), total)
// 				return
// 			}

// 			assert.NoError(t, err)
// 			assert.Len(t, links, tt.expectedCount)
// 			assert.Equal(t, tt.expectedTotal, total)
// 		})
// 	}
// }
