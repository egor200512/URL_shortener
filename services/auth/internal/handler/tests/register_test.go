package tests

import (
	"context"
	"fmt"
	"strings"
	"testing"

	desc "github.com/egor200512/URL_shortener/services/auth/internal/gen_auth"
	"github.com/egor200512/URL_shortener/services/auth/internal/handler"
	pkg "github.com/egor200512/URL_shortener/shared/pkg/validation"

	"github.com/egor200512/URL_shortener/services/auth/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRegister(t *testing.T) {
	tests := []struct {
		name       string
		email      string
		password   string
		serviceErr error
		expErr     error
	}{
		{
			name:       "Success",
			email:      "test@gmail.com",
			password:   "password123",
			serviceErr: nil,
			expErr:     nil,
		},
		{
			name:       "Empty email",
			email:      "",
			password:   "password",
			serviceErr: nil,
			expErr:     status.Error(codes.InvalidArgument, "validation failed: email is empty"),
		},
		{
			name:       "Empty password",
			email:      "test@gmail.com",
			password:   "",
			serviceErr: nil,
			expErr:     status.Error(codes.InvalidArgument, "validation failed: password is empty"),
		},
		{
			name:       "Email too long",
			email:      strings.Repeat("a", 51),
			password:   "password",
			serviceErr: nil,
			expErr:     status.Error(codes.InvalidArgument, "validation failed: email too long (max 50 chars)"),
		},
		{
			name:       "Password too long",
			email:      "test@gmail.com",
			password:   strings.Repeat("a", 51),
			serviceErr: nil,
			expErr:     status.Error(codes.InvalidArgument, "validation failed: password too long (max 50 chars)"),
		},
		{
			name:       "Invalid email format",
			email:      "invalid-email",
			password:   "password",
			serviceErr: nil,
			expErr:     status.Error(codes.InvalidArgument, "validation failed: invalid email format:"),
		},
		{
			name:       "Service error (user exists)",
			email:      "exists@gmail.com",
			password:   "password123",
			serviceErr: fmt.Errorf("user already exists"),
			expErr:     status.Error(codes.Internal, "failed to register: user already exists"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			authServiceMock := mocks.NewMockIAuthService(t)

			if err := pkg.ValidateEmailPassword(test.email, test.password); err == nil {
				authServiceMock.EXPECT().
					Register(mock.Anything, test.email, test.password).
					Return(test.serviceErr).Once()
			}

			h := handler.NewAuthRouter(authServiceMock)

			resp, err := h.Register(context.Background(), &desc.RegisterRequest{
				Email:    test.email,
				Password: test.password,
			})

			if test.expErr == nil {
				require.NoError(t, err)
				require.NotNil(t, resp)
			} else {
				require.Error(t, err)
				require.Nil(t, resp)
				st, ok := status.FromError(err)
				require.True(t, ok)
				expSt, _ := status.FromError(test.expErr)
				assert.Equal(t, expSt.Code(), st.Code())
				assert.Contains(t, st.Message(), expSt.Message())
			}
		})
	}
}
