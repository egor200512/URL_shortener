package models

import "time"

type LinkEvent struct {
	ID           string
	EventType    string
	UserID       string
	ShortLink    string
	OriginalLink string
	ExecutedAt   time.Time
}

type LinkStats struct {
	ShortLink    string
	CreatedCount int32
	FetchedCount int32
	DeletedCount int32
	TotalCount   int32
}
