package links

import (
	"context"
	"fmt"
)

func (repo *linksRepository) DeleteLink(ctx context.Context, shortLink string, userID string) error {
	query := fmt.Sprintf(`DELETE FROM %s where %s = $1 and %s = $2`, LINKS, SHORT_LINK, USER_ID)

	if _, err := repo.pool.Exec(ctx, query, shortLink, userID); err != nil {
		return err
	}

	return nil
}
