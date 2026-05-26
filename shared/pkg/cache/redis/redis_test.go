package redis

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/egor200512/URL_shortener/shared/pkg/cache"
	"github.com/go-redis/redismock/v9"
	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
)

func TestRedisCli_GetShort(t *testing.T) {
	t.Parallel()

	link := testCacheLink()
	rawLink, err := json.Marshal(link)
	if err != nil {
		t.Fatalf("marshal link: %v", err)
	}

	tests := []struct {
		name       string
		setup      func(redismock.ClientMock)
		wantLink   *cache.Link
		wantErrSub string
	}{
		{
			name: "cache miss",
			setup: func(mock redismock.ClientMock) {
				mock.ExpectGet("abc123").RedisNil()
			},
		},
		{
			name: "bad json",
			setup: func(mock redismock.ClientMock) {
				mock.ExpectGet("abc123").SetVal("{bad")
			},
			wantErrSub: "invalid character",
		},
		{
			name: "valid json",
			setup: func(mock redismock.ClientMock) {
				mock.ExpectGet("abc123").SetVal(string(rawLink))
			},
			wantLink: link,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, mock := redismock.NewClientMock()
			tt.setup(mock)

			redisCli := &redisCli{Cli: client}
			got, err := redisCli.GetShort(context.Background(), "abc123")
			assertErrContains(t, err, tt.wantErrSub)
			assertCacheLink(t, got, tt.wantLink)
			assertRedisExpectations(t, mock)
		})
	}
}

func TestRedisCli_SetShort(t *testing.T) {
	t.Parallel()

	link := testCacheLink()
	rawLink, err := json.Marshal(link)
	if err != nil {
		t.Fatalf("marshal link: %v", err)
	}

	tests := []struct {
		name       string
		link       *cache.Link
		setup      func(redismock.ClientMock)
		wantErrSub string
	}{
		{
			name:       "nil link",
			wantErrSub: goredis.Nil.Error(),
		},
		{
			name: "valid link",
			link: link,
			setup: func(mock redismock.ClientMock) {
				mock.ExpectSet("abc123", rawLink, 10*time.Minute).SetVal("OK")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, mock := redismock.NewClientMock()
			if tt.setup != nil {
				tt.setup(mock)
			}

			redisCli := &redisCli{Cli: client, ttl: 10 * time.Minute}
			err := redisCli.SetShort(context.Background(), "abc123", tt.link)
			assertErrContains(t, err, tt.wantErrSub)
			assertRedisExpectations(t, mock)
		})
	}
}

func TestRedisCli_DelShort(t *testing.T) {
	t.Parallel()

	client, mock := redismock.NewClientMock()
	mock.ExpectDel("abc123").SetVal(1)

	redisCli := &redisCli{Cli: client}
	err := redisCli.DelShort(context.Background(), "abc123")
	assertErrContains(t, err, "")
	assertRedisExpectations(t, mock)
}

func testCacheLink() *cache.Link {
	return &cache.Link{
		ID:               uuid.New(),
		UserID:           uuid.New(),
		ShortLink:        "abc123",
		OriginalLinkHost: "example.com",
		OriginalLink:     "example.com/path",
		CreatedAt:        time.Date(2026, 5, 26, 10, 0, 0, 0, time.UTC),
	}
}

func assertCacheLink(t *testing.T, got *cache.Link, want *cache.Link) {
	t.Helper()

	if got == nil || want == nil {
		if got != want {
			t.Fatalf("link = %#v, want %#v", got, want)
		}
		return
	}

	if got.ID != want.ID ||
		got.UserID != want.UserID ||
		got.ShortLink != want.ShortLink ||
		got.OriginalLinkHost != want.OriginalLinkHost ||
		got.OriginalLink != want.OriginalLink ||
		!got.CreatedAt.Equal(want.CreatedAt) {
		t.Fatalf("link = %#v, want %#v", got, want)
	}
}

func assertRedisExpectations(t *testing.T, mock redismock.ClientMock) {
	t.Helper()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet redis expectations: %v", err)
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
