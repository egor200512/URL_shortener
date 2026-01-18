package tests

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/egor200512/URL_shortener/services/auth/internal/handler"
	desc "github.com/egor200512/URL_shortener/shared/gen/auth"
	pkg "github.com/egor200512/URL_shortener/shared/pkg/validation"

	"github.com/egor200512/URL_shortener/services/auth/internal/mocks"
	"github.com/egor200512/URL_shortener/services/auth/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestLogin(t *testing.T) {
	testToken := &models.AccessToken{Token: "test_access_token"}

	tests := []struct {
		name       string
		email      string
		password   string
		token      *models.AccessToken
		serviceErr error
		expErr     error
	}{
		{
			name:       "Success",
			email:      "test@gmail.com",
			password:   "password123",
			token:      testToken,
			serviceErr: nil,
			expErr:     nil,
		},
		{
			name:       "Empty email",
			email:      "",
			password:   "password",
			token:      nil,
			serviceErr: nil,
			expErr:     status.Error(codes.InvalidArgument, "validation failed: email is empty"),
		},
		{
			name:       "Empty password",
			email:      "test@gmail.com",
			password:   "",
			token:      nil,
			serviceErr: nil,
			expErr:     status.Error(codes.InvalidArgument, "validation failed: password is empty"),
		},
		{
			name:       "Email too long",
			email:      strings.Repeat("a", 51),
			password:   "password",
			token:      nil,
			serviceErr: nil,
			expErr:     status.Error(codes.InvalidArgument, "validation failed: email too long (max 50 chars)"),
		},
		{
			name:       "Password too long",
			email:      "test@gmail.com",
			password:   strings.Repeat("a", 51),
			token:      nil,
			serviceErr: nil,
			expErr:     status.Error(codes.InvalidArgument, "validation failed: password too long (max 50 chars)"),
		},
		{
			name:       "Invalid email",
			email:      "invalid-email",
			password:   "password",
			token:      nil,
			serviceErr: nil,
			expErr:     status.Error(codes.InvalidArgument, "validation failed: invalid email format:"),
		},
		{
			name:       "User not found",
			email:      "notfound@gmail.com",
			password:   "password123",
			token:      nil,
			serviceErr: fmt.Errorf("user not found"),
			expErr:     status.Error(codes.Internal, "failed to login: user not found"),
		},
		{
			name:       "Wrong password",
			email:      "test@gmail.com",
			password:   "wrongpass",
			token:      nil,
			serviceErr: fmt.Errorf("invalid credentials"),
			expErr:     status.Error(codes.Internal, "failed to login: invalid credentials"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			authServiceMock := mocks.NewMockIAuthService(t)

			// Mock ТОЛЬКО для валидных входных данных
			if err := pkg.ValidateEmailPassword(test.email, test.password); err == nil {
				authServiceMock.EXPECT().
					Login(mock.Anything, test.email, test.password).
					Return(test.token, test.serviceErr).Once()
			}

			h := handler.NewAuthRouter(authServiceMock)

			resp, err := h.Login(context.Background(), &desc.LoginRequest{
				Email:    test.email,
				Password: test.password,
			})

			if test.expErr == nil {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, testToken.Token, resp.AccessToken)
			} else {
				assert.Error(t, err)
				assert.Nil(t, resp)
				st, ok := status.FromError(err)
				if ok {
					expSt, _ := status.FromError(test.expErr)
					assert.Equal(t, expSt.Code(), st.Code())
					assert.Contains(t, st.Message(), expSt.Message())
				}
			}
		})
	}
}
