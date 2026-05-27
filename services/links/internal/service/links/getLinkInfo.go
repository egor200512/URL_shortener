package links

import (
	"context"
	"log/slog"

	"github.com/egor200512/URL_shortener/services/links/models"
	"github.com/egor200512/URL_shortener/shared/pkg/events"
)

func (service *linksService) GetLinkInfo(ctx context.Context, shortLink string) (*models.Link, error) {
	cachedLink, err := service.cache.GetShort(ctx, shortLink)
	if err != nil {
		slog.Warn("failed to get link from cache", "short_link", shortLink, "error", err)
	}

	if cachedLink != nil {
		slog.Info("link found in cache", "short_link", shortLink)
		l := modelLinkFromCache(cachedLink)
		service.publishFetchedEvent(ctx, l)
		return l, nil
	}

	l, err := service.linksRepo.GetByShortLink(ctx, shortLink)
	if err != nil {
		return nil, err
	}

	if l == nil {
		return nil, nil
	}

	if err = service.cache.SetShort(ctx, shortLink, cacheLinkFromModel(l)); err != nil {
		slog.Warn("failed to cache link", "short_link", shortLink, "error", err)
	}

	service.publishFetchedEvent(ctx, l)

	return l, nil
}

func (service *linksService) publishFetchedEvent(ctx context.Context, link *models.Link) {
	service.publishLinkEvent(ctx, events.LinkFetched, link)
}
