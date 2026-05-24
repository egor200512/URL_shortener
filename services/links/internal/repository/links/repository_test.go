package links

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/egor200512/URL_shortener/services/links/models"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock"
)

func TestLinksRepository_GetByShortLink(t *testing.T) {
	t.Parallel()

	query := regexp.QuoteMeta("SELECT * FROM links.short_links WHERE short_link = $1")
	link := testLink("abc123")

	tests := []struct {
		name       string
		setup      func(pgxmock.PgxPoolIface)
		wantLink   bool
		wantErrSub string
	}{
		{name: "found", setup: func(db pgxmock.PgxPoolIface) {
			db.ExpectQuery(query).WithArgs(link.ShortLink).WillReturnRows(linkRows().AddRow(link.ID, link.UserID, link.ShortLink, link.OriginalLinkHost, link.OriginalLink, link.CreatedAt))
		}, wantLink: true},
		{name: "not found", setup: func(db pgxmock.PgxPoolIface) {
			db.ExpectQuery(query).WithArgs(link.ShortLink).WillReturnRows(linkRows())
		}},
		{name: "query error", setup: func(db pgxmock.PgxPoolIface) {
			db.ExpectQuery(query).WithArgs(link.ShortLink).WillReturnError(errors.New("query failed"))
		}, wantErrSub: "query failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := newMockPool(t)
			defer db.Close()
			tt.setup(db)

			repo := &linksRepository{pool: db}
			got, err := repo.GetByShortLink(context.Background(), link.ShortLink)
			assertErrContains(t, err, tt.wantErrSub)
			assertLinkResult(t, got, link, tt.wantLink)
			assertExpectations(t, db)
		})
	}
}

func TestLinksRepository_GetByOriginalLink(t *testing.T) {
	t.Parallel()

	query := regexp.QuoteMeta("SELECT * FROM links.short_links WHERE original_link = $1")
	link := testLink("abc123")

	db := newMockPool(t)
	defer db.Close()
	db.ExpectQuery(query).WithArgs(link.OriginalLink).WillReturnRows(linkRows().AddRow(link.ID, link.UserID, link.ShortLink, link.OriginalLinkHost, link.OriginalLink, link.CreatedAt))

	repo := &linksRepository{pool: db}
	got, err := repo.GetByOriginalLink(context.Background(), link.OriginalLink)
	assertErrContains(t, err, "")
	assertLinkResult(t, got, link, true)
	assertExpectations(t, db)
}

func TestLinksRepository_GetUserLinks(t *testing.T) {
	t.Parallel()

	query := `(?s)SELECT original_link AS original_link, COUNT\(\*\) OVER\(\) AS total_count.*FROM links\.short_links.*WHERE user_id = \$1.*ORDER BY created_at DESC.*LIMIT \$2 OFFSET \$3`
	userID := uuid.New().String()

	tests := []struct {
		name       string
		setup      func(pgxmock.PgxPoolIface)
		wantLinks  []string
		wantTotal  int32
		wantErrSub string
	}{
		{name: "found", setup: func(db pgxmock.PgxPoolIface) {
			rows := pgxmock.NewRows([]string{"original_link", "total_count"}).
				AddRow("example.com/a", int32(2)).
				AddRow("example.com/b", int32(2))
			db.ExpectQuery(query).WithArgs(userID, int32(10), int32(0)).WillReturnRows(rows)
		}, wantLinks: []string{"example.com/a", "example.com/b"}, wantTotal: 2},
		{name: "empty", setup: func(db pgxmock.PgxPoolIface) {
			db.ExpectQuery(query).WithArgs(userID, int32(10), int32(0)).WillReturnRows(pgxmock.NewRows([]string{"original_link", "total_count"}))
		}, wantLinks: []string{}, wantTotal: 0},
		{name: "query error", setup: func(db pgxmock.PgxPoolIface) {
			db.ExpectQuery(query).WithArgs(userID, int32(10), int32(0)).WillReturnError(errors.New("query failed"))
		}, wantErrSub: "query failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := newMockPool(t)
			defer db.Close()
			tt.setup(db)

			repo := &linksRepository{pool: db}
			got, total, err := repo.GetUserLinks(context.Background(), userID, 10, 0)
			assertErrContains(t, err, tt.wantErrSub)

			if tt.wantErrSub == "" {
				if total != tt.wantTotal {
					t.Fatalf("total = %d, want %d", total, tt.wantTotal)
				}
				if strings.Join(got, ",") != strings.Join(tt.wantLinks, ",") {
					t.Fatalf("links = %v, want %v", got, tt.wantLinks)
				}
			}
			assertExpectations(t, db)
		})
	}
}

func TestLinksRepository_InsertLink(t *testing.T) {
	t.Parallel()

	query := regexp.QuoteMeta("INSERT INTO links.short_links (user_id, short_link, original_link_host, original_link) VALUES ($1, $2, $3, $4) RETURNING id, user_id, short_link, original_link_host, original_link, created_at")
	link := testLink("abc123")
	req := &models.CreateLinkReq{
		UserID:           link.UserID.String(),
		ShortLink:        link.ShortLink,
		OriginalLinkHost: link.OriginalLinkHost,
		OriginalLink:     link.OriginalLink,
	}

	tests := []struct {
		name       string
		setup      func(pgxmock.PgxPoolIface)
		wantErrSub string
	}{
		{name: "success", setup: func(db pgxmock.PgxPoolIface) {
			db.ExpectQuery(query).
				WithArgs(req.UserID, req.ShortLink, req.OriginalLinkHost, req.OriginalLink).
				WillReturnRows(linkRows().AddRow(link.ID, link.UserID, link.ShortLink, link.OriginalLinkHost, link.OriginalLink, link.CreatedAt))
		}},
		{name: "query error", setup: func(db pgxmock.PgxPoolIface) {
			db.ExpectQuery(query).
				WithArgs(req.UserID, req.ShortLink, req.OriginalLinkHost, req.OriginalLink).
				WillReturnError(errors.New("insert failed"))
		}, wantErrSub: "insert failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := newMockPool(t)
			defer db.Close()
			tt.setup(db)

			repo := &linksRepository{pool: db}
			got, err := repo.InsertLink(context.Background(), req)
			assertErrContains(t, err, tt.wantErrSub)
			assertLinkResult(t, got, link, tt.wantErrSub == "")
			assertExpectations(t, db)
		})
	}
}

func TestLinksRepository_DeleteLink(t *testing.T) {
	t.Parallel()

	query := regexp.QuoteMeta("DELETE FROM links.short_links where short_link = $1 and user_id = $2")
	userID := uuid.New().String()

	tests := []struct {
		name       string
		setup      func(pgxmock.PgxPoolIface)
		wantErrSub string
	}{
		{name: "success", setup: func(db pgxmock.PgxPoolIface) {
			db.ExpectExec(query).WithArgs("abc123", userID).WillReturnResult(pgxmock.NewResult("DELETE", 1))
		}},
		{name: "exec error", setup: func(db pgxmock.PgxPoolIface) {
			db.ExpectExec(query).WithArgs("abc123", userID).WillReturnError(errors.New("delete failed"))
		}, wantErrSub: "delete failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := newMockPool(t)
			defer db.Close()
			tt.setup(db)

			repo := &linksRepository{pool: db}
			err := repo.DeleteLink(context.Background(), "abc123", userID)
			assertErrContains(t, err, tt.wantErrSub)
			assertExpectations(t, db)
		})
	}
}

func linkRows() *pgxmock.Rows {
	return pgxmock.NewRows([]string{ID, USER_ID, SHORT_LINK, ORIGINAL_LINK_HOST, ORIGINAL_LINK, CREATED_AT})
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

func assertLinkResult(t *testing.T, got *models.Link, want *models.Link, wantLink bool) {
	t.Helper()

	if !wantLink {
		if got != nil {
			t.Fatalf("link = %#v, want nil", got)
		}
		return
	}
	if got == nil {
		t.Fatal("link = nil, want link")
	}
	if got.ID != want.ID || got.UserID != want.UserID || got.ShortLink != want.ShortLink || got.OriginalLink != want.OriginalLink {
		t.Fatalf("link = %#v, want %#v", got, want)
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
	if !strings.Contains(err.Error(), wantSub) {
		t.Fatalf("err = %q, want containing %q", err.Error(), wantSub)
	}
}
