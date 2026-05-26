package links

import (
	"context"
	"fmt"
)

func (repo *linksRepository) DeleteLink(ctx context.Context, shortLink string, userID string) error {
	query := fmt.Sprintf(`DELETE FROM %s where %s = $1 and %s = $2`, LINKS, SHORT_LINK, USER_ID)

	commandTag, err := repo.pool.Exec(ctx, query, shortLink, userID)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("link %s was not deleted", shortLink)
	}

	return nil
}
