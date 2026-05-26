package analytics

import (
	"context"
	"fmt"

	"github.com/egor200512/URL_shortener/services/analytics/internal/models"
	"github.com/georgysavva/scany/pgxscan"
)

func (repo *analyticsRepository) GetLinkStats(ctx context.Context, shortLink string) (*models.LinkStats, error) {
	type eventCount struct {
		EventType string `db:"event_type"`
		Count     int32  `db:"count"`
	}

	query := fmt.Sprintf(
		`SELECT %s, COUNT(*) FROM %s WHERE %s = $1 GROUP BY %s`,
		EVENT_TYPE,
		LINK_EVENTS,
		SHORT_LINK,
		EVENT_TYPE,
	)

	stats := &models.LinkStats{
		ShortLink: shortLink,
	}

	rows := make([]eventCount, 0)
	if err := pgxscan.Select(ctx, repo.pool, &rows, query, shortLink); err != nil {
		return nil, err
	}

	for _, row := range rows {
		switch row.EventType {
		case "created":
			stats.CreatedCount = row.Count
		case "fetched":
			stats.FetchedCount = row.Count
		case "deleted":
			stats.DeletedCount = row.Count
		}

		stats.TotalCount += row.Count
	}

	return stats, nil
}
