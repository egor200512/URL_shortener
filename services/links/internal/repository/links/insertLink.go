package links

import (
	"context"
	"fmt"
	"net/url"

	"github.com/google/uuid"
)

func (repo *linksRepo) InsertLink(ctx context.Context, u *url.URL, shortLink string, user_id string) error {
	query := fmt.Sprintf(`INSERT INTO %s (%s, %s, %s, %s) VALUES ($1, $2, $3, $4)`,
		LINKS,
		USER_ID,
		SHORT_LINK,
		ORIGINAL_LINK_HOST,
		ORIGINAL_LINK,
	)
	u_id, err := uuid.Parse(user_id)
	if err != nil {
		return err
	}
	if _, err := repo.pool.Exec(ctx, query, u_id, shortLink, u.Host, u.Host+u.Path); err != nil {
		return err
	}

	return nil
}
