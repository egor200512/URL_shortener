package handler

import (
	desc "github.com/egor200512/URL_shortener/services/links/internal/gen_links"
	"github.com/egor200512/URL_shortener/services/links/internal/service"
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
