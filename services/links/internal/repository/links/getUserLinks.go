package links

import (
	"context"
	"fmt"

	"github.com/georgysavva/scany/pgxscan"
)

type linkRow struct {
	OriginalLink string `db:"original_link"`
	TotalCount   int32  `db:"total_count"`
}

func (repo *linksRepo) GetUserLinks(ctx context.Context, userID string, limit, offset int32) ([]string, int32, error) {
	query := fmt.Sprintf(
		`SELECT %s AS original_link, COUNT(*) OVER() AS total_count
		FROM %s
		WHERE %s = $1
		ORDER BY %s DESC
		LIMIT $2 OFFSET $3`,
		ORIGINAL_LINK,
		LINKS,
		USER_ID,
		CREATED_AT,
	)

	var rowsData []linkRow
	if err := pgxscan.Select(ctx, repo.pool, &rowsData, query, userID, limit, offset); err != nil {
		return nil, 0, err
	}

	if len(rowsData) == 0 {
		return []string{}, 0, nil
	}

	links := make([]string, 0, len(rowsData))
	for _, row := range rowsData {
		links = append(links, row.OriginalLink)
	}

	return links, rowsData[0].TotalCount, nil
}
