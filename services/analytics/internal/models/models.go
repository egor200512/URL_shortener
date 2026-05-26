package models

type LinkStats struct {
	ShortLink    string
	CreatedCount int32
	FetchedCount int32
	DeletedCount int32
	TotalCount   int32
}
