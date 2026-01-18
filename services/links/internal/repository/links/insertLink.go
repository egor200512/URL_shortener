package links

import (
	"context"
	"fmt"
	"net/url"
)

func (repo *linksRepo) InsertLink(ctx context.Context, u *url.URL, shortLink string) error {
	query := fmt.Sprintf(`INSERT (%s, %s, %s, %s, %s) INTO %s VALUES ($1, $2, $3, $4, $5)`,
		USER_ID,
		SHORT_LINK,
		ORIGINAL_LINK_HOST,
		ORIGINAL_LINK,
		CREATED_AT,
		LINKS,
	)

	if _, err := repo.pool.Exec(ctx, query); err != nil {
		return err
	}

	return nil
}
