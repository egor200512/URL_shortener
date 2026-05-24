package auth

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/egor200512/URL_shortener/services/auth/internal/mocks"
	"github.com/egor200512/URL_shortener/services/auth/internal/models"
	jwtpkg "github.com/egor200512/URL_shortener/shared/pkg/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Register(t *testing.T) {
	t.Parallel()

	const (
		email    = "user@example.com"
		password = "password"
	)

	tests := []struct {
		name       string
		user       *models.User
		checkErr   error
		insertErr  error
		wantErrSub string
		wantInsert bool
	}{
		{
			name:       "success",
			wantInsert: true,
		},
		{
			name:       "check registration error",
			checkErr:   errors.New("db failed"),
			wantErrSub: "failed to check registration",
		},
		{
			name:       "user already exists",
			user:       &models.User{ID: uuid.New(), Email: email},
			wantErrSub: "already exists",
		},
		{
			name:       "insert user error",
			insertErr:  errors.New("insert failed"),
			wantErrSub: "insert failed",
			wantInsert: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewIAuthRepo(t)
			jwtConf := mocks.NewIJwtConf(t)

			repo.On("CheckRegistration", mock.Anything, email).Return(tt.user, tt.checkErr).Once()
			if tt.wantInsert {
				repo.On("InsertUser", mock.Anything, mock.MatchedBy(func(req *models.RegisterRequest) bool {
					if req == nil {
						return false
					}
					if req.Email != email {
						return false
					}
					if len(req.Salt) != 16 {
						return false
					}
					return bcrypt.CompareHashAndPassword(req.SaltPassHash, []byte(password+string(req.Salt))) == nil
				})).Return(tt.insertErr).Once()
			}

			svc := NewAuthService(repo, jwtConf)
			err := svc.Register(context.Background(), email, password)
			assertErrContains(t, err, tt.wantErrSub)
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	t.Parallel()

	const (
		email    = "user@example.com"
		password = "password"
		secret   = "test-secret"
	)

	userID := uuid.New()
	user := newTestUser(t, userID, email, password)

	tests := []struct {
		name       string
		user       *models.User
		password   string
		checkErr   error
		wantErrSub string
		wantJWT    bool
	}{
		{
			name:     "success",
			user:     user,
			password: password,
			wantJWT:  true,
		},
		{
			name:       "check registration error",
			password:   password,
			checkErr:   errors.New("db failed"),
			wantErrSub: "db failed",
		},
		{
			name:       "user does not exist",
			password:   password,
			wantErrSub: "doesn't exist",
		},
		{
			name:       "wrong password",
			user:       user,
			password:   "wrong-password",
			wantErrSub: "wrong password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewIAuthRepo(t)
			jwtConf := mocks.NewIJwtConf(t)

			repo.On("CheckRegistration", mock.Anything, email).Return(tt.user, tt.checkErr).Once()
			if tt.wantJWT {
				jwtConf.On("Secret").Return(secret).Once()
				jwtConf.On("AccessExp").Return(time.Duration(1)).Once()
			}

			svc := NewAuthService(repo, jwtConf)
			token, err := svc.Login(context.Background(), email, tt.password)
			assertErrContains(t, err, tt.wantErrSub)

			if !tt.wantJWT {
				if token != nil {
					t.Fatalf("token = %#v, want nil", token)
				}
				return
			}

			if token == nil {
				t.Fatal("expected token")
			}
			if token.Token == "" {
				t.Fatal("expected non-empty access token")
			}
			if string(token.Salt) != string(tt.user.Salt) {
				t.Fatalf("token salt = %q, want %q", string(token.Salt), string(tt.user.Salt))
			}

			claims, err := jwtpkg.VerifyToken(token.Token, []byte(secret))
			if err != nil {
				t.Fatalf("verify generated token: %v", err)
			}
			if claims.ID != userID {
				t.Fatalf("claims ID = %s, want %s", claims.ID, userID)
			}
			if claims.Email != email {
				t.Fatalf("claims email = %q, want %q", claims.Email, email)
			}
		})
	}
}

func TestAuthService_VerifyToken(t *testing.T) {
	t.Parallel()

	const secret = "test-secret"

	userID := uuid.New()
	validToken, err := jwtpkg.GenerateToken(userID, "user@example.com", []byte(secret), time.Duration(1))
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	tests := []struct {
		name       string
		token      string
		secret     string
		wantID     uuid.UUID
		wantErrSub string
	}{
		{
			name:   "success",
			token:  validToken,
			secret: secret,
			wantID: userID,
		},
		{
			name:       "invalid token",
			token:      "bad-token",
			secret:     secret,
			wantErrSub: "invalid token",
		},
		{
			name:       "wrong secret",
			token:      validToken,
			secret:     "wrong-secret",
			wantErrSub: "invalid token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewIAuthRepo(t)
			jwtConf := mocks.NewIJwtConf(t)
			jwtConf.On("Secret").Return(tt.secret).Once()

			svc := NewAuthService(repo, jwtConf)
			gotID, err := svc.VerifyToken(context.Background(), tt.token)
			assertErrContains(t, err, tt.wantErrSub)

			if tt.wantErrSub != "" {
				if gotID != uuid.Nil {
					t.Fatalf("id = %s, want nil uuid", gotID)
				}
				return
			}
			if gotID != tt.wantID {
				t.Fatalf("id = %s, want %s", gotID, tt.wantID)
			}
		})
	}
}

func newTestUser(t *testing.T, id uuid.UUID, email string, password string) *models.User {
	t.Helper()

	salt := []byte("test-salt")
	hash, err := bcrypt.GenerateFromPassword([]byte(password+string(salt)), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("generate password hash: %v", err)
	}

	return &models.User{
		ID:           id,
		Email:        email,
		Salt:         salt,
		SaltPassHash: hash,
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
