package links

import (
	"database/sql"

	"github.com/egor200512/URL_shortener/services/links/models"
	"github.com/egor200512/URL_shortener/shared/pkg/cache"
)

func cacheLinkFromModel(link *models.Link) *cache.Link {
	if link == nil {
		return nil
	}

	return &cache.Link{
		ID:               link.ID,
		UserID:           link.UserID,
		ShortLink:        link.ShortLink,
		OriginalLinkHost: link.OriginalLinkHost,
		OriginalLink:     link.OriginalLink,
		CreatedAt:        link.CreatedAt.Time,
	}
}

func modelLinkFromCache(link *cache.Link) *models.Link {
	if link == nil {
		return nil
	}

	return &models.Link{
		ID:               link.ID,
		UserID:           link.UserID,
		ShortLink:        link.ShortLink,
		OriginalLinkHost: link.OriginalLinkHost,
		OriginalLink:     link.OriginalLink,
		CreatedAt: sql.NullTime{
			Time:  link.CreatedAt,
			Valid: !link.CreatedAt.IsZero(),
		},
	}
}
