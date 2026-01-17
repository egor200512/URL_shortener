package handler

import (
	"context"
	"net/url"

	desc "github.com/egor200512/URL_shortener/services/links/internal/gen_links"
	"github.com/k0kubun/pp/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (handler *LinksHandler) CreateShortLink(ctx context.Context, in *desc.CreateShortLinkRequest) (*desc.CreateShortLinkResponse, error) {
	u, err := url.ParseRequestURI(in.OriginalLink)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	pp.Println(u)

	return &desc.CreateShortLinkResponse{}, nil
}
