package analytics

import (
	"context"
	"fmt"

	"github.com/egor200512/URL_shortener/services/analytics/internal/models"
	"github.com/georgysavva/scany/pgxscan"
)

func (repo *analyticsRepository) GetEvents(ctx context.Context, limit, offset int32) ([]*models.LinkEvent, int32, error) {
	type eventTotal struct {
		Total int32 `db:"total"`
	}

	countQuery := fmt.Sprintf(`SELECT COUNT(*) AS total FROM %s`, LINK_EVENTS)

	var total eventTotal
	if err := pgxscan.Get(ctx, repo.pool, &total, countQuery); err != nil {
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

	events := make([]*models.LinkEvent, 0, limit)
	if err := pgxscan.Select(ctx, repo.pool, &events, query, limit, offset); err != nil {
		return nil, 0, err
	}

	return events, total.Total, nil
}
