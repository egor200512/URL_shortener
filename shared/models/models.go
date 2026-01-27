package models

import "time"

type JetStreamUnit struct {
	UserID       string    `json:"user_id"`
	ShortLink    string    `json:"short_link"`
	OriginalLink string    `json:"original_link"`
	ExecutedAt   time.Time `json:"executed_at"`
}

type LinkEvent struct {
	EventType    string    `json:"event_type"`
	UserID       string    `json:"user_id"`
	ShortLink    string    `json:"short_link"`
	OriginalLink string    `json:"original_link"`
	ExecutedAt   time.Time `json:"executed_at"`
}
