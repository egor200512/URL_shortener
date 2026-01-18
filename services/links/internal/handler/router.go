package handler

import (
	"github.com/egor200512/URL_shortener/services/links/internal/service"
	desc "github.com/egor200512/URL_shortener/shared/gen/links"
)

type LinksHandler struct {
	desc.UnimplementedLinksServiceServer
	linksService service.ILinksService
}

func NewLinksRouter(linksService service.ILinksService) *LinksHandler {
	return &LinksHandler{
		linksService: linksService,
	}
}
