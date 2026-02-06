package analytics

import (
	"context"
	"fmt"

	"github.com/egor200512/URL_shortener/shared/models"
)

func (repo *analyticsRepository) GetEvents(ctx context.Context, limit, offset int32) ([]*models.LinkEvent, error) {
	query := fmt.Sprintf(
		`SELECT %s, %s, %s, %s, %s FROM %s ORDER BY %s DESC LIMIT $1 OFFSET $2`,
		EVENT_TYPE,
		USER_ID,
		SHORT_LINK,
		ORIGINAL_LINK,
		EXECUTED_AT,
		EVENTS,
		EXECUTED_AT,
	)

	rows, err := repo.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]*models.LinkEvent, 0, int(limit))
	for rows.Next() {
		event := new(models.LinkEvent)
		if err := rows.Scan(
			&event.EventType,
			&event.UserID,
			&event.ShortLink,
			&event.OriginalLink,
			&event.ExecutedAt,
		); err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
