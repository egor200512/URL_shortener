package analytics

import (
	"context"
	"fmt"

	"github.com/egor200512/URL_shortener/services/analytics/internal/models"
)

func (repo *analyticsRepository) InsertEvent(ctx context.Context, event *models.LinkEvent) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (%s, %s, %s, %s, %s) VALUES ($1, $2, $3, $4, $5)`,
		LINK_EVENTS,
		EVENT_TYPE,
		USER_ID,
		SHORT_LINK,
		ORIGINAL_LINK,
		EXECUTED_AT,
	)

	_, err := repo.pool.Exec(
		ctx,
		query,
		event.EventType,
		event.UserID,
		event.ShortLink,
		event.OriginalLink,
		event.ExecutedAt,
	)
	return err
}
