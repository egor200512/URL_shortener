package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	mr "github.com/egor200512/URL_shortener/services/auth/internal/mocks"
	s "github.com/egor200512/URL_shortener/services/auth/internal/service/auth"
	mc "github.com/egor200512/URL_shortener/shared/mocks"

	"github.com/egor200512/URL_shortener/services/auth/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Login(t *testing.T) {
	// Подготовка тестовых данных
	testUserID := uuid.New()
	testPassword := "password123"
	testSalt := []byte("testsalt")

	// Генерируем правильный хеш для успешного теста
	validHash, _ := bcrypt.GenerateFromPassword([]byte(testPassword+string(testSalt)), bcrypt.DefaultCost)

	tests := []struct {
		name           string
		email          string
		password       string
		mockCheckUser  func(*mr.MockIAuthRepo)
		mockJwtConf    func(*mc.MockIJwtConf)
		expErrContains string
		expToken       bool
	}{
		{
			name:     "Success - valid credentials",
			email:    "valid@test.com",
			password: testPassword,
			mockCheckUser: func(m *mr.MockIAuthRepo) {
				m.EXPECT().CheckRegistration(mock.Anything, "valid@test.com").
					Return(&models.User{
						ID:           testUserID,
						Email:        "valid@test.com",
						Salt:         testSalt,
						SaltPassHash: validHash,
					}, nil).Once()
			},
			mockJwtConf: func(m *mc.MockIJwtConf) {
				m.EXPECT().Secret().Return("test-secret").Once()
				m.EXPECT().AccessExp().Return(15 * time.Minute).Once()
			},
			expErrContains: "",
			expToken:       true,
		},
		{
			name:     "User doesn't exist",
			email:    "notfound@test.com",
			password: testPassword,
			mockCheckUser: func(m *mr.MockIAuthRepo) {
				m.EXPECT().CheckRegistration(mock.Anything, "notfound@test.com").
					Return(nil, nil).Once() // No user found
			},
			mockJwtConf:    func(m *mc.MockIJwtConf) { /* not called */ },
			expErrContains: "user notfound@test.com doesn't exist",
			expToken:       false,
		},
		{
			name:     "Repo check error",
			email:    "error@test.com",
			password: testPassword,
			mockCheckUser: func(m *mr.MockIAuthRepo) {
				m.EXPECT().CheckRegistration(mock.Anything, "error@test.com").
					Return(nil, fmt.Errorf("db connection error")).Once()
			},
			mockJwtConf:    func(m *mc.MockIJwtConf) { /* not called */ },
			expErrContains: "db connection error",
			expToken:       false,
		},
		{
			name:     "Wrong password",
			email:    "valid@test.com",
			password: "wrongpassword",
			mockCheckUser: func(m *mr.MockIAuthRepo) {
				m.EXPECT().CheckRegistration(mock.Anything, "valid@test.com").
					Return(&models.User{
						ID:           testUserID,
						Email:        "valid@test.com",
						Salt:         testSalt,
						SaltPassHash: validHash,
					}, nil).Once()
			},
			mockJwtConf:    func(m *mc.MockIJwtConf) { /* not called */ },
			expErrContains: "wrong password",
			expToken:       false,
		},
		{
			name:     "Invalid hash in database",
			email:    "badhash@test.com",
			password: testPassword,
			mockCheckUser: func(m *mr.MockIAuthRepo) {
				m.EXPECT().CheckRegistration(mock.Anything, "badhash@test.com").
					Return(&models.User{
						ID:           testUserID,
						Email:        "badhash@test.com",
						Salt:         testSalt,
						SaltPassHash: []byte("invalid-hash-format"),
					}, nil).Once()
			},
			mockJwtConf:    func(m *mc.MockIJwtConf) { /* not called */ },
			expErrContains: "wrong password",
			expToken:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := mr.NewMockIAuthRepo(t)
			jwtConfMock := mc.NewMockIJwtConf(t)

			tt.mockCheckUser(repoMock)
			tt.mockJwtConf(jwtConfMock)

			svc := s.NewAuthService(repoMock, jwtConfMock)

			token, err := svc.Login(context.Background(), tt.email, tt.password)

			if tt.expErrContains == "" {
				assert.NoError(t, err)
				assert.NotNil(t, token)
				assert.NotEmpty(t, token.Token)
				assert.Equal(t, testSalt, token.Salt)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expErrContains)
				assert.Nil(t, token)
			}
		})
	}
}
