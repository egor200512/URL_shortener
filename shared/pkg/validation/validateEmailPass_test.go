package pkg

import (
	"strings"
	"testing"
)

func TestValidateEmailPassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		email      string
		password   string
		wantErrSub string
	}{
		{
			name:       "empty email",
			email:      "",
			password:   "password",
			wantErrSub: "email is empty",
		},
		{
			name:       "empty password",
			email:      "user@example.com",
			password:   "",
			wantErrSub: "password is empty",
		},
		{
			name:       "email too long",
			email:      strings.Repeat("a", maxLen+1) + "@example.com",
			password:   "password",
			wantErrSub: "email too long",
		},
		{
			name:       "password too long",
			email:      "user@example.com",
			password:   strings.Repeat("a", maxLen+1),
			wantErrSub: "password too long",
		},
		{
			name:       "invalid email",
			email:      "bad-email",
			password:   "password",
			wantErrSub: "invalid email format",
		},
		{
			name:     "valid email and password",
			email:    "user@example.com",
			password: "password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateEmailPassword(tt.email, tt.password)
			assertErrContains(t, err, tt.wantErrSub)
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
