package links

import (
	"context"
	"errors"
	"testing"

	mocksR "github.com/egor200512/URL_shortener/services/links/internal/mocks"
	service "github.com/egor200512/URL_shortener/services/links/internal/service/links"
	mocksC "github.com/egor200512/URL_shortener/shared/mocks"
	"github.com/egor200512/URL_shortener/shared/pkg/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLinksService_GetUserLinks(t *testing.T) {
	userID := uuid.New().String()

	tests := []struct {
		name          string
		userID        any
		limit         int32
		offset        int32
		setupMocks    func(*mocksR.MockILinksRepo)
		expectedLinks []string
		expectedTotal int32
		expectedError string
	}{
		{
			name:   "success - returns links",
			userID: userID,
			limit:  10,
			offset: 0,
			setupMocks: func(repo *mocksR.MockILinksRepo) {
				repo.EXPECT().GetUserLinks(mock.Anything, userID, int32(10), int32(0)).
					Return([]string{"example.com/path"}, int32(1), nil).Once()
			},
			expectedLinks: []string{"example.com/path"},
			expectedTotal: 1,
		},
		{
			name:          "error - missing user id",
			userID:        nil,
			limit:         10,
			offset:        0,
			setupMocks:    func(*mocksR.MockILinksRepo) {},
			expectedError: "failed to get userID from ctx",
		},
		{
			name:          "error - wrong user id type",
			userID:        12345,
			limit:         10,
			offset:        0,
			setupMocks:    func(*mocksR.MockILinksRepo) {},
			expectedError: "failed to get userID from ctx",
		},
		{
			name:   "error - repo failure",
			userID: userID,
			limit:  5,
			offset: 5,
			setupMocks: func(repo *mocksR.MockILinksRepo) {
				repo.EXPECT().GetUserLinks(mock.Anything, userID, int32(5), int32(5)).
					Return(nil, int32(0), errors.New("db error")).Once()
			},
			expectedError: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocksR.NewMockILinksRepo(t)
			cache := mocksC.NewMockICache(t)
			jwtConf := mocksC.NewMockIJwtConf(t)
			tt.setupMocks(repo)

			svc := service.NewLinksService(repo, cache, &fakeProducer{}, jwtConf)

			ctx := context.Background()
			if tt.userID != nil {
				ctx = context.WithValue(ctx, jwt.ClaimsCtxKey, tt.userID)
			}

			links, total, err := svc.GetUserLinks(ctx, tt.limit, tt.offset)
			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, links)
				assert.Equal(t, int32(0), total)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedLinks, links)
			assert.Equal(t, tt.expectedTotal, total)
		})
	}
}
