package service

import (
	"context"
	"net/url"
)

type ILinksService interface {
	CreateLink(context.Context, *url.URL) (string, error)

	// GetUserLinks()
	// GetLink()
	// UpdateLink()
	// DeleteLink()
}
