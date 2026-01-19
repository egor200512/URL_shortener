package links

import (
	"context"
	"fmt"

	"github.com/egor200512/URL_shortener/services/links/models"
)

func (repo *linksRepo) InsertLink(ctx context.Context, req *models.CreateLinkReq) error {
	query := fmt.Sprintf(`INSERT INTO %s (%s, %s, %s, %s) VALUES ($1, $2, $3, $4)`,
		LINKS,
		USER_ID,
		SHORT_LINK,
		ORIGINAL_LINK_HOST,
		ORIGINAL_LINK,
	)
	// u_id, err := uuid.Parse(user_id)
	// if err != nil {
	// 	return err
	// }
	if _, err := repo.pool.Exec(ctx, query, req.UserID, req.ShortLink, req.OriginalLinkHost, req.OriginalLink); err != nil {
		return err
	}

	return nil
}
