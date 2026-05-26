package models

import "time"

type LinkEvent struct {
	ID           string    `db:"id"`
	EventType    string    `db:"event_type"`
	UserID       string    `db:"user_id"`
	ShortLink    string    `db:"short_link"`
	OriginalLink string    `db:"original_link"`
	ExecutedAt   time.Time `db:"executed_at"`
}

type LinkStats struct {
	ShortLink    string `db:"short_link"`
	CreatedCount int32  `db:"created_count"`
	FetchedCount int32  `db:"fetched_count"`
	DeletedCount int32  `db:"deleted_count"`
	TotalCount   int32  `db:"total_count"`
}
