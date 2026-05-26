package analytics

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/egor200512/URL_shortener/services/analytics/internal/models"
	"github.com/pashagolub/pgxmock"
)

func TestAnalyticsRepository_InsertEvent(t *testing.T) {
	t.Parallel()

	query := regexp.QuoteMeta("INSERT INTO analytics.link_events (event_type, user_id, short_link, original_link, executed_at) VALUES ($1, $2, $3, $4, $5)")
	event := testEvent()

	tests := []struct {
		name       string
		setup      func(pgxmock.PgxPoolIface)
		wantErrSub string
	}{
		{name: "success", setup: func(db pgxmock.PgxPoolIface) {
			db.ExpectExec(query).
				WithArgs(event.EventType, event.UserID, event.ShortLink, event.OriginalLink, event.ExecutedAt).
				WillReturnResult(pgxmock.NewResult("INSERT", 1))
		}},
		{name: "exec error", setup: func(db pgxmock.PgxPoolIface) {
			db.ExpectExec(query).
				WithArgs(event.EventType, event.UserID, event.ShortLink, event.OriginalLink, event.ExecutedAt).
				WillReturnError(errors.New("insert failed"))
		}, wantErrSub: "insert failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := newMockPool(t)
			defer db.Close()
			tt.setup(db)

			repo := &analyticsRepository{pool: db}
			err := repo.InsertEvent(context.Background(), event)
			assertErrContains(t, err, tt.wantErrSub)
			assertExpectations(t, db)
		})
	}
}

func TestAnalyticsRepository_GetEvents(t *testing.T) {
	t.Parallel()

	countQuery := regexp.QuoteMeta("SELECT COUNT(*) AS total FROM analytics.link_events")
	eventsQuery := regexp.QuoteMeta("SELECT id::text, event_type, user_id::text, short_link, original_link, executed_at FROM analytics.link_events ORDER BY executed_at DESC LIMIT $1 OFFSET $2")
	event := testEvent()

	tests := []struct {
		name       string
		setup      func(pgxmock.PgxPoolIface)
		wantLen    int
		wantTotal  int32
		wantErrSub string
	}{
		{name: "success", setup: func(db pgxmock.PgxPoolIface) {
			db.ExpectQuery(countQuery).WillReturnRows(pgxmock.NewRows([]string{"total"}).AddRow(int32(1)))
			db.ExpectQuery(eventsQuery).WithArgs(int32(10), int32(0)).
				WillReturnRows(eventRows().AddRow(event.ID, event.EventType, event.UserID, event.ShortLink, event.OriginalLink, event.ExecutedAt))
		}, wantLen: 1, wantTotal: 1},
		{name: "empty", setup: func(db pgxmock.PgxPoolIface) {
			db.ExpectQuery(countQuery).WillReturnRows(pgxmock.NewRows([]string{"total"}).AddRow(int32(0)))
			db.ExpectQuery(eventsQuery).WithArgs(int32(10), int32(0)).WillReturnRows(eventRows())
		}, wantLen: 0},
		{name: "count error", setup: func(db pgxmock.PgxPoolIface) {
			db.ExpectQuery(countQuery).WillReturnError(errors.New("count failed"))
		}, wantErrSub: "count failed"},
		{name: "query error", setup: func(db pgxmock.PgxPoolIface) {
			db.ExpectQuery(countQuery).WillReturnRows(pgxmock.NewRows([]string{"total"}).AddRow(int32(1)))
			db.ExpectQuery(eventsQuery).WithArgs(int32(10), int32(0)).WillReturnError(errors.New("query failed"))
		}, wantErrSub: "query failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := newMockPool(t)
			defer db.Close()
			tt.setup(db)

			repo := &analyticsRepository{pool: db}
			got, total, err := repo.GetEvents(context.Background(), 10, 0)
			assertErrContains(t, err, tt.wantErrSub)

			if tt.wantErrSub == "" {
				if len(got) != tt.wantLen || total != tt.wantTotal {
					t.Fatalf("events len=%d total=%d, want len=%d total=%d", len(got), total, tt.wantLen, tt.wantTotal)
				}
				if tt.wantLen > 0 && got[0].ID != event.ID {
					t.Fatalf("event = %#v, want %#v", got[0], event)
				}
			}
			assertExpectations(t, db)
		})
	}
}

func TestAnalyticsRepository_GetLinkStats(t *testing.T) {
	t.Parallel()

	query := regexp.QuoteMeta("SELECT event_type, COUNT(*) FROM analytics.link_events WHERE short_link = $1 GROUP BY event_type")

	tests := []struct {
		name       string
		setup      func(pgxmock.PgxPoolIface)
		want       *models.LinkStats
		wantErrSub string
	}{
		{name: "success", setup: func(db pgxmock.PgxPoolIface) {
			rows := pgxmock.NewRows([]string{"event_type", "count"}).
				AddRow("created", int32(1)).
				AddRow("fetched", int32(2)).
				AddRow("deleted", int32(3))
			db.ExpectQuery(query).WithArgs("abc123").WillReturnRows(rows)
		}, want: &models.LinkStats{ShortLink: "abc123", CreatedCount: 1, FetchedCount: 2, DeletedCount: 3, TotalCount: 6}},
		{name: "empty", setup: func(db pgxmock.PgxPoolIface) {
			db.ExpectQuery(query).WithArgs("abc123").WillReturnRows(pgxmock.NewRows([]string{"event_type", "count"}))
		}, want: &models.LinkStats{ShortLink: "abc123"}},
		{name: "query error", setup: func(db pgxmock.PgxPoolIface) {
			db.ExpectQuery(query).WithArgs("abc123").WillReturnError(errors.New("query failed"))
		}, wantErrSub: "query failed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := newMockPool(t)
			defer db.Close()
			tt.setup(db)

			repo := &analyticsRepository{pool: db}
			got, err := repo.GetLinkStats(context.Background(), "abc123")
			assertErrContains(t, err, tt.wantErrSub)

			if tt.wantErrSub == "" && *got != *tt.want {
				t.Fatalf("stats = %#v, want %#v", got, tt.want)
			}
			assertExpectations(t, db)
		})
	}
}

func eventRows() *pgxmock.Rows {
	return pgxmock.NewRows([]string{ID, EVENT_TYPE, USER_ID, SHORT_LINK, ORIGINAL_LINK, EXECUTED_AT})
}

func testEvent() *models.LinkEvent {
	return &models.LinkEvent{
		ID:           "event-id",
		EventType:    "created",
		UserID:       "user-id",
		ShortLink:    "abc123",
		OriginalLink: "example.com/path",
		ExecutedAt:   time.Date(2026, 5, 26, 10, 0, 0, 0, time.UTC),
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
