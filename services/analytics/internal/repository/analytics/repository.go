package analytics

import (
	"context"
	"fmt"

	"github.com/egor200512/URL_shortener/shared/models"
)

const (
	EVENTS        = "analytics.link_events"
	ID            = "id"
	EVENT_TYPE    = "event_type"
	USER_ID       = "user_id"
	SHORT_LINK    = "short_link"
	ORIGINAL_LINK = "original_link"
	EXECUTED_AT   = "executed_at"
)

func (repo *analyticsRepository) InsertEvent(ctx context.Context, payload *models.LinkEvent) error {
	query := fmt.Sprintf(`INSERT INTO %s (%s, %s, %s, %s) VALUES ($1, $2, $3, $4)`,
		EVENTS,
		EVENT_TYPE,
		USER_ID,
		SHORT_LINK,
		ORIGINAL_LINK,
	)

	_, err := repo.pool.Exec(
		ctx,
		query,
		payload.EventType,
		payload.UserID,
		payload.ShortLink,
		payload.OriginalLink,
	)
	if err != nil {
		return err
	}

	return nil
}
