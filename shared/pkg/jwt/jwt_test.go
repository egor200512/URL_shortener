package jwt

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGenerateToken_ExpiresAt(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	duration := 5 * time.Minute
	before := time.Now().Add(duration).Unix()

	token, err := GenerateToken(userID, "user@example.com", []byte("secret"), duration)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	after := time.Now().Add(duration).Unix()
	claims, err := VerifyToken(token, []byte("secret"))
	if err != nil {
		t.Fatalf("VerifyToken() error = %v", err)
	}

	if claims.ExpiresAt < before || claims.ExpiresAt > after {
		t.Fatalf("ExpiresAt = %d, want between %d and %d", claims.ExpiresAt, before, after)
	}
}

func TestVerifyToken(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	const email = "user@example.com"
	const secret = "secret"

	validToken, err := GenerateToken(userID, email, []byte(secret), time.Minute)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v", err)
	}

	expiredToken, err := GenerateToken(userID, email, []byte(secret), -time.Minute)
	if err != nil {
		t.Fatalf("GenerateToken() expired error = %v", err)
	}

	tests := []struct {
		name       string
		token      string
		secret     string
		wantID     uuid.UUID
		wantEmail  string
		wantErrSub string
	}{
		{
			name:      "valid token",
			token:     validToken,
			secret:    secret,
			wantID:    userID,
			wantEmail: email,
		},
		{
			name:       "expired token",
			token:      expiredToken,
			secret:     secret,
			wantErrSub: "invalid token",
		},
		{
			name:       "wrong secret",
			token:      validToken,
			secret:     "wrong-secret",
			wantErrSub: "invalid token",
		},
		{
			name:       "garbage token",
			token:      "bad-token",
			secret:     secret,
			wantErrSub: "invalid token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			claims, err := VerifyToken(tt.token, []byte(tt.secret))
			assertErrContains(t, err, tt.wantErrSub)
			if tt.wantErrSub != "" {
				return
			}

			if claims.ID != tt.wantID {
				t.Fatalf("ID = %s, want %s", claims.ID, tt.wantID)
			}
			if claims.Email != tt.wantEmail {
				t.Fatalf("Email = %q, want %q", claims.Email, tt.wantEmail)
			}
		})
	}
}

func assertErrContains(t *testing.T, err error, wantSub string) {
	t.Helper()

	if wantSub == "" {
		if err != nil {
			t.Fatalf("err = %v, want nil", err)
		}
		return
	}

	if err == nil {
		t.Fatalf("err = nil, want containing %q", wantSub)
	}
	if !strings.Contains(err.Error(), wantSub) {
		t.Fatalf("err = %q, want containing %q", err.Error(), wantSub)
	}
}
