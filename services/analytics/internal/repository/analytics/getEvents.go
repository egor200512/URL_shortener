package analytics

import (
	"context"
	"fmt"

	"github.com/egor200512/URL_shortener/services/analytics/internal/models"
)

func (repo *analyticsRepository) GetEvents(ctx context.Context, limit, offset int32) ([]*models.LinkEvent, int32, error) {
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s`, LINK_EVENTS)

	var total int32
	if err := repo.pool.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(
		`SELECT %s::text, %s, %s::text, %s, %s, %s FROM %s ORDER BY %s DESC LIMIT $1 OFFSET $2`,
		ID,
		EVENT_TYPE,
		USER_ID,
		SHORT_LINK,
		ORIGINAL_LINK,
		EXECUTED_AT,
		LINK_EVENTS,
		EXECUTED_AT,
	)

	rows, err := repo.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	events := make([]*models.LinkEvent, 0, limit)
	for rows.Next() {
		event := &models.LinkEvent{}
		if err := rows.Scan(
			&event.ID,
			&event.EventType,
			&event.UserID,
			&event.ShortLink,
			&event.OriginalLink,
			&event.ExecutedAt,
		); err != nil {
			return nil, 0, err
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return events, total, nil
}
