package analytics

import (
	"context"
	"fmt"

	"github.com/egor200512/URL_shortener/services/analytics/internal/models"
)

func (repo *analyticsRepository) GetLinkStats(ctx context.Context, shortLink string) (*models.LinkStats, error) {
	query := fmt.Sprintf(
		`SELECT %s, COUNT(*) FROM %s WHERE %s = $1 GROUP BY %s`,
		EVENT_TYPE,
		LINK_EVENTS,
		SHORT_LINK,
		EVENT_TYPE,
	)

	rows, err := repo.pool.Query(ctx, query, shortLink)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := &models.LinkStats{
		ShortLink: shortLink,
	}

	for rows.Next() {
		var eventType string
		var count int32

		if err := rows.Scan(&eventType, &count); err != nil {
			return nil, err
		}

		switch eventType {
		case "created":
			stats.CreatedCount = count
		case "fetched":
			stats.FetchedCount = count
		case "deleted":
			stats.DeletedCount = count
		}

		stats.TotalCount += count
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}
