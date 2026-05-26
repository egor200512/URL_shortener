package links

import (
	"context"
	"fmt"

	"github.com/egor200512/URL_shortener/services/links/models"
	"github.com/georgysavva/scany/pgxscan"
)

func (repo *linksRepository) GetByOriginalLink(ctx context.Context, originalLink string) (*models.Link, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE %s = $1`, LINKS, ORIGINAL_LINK)

	rows, err := repo.pool.Query(ctx, query, originalLink)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	link := &models.Link{}
	if err := pgxscan.ScanRow(link, rows); err != nil {
		return nil, err
	}

	return link, nil
}
