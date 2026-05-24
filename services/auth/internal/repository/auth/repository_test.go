package auth

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/egor200512/URL_shortener/services/auth/internal/models"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock"
)

func TestAuthRepo_CheckRegistration(t *testing.T) {
	t.Parallel()

	const email = "user@example.com"
	query := regexp.QuoteMeta("SELECT * FROM auth.users WHERE email = $1")

	tests := []struct {
		name       string
		setup      func(pgxmock.PgxPoolIface, uuid.UUID, []byte, []byte)
		wantUser   bool
		wantErrSub string
	}{
		{
			name: "found",
			setup: func(db pgxmock.PgxPoolIface, id uuid.UUID, salt []byte, hash []byte) {
				rows := pgxmock.NewRows([]string{ID, EMAIL, SALT, SALT_PASS_HASH}).
					AddRow(id, email, salt, hash)

				db.ExpectQuery(query).WithArgs(email).WillReturnRows(rows)
			},
			wantUser: true,
		},
		{
			name: "not found",
			setup: func(db pgxmock.PgxPoolIface, id uuid.UUID, salt []byte, hash []byte) {
				rows := pgxmock.NewRows([]string{ID, EMAIL, SALT, SALT_PASS_HASH})

				db.ExpectQuery(query).WithArgs(email).WillReturnRows(rows)
			},
		},
		{
			name: "query error",
			setup: func(db pgxmock.PgxPoolIface, id uuid.UUID, salt []byte, hash []byte) {
				db.ExpectQuery(query).WithArgs(email).WillReturnError(errors.New("query failed"))
			},
			wantErrSub: "query failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := newMockPool(t)
			defer db.Close()

			id := uuid.New()
			salt := []byte("salt")
			hash := []byte("hash")

			tt.setup(db, id, salt, hash)

			repo := &authRepo{pool: db}
			user, err := repo.CheckRegistration(context.Background(), email)
			assertErrContains(t, err, tt.wantErrSub)

			if tt.wantUser {
				if user == nil {
					t.Fatal("user = nil, want user")
				}
				if user.ID != id {
					t.Fatalf("user ID = %s, want %s", user.ID, id)
				}
				if user.Email != email {
					t.Fatalf("user email = %q, want %q", user.Email, email)
				}
				if string(user.Salt) != string(salt) {
					t.Fatalf("user salt = %q, want %q", string(user.Salt), string(salt))
				}
				if string(user.SaltPassHash) != string(hash) {
					t.Fatalf("user hash = %q, want %q", string(user.SaltPassHash), string(hash))
				}
				assertExpectations(t, db)
				return
			}

			if user != nil {
				t.Fatalf("user = %#v, want nil", user)
			}
			assertExpectations(t, db)
		})
	}
}

func TestAuthRepo_InsertUser(t *testing.T) {
	t.Parallel()

	query := regexp.QuoteMeta("INSERT INTO auth.users (email, salt, salt_password_hash) VALUES ($1, $2, $3)")
	req := &models.RegisterRequest{
		Email:        "user@example.com",
		Salt:         []byte("salt"),
		SaltPassHash: []byte("hash"),
	}

	tests := []struct {
		name       string
		setup      func(pgxmock.PgxPoolIface)
		wantErrSub string
	}{
		{
			name: "success",
			setup: func(db pgxmock.PgxPoolIface) {
				db.ExpectExec(query).
					WithArgs(req.Email, req.Salt, req.SaltPassHash).
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
			},
		},
		{
			name: "exec error",
			setup: func(db pgxmock.PgxPoolIface) {
				db.ExpectExec(query).
					WithArgs(req.Email, req.Salt, req.SaltPassHash).
					WillReturnError(errors.New("insert failed"))
			},
			wantErrSub: "insert failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := newMockPool(t)
			defer db.Close()

			tt.setup(db)

			repo := &authRepo{pool: db}
			err := repo.InsertUser(context.Background(), req)
			assertErrContains(t, err, tt.wantErrSub)
			assertExpectations(t, db)
		})
	}
}

func newMockPool(t *testing.T) pgxmock.PgxPoolIface {
	t.Helper()

	db, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("new pgx mock pool: %v", err)
	}
	return db
}

func assertExpectations(t *testing.T, db pgxmock.PgxPoolIface) {
	t.Helper()

	if err := db.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet db expectations: %v", err)
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
	if !regexp.MustCompile(regexp.QuoteMeta(wantSub)).MatchString(err.Error()) {
		t.Fatalf("err = %q, want containing %q", err.Error(), wantSub)
	}
}
