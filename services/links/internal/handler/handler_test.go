package handler

import (
	"context"
	"database/sql"
	"errors"
	"net/url"
	"testing"
	"time"

	"github.com/egor200512/URL_shortener/services/links/internal/mocks"
	"github.com/egor200512/URL_shortener/services/links/models"
	desc "github.com/egor200512/URL_shortener/shared/gen/links"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestLinksHandler_CreateLink(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		req       *desc.CreateLinkRequest
		shortLink string
		svcErr    error
		wantCode  codes.Code
		wantCall  bool
	}{
		{
			name:      "success",
			req:       &desc.CreateLinkRequest{OriginalLink: "https://example.com/path"},
			shortLink: "abc123",
			wantCode:  codes.OK,
			wantCall:  true,
		},
		{
			name:     "invalid url",
			req:      &desc.CreateLinkRequest{OriginalLink: "://bad"},
			wantCode: codes.InvalidArgument,
		},
		{
			name:     "service error",
			req:      &desc.CreateLinkRequest{OriginalLink: "https://example.com/path"},
			svcErr:   errors.New("already exists"),
			wantCode: codes.InvalidArgument,
			wantCall: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := mocks.NewILinksService(t)
			if tt.wantCall {
				svc.On("CreateLink", mock.Anything, mock.MatchedBy(func(u *url.URL) bool {
					return u != nil && u.String() == tt.req.OriginalLink
				})).Return(tt.shortLink, tt.svcErr).Once()
			}

			h := NewLinksRouter(svc)
			resp, err := h.CreateLink(context.Background(), tt.req)
			assertCode(t, err, tt.wantCode)

			if tt.wantCode == codes.OK {
				if resp == nil {
					t.Fatal("expected response")
				}
				if resp.ShortLink != tt.shortLink {
					t.Fatalf("short link = %q, want %q", resp.ShortLink, tt.shortLink)
				}
			}
		})
	}
}

func TestLinksHandler_GetOriginalLink(t *testing.T) {
	t.Parallel()

	link := testLink()

	tests := []struct {
		name     string
		req      *desc.GetOriginalLinkRequest
		link     *models.Link
		svcErr   error
		wantCode codes.Code
		wantCall bool
	}{
		{name: "success", req: &desc.GetOriginalLinkRequest{ShortLink: "abc123"}, link: link, wantCode: codes.OK, wantCall: true},
		{name: "empty short link", req: &desc.GetOriginalLinkRequest{}, wantCode: codes.InvalidArgument},
		{name: "service error", req: &desc.GetOriginalLinkRequest{ShortLink: "abc123"}, svcErr: errors.New("db failed"), wantCode: codes.Internal, wantCall: true},
		{name: "not found", req: &desc.GetOriginalLinkRequest{ShortLink: "abc123"}, wantCode: codes.NotFound, wantCall: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := mocks.NewILinksService(t)
			if tt.wantCall {
				svc.On("GetLinkInfo", mock.Anything, tt.req.ShortLink).Return(tt.link, tt.svcErr).Once()
			}

			h := NewLinksRouter(svc)
			resp, err := h.GetOriginalLink(context.Background(), tt.req)
			assertCode(t, err, tt.wantCode)

			if tt.wantCode == codes.OK && resp.OriginalUrl != link.OriginalLink {
				t.Fatalf("original url = %q, want %q", resp.OriginalUrl, link.OriginalLink)
			}
		})
	}
}

func TestLinksHandler_GetLinkInfo(t *testing.T) {
	t.Parallel()

	link := testLink()

	tests := []struct {
		name     string
		req      *desc.GetLinkInfoRequest
		link     *models.Link
		svcErr   error
		wantCode codes.Code
		wantCall bool
	}{
		{name: "success", req: &desc.GetLinkInfoRequest{ShortLink: link.ShortLink}, link: link, wantCode: codes.OK, wantCall: true},
		{name: "empty short link", req: &desc.GetLinkInfoRequest{}, wantCode: codes.InvalidArgument},
		{name: "service error", req: &desc.GetLinkInfoRequest{ShortLink: link.ShortLink}, svcErr: errors.New("db failed"), wantCode: codes.Internal, wantCall: true},
		{name: "not found", req: &desc.GetLinkInfoRequest{ShortLink: link.ShortLink}, wantCode: codes.NotFound, wantCall: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := mocks.NewILinksService(t)
			if tt.wantCall {
				svc.On("GetLinkInfo", mock.Anything, tt.req.ShortLink).Return(tt.link, tt.svcErr).Once()
			}

			h := NewLinksRouter(svc)
			resp, err := h.GetLinkInfo(context.Background(), tt.req)
			assertCode(t, err, tt.wantCode)

			if tt.wantCode == codes.OK {
				if resp.Id != link.ID.String() {
					t.Fatalf("id = %q, want %q", resp.Id, link.ID.String())
				}
				if resp.OriginalLink != link.OriginalLink {
					t.Fatalf("original link = %q, want %q", resp.OriginalLink, link.OriginalLink)
				}
			}
		})
	}
}

func TestLinksHandler_GetUserLinks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		req      *desc.GetUserLinksRequest
		links    []string
		total    int32
		svcErr   error
		wantCode codes.Code
		wantCall bool
	}{
		{name: "success", req: &desc.GetUserLinksRequest{Limit: 10, Offset: 0}, links: []string{"a", "b"}, total: 2, wantCode: codes.OK, wantCall: true},
		{name: "bad limit", req: &desc.GetUserLinksRequest{Limit: 0}, wantCode: codes.InvalidArgument},
		{name: "bad offset", req: &desc.GetUserLinksRequest{Limit: 10, Offset: -1}, wantCode: codes.InvalidArgument},
		{name: "service error", req: &desc.GetUserLinksRequest{Limit: 10, Offset: 0}, svcErr: errors.New("db failed"), wantCode: codes.Internal, wantCall: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := mocks.NewILinksService(t)
			if tt.wantCall {
				svc.On("GetUserLinks", mock.Anything, tt.req.Limit, tt.req.Offset).Return(tt.links, tt.total, tt.svcErr).Once()
			}

			h := NewLinksRouter(svc)
			resp, err := h.GetUserLinks(context.Background(), tt.req)
			assertCode(t, err, tt.wantCode)

			if tt.wantCode == codes.OK {
				if resp.TotalCount != tt.total {
					t.Fatalf("total = %d, want %d", resp.TotalCount, tt.total)
				}
				if len(resp.Links) != len(tt.links) {
					t.Fatalf("links len = %d, want %d", len(resp.Links), len(tt.links))
				}
			}
		})
	}
}

func TestLinksHandler_DeleteLink(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		req      *desc.DeleteLinkRequest
		svcErr   error
		wantCode codes.Code
		wantCall bool
	}{
		{name: "success", req: &desc.DeleteLinkRequest{ShortLink: "abc123"}, wantCode: codes.OK, wantCall: true},
		{name: "empty short link", req: &desc.DeleteLinkRequest{}, wantCode: codes.InvalidArgument},
		{name: "service error", req: &desc.DeleteLinkRequest{ShortLink: "abc123"}, svcErr: errors.New("db failed"), wantCode: codes.Internal, wantCall: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := mocks.NewILinksService(t)
			if tt.wantCall {
				svc.On("DeleteLink", mock.Anything, tt.req.ShortLink).Return(tt.svcErr).Once()
			}

			h := NewLinksRouter(svc)
			resp, err := h.DeleteLink(context.Background(), tt.req)
			assertCode(t, err, tt.wantCode)

			if tt.wantCode == codes.OK && resp == nil {
				t.Fatal("expected response")
			}
		})
	}
}

func testLink() *models.Link {
	return &models.Link{
		ID:               uuid.New(),
		UserID:           uuid.New(),
		ShortLink:        "abc123",
		OriginalLinkHost: "example.com",
		OriginalLink:     "example.com/path",
		CreatedAt:        sql.NullTime{Time: time.Now(), Valid: true},
	}
}

func assertCode(t *testing.T, err error, want codes.Code) {
	t.Helper()

	if got := status.Code(err); got != want {
		t.Fatalf("status code = %s, want %s; err = %v", got, want, err)
	}
}
