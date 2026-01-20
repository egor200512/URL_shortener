package links

import (
	"context"
	"fmt"

	"github.com/egor200512/URL_shortener/services/links/models"
)

func (repo *linksRepository) InsertLink(ctx context.Context, req *models.CreateLinkReq) (*models.Link, error) {
	query := fmt.Sprintf(`INSERT INTO %s (%s, %s, %s, %s) VALUES ($1, $2, $3, $4) RETURNING %s, %s, %s, %s, %s, %s`,
		LINKS,
		USER_ID,
		SHORT_LINK,
		ORIGINAL_LINK_HOST,
		ORIGINAL_LINK,
		ID,
		USER_ID,
		SHORT_LINK,
		ORIGINAL_LINK_HOST,
		ORIGINAL_LINK,
		CREATED_AT,
	)

	row := repo.pool.QueryRow(ctx, query, req.UserID, req.ShortLink, req.OriginalLinkHost, req.OriginalLink)

	link := &models.Link{}
	if err := row.Scan(
		&link.ID,
		&link.UserID,
		&link.ShortLink,
		&link.OriginalLinkHost,
		&link.OriginalLink,
		&link.CreatedAt,
	); err != nil {
		return nil, err
	}

	return link, nil
}
