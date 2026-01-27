package links

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/egor200512/URL_shortener/services/links/models"
	m "github.com/egor200512/URL_shortener/shared/models"
	events "github.com/egor200512/URL_shortener/shared/pkg/broker/nats/events"
	"github.com/egor200512/URL_shortener/shared/pkg/jwt"
	pkg "github.com/egor200512/URL_shortener/shared/pkg/links"
)

func (service *linksService) CreateLink(ctx context.Context, u *url.URL) (string, error) {
	l, err := service.linksRepo.GetByOriginalLink(ctx, u.Host+u.Path)
	if err != nil {
		return "", err
	}

	if l != nil {
		return "", fmt.Errorf("%s alredy exists", u.Host+u.Path)
	}

	var shortLink string

	for {
		shortLink, err = pkg.GenerateShortLink()
		if err != nil {
			return "", err
		}

		l, err := service.linksRepo.GetByShortLink(ctx, shortLink)
		if err != nil {
			return "", err
		}

		if l == nil {
			break
		}
	}

	userID, ok := ctx.Value(jwt.ClaimsCtxKey).(string)
	if !ok {
		return "", errors.New("failed to get userID from ctx")
	}

	req := &models.CreateLinkReq{
		UserID:           userID,
		ShortLink:        shortLink,
		OriginalLinkHost: u.Host,
		OriginalLink:     u.Host + u.Path,
	}

	created, err := service.linksRepo.InsertLink(ctx, req)
	if err != nil {
		return "", err
	}

	if err = service.cache.SetShort(ctx, shortLink, created); err != nil {
		log.Printf("failed to cache link: %s\n", err.Error())
	}

	msg := &m.JetStreamUnit{
		UserID:       created.UserID.String(),
		ShortLink:    created.ShortLink,
		OriginalLink: created.OriginalLink,
		ExecutedAt:   time.Now(),
	}

	encodedMsg, err := json.Marshal(msg)
	if err != nil {
		log.Printf("failed to marshal event: %s\n", err.Error())
	} else {
		if err = service.broker.Publish(
			ctx,
			events.Created(service.broker.Prefix()),
			encodedMsg,
		); err != nil {
			log.Printf("failed to publish event: %s\n", err.Error())
		}
	}

	return shortLink, nil
}
