package links

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/egor200512/URL_shortener/services/links/internal/mocks"
	"github.com/egor200512/URL_shortener/services/links/models"
	cachepkg "github.com/egor200512/URL_shortener/shared/pkg/cache"
	"github.com/egor200512/URL_shortener/shared/pkg/events"
	"github.com/egor200512/URL_shortener/shared/pkg/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

func TestLinksService_CreateLink(t *testing.T) {
	t.Parallel()

	const userID = "user-id"
	u := mustParseURL(t, "https://example.com/path")
	created := testLink("abc123")

	tests := []struct {
		name        string
		ctx         context.Context
		existing    *models.Link
		lookupErr   error
		insertErr   error
		cacheErr    error
		wantErrSub  string
		wantInsert  bool
		wantPublish bool
	}{
		{name: "success", ctx: contextWithUser(userID), wantInsert: true, wantPublish: true},
		{name: "original lookup error", ctx: contextWithUser(userID), lookupErr: errors.New("lookup failed"), wantErrSub: "lookup failed"},
		{name: "already exists", ctx: contextWithUser(userID), existing: created, wantErrSub: "already exists"},
		{name: "missing user id", ctx: context.Background(), wantErrSub: "failed to get userID"},
		{name: "insert error", ctx: contextWithUser(userID), insertErr: errors.New("insert failed"), wantErrSub: "insert failed", wantInsert: true},
		{name: "cache error is ignored", ctx: contextWithUser(userID), cacheErr: errors.New("cache failed"), wantInsert: true, wantPublish: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewILinksRepo(t)
			cache := mocks.NewICache(t)
			producer := mocks.NewIProducer(t)

			repo.On("GetByOriginalLink", mock.Anything, "example.com/path").Return(tt.existing, tt.lookupErr).Once()
			if tt.lookupErr == nil && tt.existing == nil {
				repo.On("GetByShortLink", mock.Anything, mock.AnythingOfType("string")).Return((*models.Link)(nil), nil).Once()
			}
			if tt.wantInsert {
				repo.On("InsertLink", mock.Anything, mock.MatchedBy(func(req *models.CreateLinkReq) bool {
					return req != nil &&
						req.UserID == userID &&
						req.ShortLink != "" &&
						req.OriginalLinkHost == "example.com" &&
						req.OriginalLink == "example.com/path"
				})).Return(created, tt.insertErr).Once()
			}
			if tt.wantInsert && tt.insertErr == nil {
				cache.On("SetShort", mock.Anything, mock.AnythingOfType("string"), mock.MatchedBy(func(link *cachepkg.Link) bool {
					return cacheLinkMatchesModel(link, created)
				})).Return(tt.cacheErr).Once()
			}
			if tt.wantPublish {
				expectLinkEventPublish(producer, events.LinkCreated, "links.created", created)
			}

			svc := NewLinksService(repo, cache, producer, nil)
			got, err := svc.CreateLink(tt.ctx, u)
			assertErrContains(t, err, tt.wantErrSub)

			if tt.wantErrSub == "" && got == "" {
				t.Fatal("expected non-empty short link")
			}
		})
	}
}

func TestLinksService_GetLinkInfo(t *testing.T) {
	t.Parallel()

	link := testLink("abc123")

	tests := []struct {
		name        string
		cacheLink   *cachepkg.Link
		cacheErr    error
		repoLink    *models.Link
		repoErr     error
		cacheSet    bool
		wantLink    *models.Link
		wantErrSub  string
		wantPublish bool
	}{
		{name: "cache hit", cacheLink: cacheLinkFromModel(link), wantLink: link, wantPublish: true},
		{name: "cache miss repo hit", repoLink: link, cacheSet: true, wantLink: link, wantPublish: true},
		{name: "cache error repo hit", cacheErr: errors.New("cache failed"), repoLink: link, cacheSet: true, wantLink: link, wantPublish: true},
		{name: "repo not found", wantLink: nil},
		{name: "repo error", repoErr: errors.New("db failed"), wantErrSub: "db failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewILinksRepo(t)
			cache := mocks.NewICache(t)
			producer := mocks.NewIProducer(t)

			cache.On("GetShort", mock.Anything, "abc123").Return(tt.cacheLink, tt.cacheErr).Once()
			if tt.cacheLink == nil {
				repo.On("GetByShortLink", mock.Anything, "abc123").Return(tt.repoLink, tt.repoErr).Once()
			}
			if tt.cacheSet {
				cache.On("SetShort", mock.Anything, "abc123", mock.MatchedBy(func(link *cachepkg.Link) bool {
					return cacheLinkMatchesModel(link, tt.repoLink)
				})).Return(nil).Once()
			}
			if tt.wantPublish {
				expectLinkEventPublish(producer, events.LinkFetched, "links.fetched", tt.wantLink)
			}

			svc := NewLinksService(repo, cache, producer, nil)
			got, err := svc.GetLinkInfo(context.Background(), "abc123")
			assertErrContains(t, err, tt.wantErrSub)

			if !linkMatches(got, tt.wantLink) {
				t.Fatalf("link = %#v, want %#v", got, tt.wantLink)
			}
		})
	}
}

func TestLinksService_GetUserLinks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		ctx        context.Context
		repoErr    error
		wantErrSub string
		wantCall   bool
	}{
		{name: "success", ctx: contextWithUser("user-id"), wantCall: true},
		{name: "missing user id", ctx: context.Background(), wantErrSub: "failed to get userID"},
		{name: "repo error", ctx: contextWithUser("user-id"), repoErr: errors.New("db failed"), wantErrSub: "db failed", wantCall: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewILinksRepo(t)
			cache := mocks.NewICache(t)
			if tt.wantCall {
				repo.On("GetUserLinks", mock.Anything, "user-id", int32(10), int32(2)).
					Return([]string{"a", "b"}, int32(2), tt.repoErr).Once()
			}

			svc := NewLinksService(repo, cache, nil, nil)
			links, total, err := svc.GetUserLinks(tt.ctx, 10, 2)
			assertErrContains(t, err, tt.wantErrSub)

			if tt.wantErrSub == "" && (len(links) != 2 || total != 2) {
				t.Fatalf("links=%v total=%d, want 2 links and total 2", links, total)
			}
		})
	}
}

func TestLinksService_DeleteLink(t *testing.T) {
	t.Parallel()

	link := testLink("abc123")

	tests := []struct {
		name        string
		ctx         context.Context
		link        *models.Link
		lookupErr   error
		deleteErr   error
		cacheErr    error
		wantErrSub  string
		wantDelete  bool
		wantPublish bool
	}{
		{name: "success", ctx: contextWithUser("user-id"), link: link, wantDelete: true, wantPublish: true},
		{name: "lookup error", ctx: contextWithUser("user-id"), lookupErr: errors.New("lookup failed"), wantErrSub: "lookup failed"},
		{name: "not found", ctx: contextWithUser("user-id"), wantErrSub: "doesn't exist"},
		{name: "missing user id", ctx: context.Background(), link: link, wantErrSub: "failed to get userID"},
		{name: "delete error", ctx: contextWithUser("user-id"), link: link, deleteErr: errors.New("delete failed"), wantErrSub: "delete failed", wantDelete: true},
		{name: "cache error is ignored", ctx: contextWithUser("user-id"), link: link, cacheErr: errors.New("cache failed"), wantDelete: true, wantPublish: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := mocks.NewILinksRepo(t)
			cache := mocks.NewICache(t)
			producer := mocks.NewIProducer(t)
			repo.On("GetByShortLink", mock.Anything, "abc123").Return(tt.link, tt.lookupErr).Once()
			if tt.wantDelete {
				repo.On("DeleteLink", mock.Anything, "abc123", "user-id").Return(tt.deleteErr).Once()
			}
			if tt.wantDelete && tt.deleteErr == nil {
				cache.On("DelShort", mock.Anything, "abc123").Return(tt.cacheErr).Once()
			}
			if tt.wantPublish {
				expectLinkEventPublish(producer, events.LinkDeleted, "links.deleted", tt.link)
			}

			svc := NewLinksService(repo, cache, producer, nil)
			err := svc.DeleteLink(tt.ctx, "abc123")
			assertErrContains(t, err, tt.wantErrSub)
		})
	}
}

func expectLinkEventPublish(producer *mocks.IProducer, eventType, subject string, link *models.Link) {
	producer.On("Subject", eventType).Return(subject).Once()
	producer.On("Publish", mock.Anything, subject, mock.MatchedBy(func(payload []byte) bool {
		var event events.LinkEvent
		if err := json.Unmarshal(payload, &event); err != nil {
			return false
		}
		return event.EventID != "" &&
			event.EventType == eventType &&
			event.UserID == link.UserID.String() &&
			event.ShortLink == link.ShortLink &&
			event.OriginalLink == link.OriginalLink &&
			!event.ExecutedAt.IsZero()
	})).Return(nil).Once()
}

func testLink(short string) *models.Link {
	return &models.Link{
		ID:               uuid.New(),
		UserID:           uuid.New(),
		ShortLink:        short,
		OriginalLinkHost: "example.com",
		OriginalLink:     "example.com/path",
		CreatedAt:        sql.NullTime{Time: time.Now(), Valid: true},
	}
}

func cacheLinkMatchesModel(got *cachepkg.Link, want *models.Link) bool {
	if got == nil || want == nil {
		return got == nil && want == nil
	}

	return got.ID == want.ID &&
		got.UserID == want.UserID &&
		got.ShortLink == want.ShortLink &&
		got.OriginalLinkHost == want.OriginalLinkHost &&
		got.OriginalLink == want.OriginalLink &&
		got.CreatedAt.Equal(want.CreatedAt.Time)
}

func linkMatches(got *models.Link, want *models.Link) bool {
	if got == nil || want == nil {
		return got == nil && want == nil
	}

	return got.ID == want.ID &&
		got.UserID == want.UserID &&
		got.ShortLink == want.ShortLink &&
		got.OriginalLinkHost == want.OriginalLinkHost &&
		got.OriginalLink == want.OriginalLink &&
		got.CreatedAt.Valid == want.CreatedAt.Valid &&
		got.CreatedAt.Time.Equal(want.CreatedAt.Time)
}

func contextWithUser(userID string) context.Context {
	return context.WithValue(context.Background(), jwt.ClaimsCtxKey, userID)
}

func mustParseURL(t *testing.T, raw string) *url.URL {
	t.Helper()

	u, err := url.ParseRequestURI(raw)
	if err != nil {
		t.Fatalf("parse url: %v", err)
	}
	return u
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
