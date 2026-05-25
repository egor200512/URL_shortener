package events

import "time"

const (
	LinkCreated = "created"
	LinkFetched = "fetched"
	LinkDeleted = "deleted"
)

type LinkEvent struct {
	EventID      string    `json:"event_id"`
	EventType    string    `json:"event_type"`
	UserID       string    `json:"user_id"`
	ShortLink    string    `json:"short_link"`
	OriginalLink string    `json:"original_link"`
	ExecutedAt   time.Time `json:"executed_at"`
}
