package test

import (
	"context"
	"fmt"
	"testing"

	mr "github.com/egor200512/URL_shortener/services/auth/internal/mocks"
	s "github.com/egor200512/URL_shortener/services/auth/internal/service/auth"
	mc "github.com/egor200512/URL_shortener/shared/mocks"

	"github.com/egor200512/URL_shortener/services/auth/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name           string
		email          string
		password       string
		mockCheckUser  func(*mr.MockIAuthRepo)
		mockInsertUser func(*mr.MockIAuthRepo)
		expErrContains string
	}{
		{
			name:     "Success - new user",
			email:    "new@test.com",
			password: "password123",
			mockCheckUser: func(m *mr.MockIAuthRepo) {
				m.EXPECT().CheckRegistration(mock.Anything, "new@test.com").
					Return(nil, nil).Once() // No user
			},
			mockInsertUser: func(m *mr.MockIAuthRepo) {
				m.EXPECT().InsertUser(mock.Anything, mock.AnythingOfType("*models.RegisterRequest")).
					Return(nil).Once()
			},
			expErrContains: "",
		},
		{
			name:     "User already exists",
			email:    "exists@test.com",
			password: "password123",
			mockCheckUser: func(m *mr.MockIAuthRepo) {
				m.EXPECT().CheckRegistration(mock.Anything, "exists@test.com").
					Return(&models.User{Email: "exists@test.com"}, nil).Once()
			},
			mockInsertUser: func(m *mr.MockIAuthRepo) { /* not called */ },
			expErrContains: "user exists@test.com already exists",
		},
		{
			name:     "Repo check error",
			email:    "error@test.com",
			password: "password123",
			mockCheckUser: func(m *mr.MockIAuthRepo) {
				m.EXPECT().CheckRegistration(mock.Anything, "error@test.com").
					Return(nil, fmt.Errorf("db error")).Once()
			},
			mockInsertUser: func(m *mr.MockIAuthRepo) { /* not called */ },
			expErrContains: "failed to check registration",
		},
		{
			name:     "Insert error",
			email:    "insert@test.com",
			password: "password123",
			mockCheckUser: func(m *mr.MockIAuthRepo) {
				m.EXPECT().CheckRegistration(mock.Anything, "insert@test.com").
					Return(nil, nil).Once()
			},
			mockInsertUser: func(m *mr.MockIAuthRepo) {
				m.EXPECT().InsertUser(mock.Anything, mock.AnythingOfType("*models.RegisterRequest")).
					Return(fmt.Errorf("insert failed")).Once()
			},
			expErrContains: "insert failed",
		},
		{
			name:     "Bcrypt error",
			email:    "bcrypt@test.com",
			password: string(make([]byte, 1000)), // Слишком длинный для демонстрации
			mockCheckUser: func(m *mr.MockIAuthRepo) {
				m.EXPECT().CheckRegistration(mock.Anything, "bcrypt@test.com").
					Return(nil, nil).Once()
			},
			mockInsertUser: func(m *mr.MockIAuthRepo) { /* not reached */ },
			expErrContains: "bcrypt", // Любая bcrypt ошибка
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := mr.NewMockIAuthRepo(t)
			jwtConfMock := mc.NewMockIJwtConf(t) // Если нужен

			tt.mockCheckUser(repoMock)
			tt.mockInsertUser(repoMock)

			svc := s.NewAuthService(repoMock, jwtConfMock)

			err := svc.Register(context.Background(), tt.email, tt.password)

			if tt.expErrContains == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expErrContains)
			}
		})
	}
}
