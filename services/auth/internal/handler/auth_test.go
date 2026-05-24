package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/egor200512/URL_shortener/services/auth/internal/mocks"
	"github.com/egor200512/URL_shortener/services/auth/internal/models"
	desc "github.com/egor200512/URL_shortener/shared/gen/auth"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAuthHandler_Register(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		req        *desc.RegisterRequest
		serviceErr error
		wantCode   codes.Code
		wantCall   bool
	}{
		{
			name:     "success",
			req:      &desc.RegisterRequest{Email: "user@example.com", Password: "password"},
			wantCode: codes.OK,
			wantCall: true,
		},
		{
			name:     "invalid email",
			req:      &desc.RegisterRequest{Email: "bad-email", Password: "password"},
			wantCode: codes.InvalidArgument,
		},
		{
			name:       "service error",
			req:        &desc.RegisterRequest{Email: "user@example.com", Password: "password"},
			serviceErr: errors.New("insert failed"),
			wantCode:   codes.Internal,
			wantCall:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := mocks.NewIAuthService(t)
			if tt.wantCall {
				svc.On("Register", mock.Anything, tt.req.Email, tt.req.Password).Return(tt.serviceErr).Once()
			}

			h := NewAuthRouter(svc)

			resp, err := h.Register(context.Background(), tt.req)
			assertCode(t, err, tt.wantCode)

			if tt.wantCode == codes.OK && resp == nil {
				t.Fatal("expected response")
			}
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		req        *desc.LoginRequest
		token      *models.AccessToken
		serviceErr error
		wantCode   codes.Code
		wantToken  string
		wantCall   bool
	}{
		{
			name:      "success",
			req:       &desc.LoginRequest{Email: "user@example.com", Password: "password"},
			token:     &models.AccessToken{Token: "access-token"},
			wantCode:  codes.OK,
			wantToken: "access-token",
			wantCall:  true,
		},
		{
			name:     "invalid password",
			req:      &desc.LoginRequest{Email: "user@example.com", Password: ""},
			wantCode: codes.InvalidArgument,
		},
		{
			name:       "service error",
			req:        &desc.LoginRequest{Email: "user@example.com", Password: "password"},
			serviceErr: errors.New("bad credentials"),
			wantCode:   codes.Internal,
			wantCall:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := mocks.NewIAuthService(t)
			if tt.wantCall {
				svc.On("Login", mock.Anything, tt.req.Email, tt.req.Password).Return(tt.token, tt.serviceErr).Once()
			}

			h := NewAuthRouter(svc)

			resp, err := h.Login(context.Background(), tt.req)
			assertCode(t, err, tt.wantCode)

			if tt.wantCode == codes.OK {
				if resp == nil {
					t.Fatal("expected response")
				}
				if resp.AccessToken != tt.wantToken {
					t.Fatalf("AccessToken = %q, want %q", resp.AccessToken, tt.wantToken)
				}
			}
		})
	}
}

func TestAuthHandler_VerifyToken(t *testing.T) {
	t.Parallel()

	userID := uuid.New()

	tests := []struct {
		name       string
		req        *desc.VerifyTokenRequest
		userID     uuid.UUID
		serviceErr error
		wantCode   codes.Code
		wantUserID string
		wantCall   bool
	}{
		{
			name:       "success",
			req:        &desc.VerifyTokenRequest{AccessToken: "access-token"},
			userID:     userID,
			wantCode:   codes.OK,
			wantUserID: userID.String(),
			wantCall:   true,
		},
		{
			name:     "empty token",
			req:      &desc.VerifyTokenRequest{AccessToken: ""},
			wantCode: codes.InvalidArgument,
		},
		{
			name:       "service error",
			req:        &desc.VerifyTokenRequest{AccessToken: "access-token"},
			serviceErr: errors.New("invalid token"),
			wantCode:   codes.Internal,
			wantCall:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := mocks.NewIAuthService(t)
			if tt.wantCall {
				svc.On("VerifyToken", mock.Anything, tt.req.AccessToken).Return(tt.userID, tt.serviceErr).Once()
			}

			h := NewAuthRouter(svc)

			resp, err := h.VerifyToken(context.Background(), tt.req)
			assertCode(t, err, tt.wantCode)

			if tt.wantCode == codes.OK {
				if resp == nil {
					t.Fatal("expected response")
				}
				if resp.UserId != tt.wantUserID {
					t.Fatalf("UserId = %q, want %q", resp.UserId, tt.wantUserID)
				}
			}
		})
	}
}

func assertCode(t *testing.T, err error, want codes.Code) {
	t.Helper()

	if got := status.Code(err); got != want {
		t.Fatalf("status code = %s, want %s; err = %v", got, want, err)
	}
}
