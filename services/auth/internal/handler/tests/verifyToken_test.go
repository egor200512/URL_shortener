package tests

import (
	"context"
	"testing"

	auth "github.com/egor200512/URL_shortener/services/auth/internal/handler"
	mocks "github.com/egor200512/URL_shortener/services/auth/internal/mocks"
	desc "github.com/egor200512/URL_shortener/shared/gen/auth"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAuthHandler_VerifyToken(t *testing.T) {
	tests := []struct {
		name        string
		req         *desc.VerifyTokenRequest
		mockService func(*mocks.MockIAuthService)
		expUserID   string
		expErr      bool
		expErrCode  codes.Code
	}{
		{
			name: "Success - valid token",
			req:  &desc.VerifyTokenRequest{AccessToken: "valid-token"},
			mockService: func(m *mocks.MockIAuthService) {
				userID := uuid.New()
				m.EXPECT().
					VerifyToken(mock.Anything, "valid-token").
					Return(userID, nil). // Только 2 значения!
					Once()
			},
			expUserID:  "", // заполнится в проверке
			expErr:     false,
			expErrCode: 0,
		},
		{
			name: "Empty token",
			req:  &desc.VerifyTokenRequest{AccessToken: ""},
			mockService: func(m *mocks.MockIAuthService) {
				// НЕ вызывается, т.к. хендлер проверяет заранее
			},
			expUserID:  "",
			expErr:     true,
			expErrCode: codes.InvalidArgument,
		},
		{
			name: "Service error",
			req:  &desc.VerifyTokenRequest{AccessToken: "invalid-token"},
			mockService: func(m *mocks.MockIAuthService) {
				m.EXPECT().
					VerifyToken(mock.Anything, "invalid-token").
					Return(uuid.Nil, assert.AnError).
					Once()
			},
			expUserID:  "",
			expErr:     true,
			expErrCode: codes.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svcMock := mocks.NewMockIAuthService(t)
			tt.mockService(svcMock)

			h := auth.NewAuthRouter(svcMock)

			resp, err := h.VerifyToken(context.Background(), tt.req)

			if tt.expErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
				st, ok := status.FromError(err)
				if ok {
					assert.Equal(t, tt.expErrCode, st.Code())
				}
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.NotEmpty(t, resp.UserId)

			// Проверяем что UserId валидный UUID
			parsedID, parseErr := uuid.Parse(resp.UserId)
			assert.NoError(t, parseErr)
			assert.NotEqual(t, uuid.Nil, parsedID)
		})
	}
}
