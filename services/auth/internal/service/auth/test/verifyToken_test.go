package test

import (
	"context"
	"testing"
	"time"

	mr "github.com/egor200512/URL_shortener/services/auth/internal/mocks"
	s "github.com/egor200512/URL_shortener/services/auth/internal/service/auth"
	mc "github.com/egor200512/URL_shortener/shared/mocks"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuthService_VerifyToken(t *testing.T) {
	testUserID := uuid.New()
	testSecret := "test-secret-key"

	// Генерируем РЕАЛЬНЫЕ JWT токены
	validToken := generateValidJWT(testUserID.String(), testSecret)
	expiredToken := generateExpiredJWT(testUserID.String(), testSecret)

	tests := []struct {
		name           string
		token          string
		mockJwtConf    func(*mc.MockIJwtConf)
		expUserID      uuid.UUID
		expErrContains string
	}{
		{
			name:  "Success - valid token",
			token: validToken,
			mockJwtConf: func(m *mc.MockIJwtConf) {
				m.EXPECT().Secret().Return(testSecret).Once()
			},
			expUserID:      testUserID,
			expErrContains: "",
		},
		{
			name:  "Expired token",
			token: expiredToken,
			mockJwtConf: func(m *mc.MockIJwtConf) {
				m.EXPECT().Secret().Return(testSecret).Once()
			},
			expUserID:      uuid.Nil,
			expErrContains: "token is expired",
		},
		{
			name:  "Invalid token - wrong secret",
			token: generateValidJWT(testUserID.String(), "wrong-secret"),
			mockJwtConf: func(m *mc.MockIJwtConf) {
				m.EXPECT().Secret().Return(testSecret).Once()
			},
			expUserID:      uuid.Nil,
			expErrContains: "signature is invalid",
		},
		{
			name:  "Empty token",
			token: "",
			mockJwtConf: func(m *mc.MockIJwtConf) {
				m.EXPECT().Secret().Return(testSecret).Once()
			},
			expUserID:      uuid.Nil,
			expErrContains: "token contains an invalid number of segments",
		},
		{
			name:  "Malformed token",
			token: "not.a.jwt.at.all!!!",
			mockJwtConf: func(m *mc.MockIJwtConf) {
				m.EXPECT().Secret().Return(testSecret).Once()
			},
			expUserID:      uuid.Nil,
			expErrContains: "token contains an invalid number of segments",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := mr.NewMockIAuthRepo(t)
			jwtConfMock := mc.NewMockIJwtConf(t)

			tt.mockJwtConf(jwtConfMock)

			svc := s.NewAuthService(repoMock, jwtConfMock)

			userID, err := svc.VerifyToken(context.Background(), tt.token)

			if tt.expErrContains == "" {
				assert.NoError(t, err)
				assert.Equal(t, tt.expUserID, userID)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expErrContains)
				assert.Equal(t, uuid.Nil, userID)
			}
		})
	}
}

// Генерация валидного JWT токена
func generateValidJWT(userID string, secret string) string {
	claims := jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

// Генерация истекшего токена
func generateExpiredJWT(userID string, secret string) string {
	claims := jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(-1 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}
