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
	"github.com/egor200512/URL_shortener/shared/pkg/events"
	"github.com/egor200512/URL_shortener/shared/pkg/jwt"
	pkg "github.com/egor200512/URL_shortener/shared/pkg/links"
	"github.com/google/uuid"
)

func (service *linksService) CreateLink(ctx context.Context, u *url.URL) (string, error) {
	l, err := service.linksRepo.GetByOriginalLink(ctx, u.Host+u.Path)
	if err != nil {
		return "", err
	}

	if l != nil {
		return "", fmt.Errorf("%s already exists", u.Host+u.Path)
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

	service.publishCreatedEvent(ctx, created)

	return shortLink, nil
}

func (service *linksService) publishCreatedEvent(ctx context.Context, link *models.Link) {
	if service.producer == nil || link == nil {
		return
	}

	event := &events.LinkEvent{
		EventID:      uuid.NewString(),
		EventType:    events.LinkCreated,
		UserID:       link.UserID.String(),
		ShortLink:    link.ShortLink,
		OriginalLink: link.OriginalLink,
		ExecutedAt:   time.Now(),
	}

	payload, err := json.Marshal(event)
	if err != nil {
		log.Printf("failed to marshal link event: %s\n", err.Error())
		return
	}

	if err := service.producer.Publish(ctx, service.producer.Subject(events.LinkCreated), payload); err != nil {
		log.Printf("failed to publish link event: %s\n", err.Error())
	}
}
