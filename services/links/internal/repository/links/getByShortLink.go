package links

import (
	"context"
	"fmt"

	"github.com/egor200512/URL_shortener/services/links/models"
	"github.com/georgysavva/scany/pgxscan"
)

func (repo *linksRepo) GetByShortLink(ctx context.Context, shortLink string) (*models.Link, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE %s = $1`, LINKS, SHORT_LINK)

	rows, err := repo.pool.Query(ctx, query, shortLink)
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
