package handler

import (
	"context"
	"net/url"

	desc "github.com/egor200512/URL_shortener/shared/gen/links"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (handler *LinksHandler) CreateLink(ctx context.Context, in *desc.CreateLinkRequest) (*desc.CreateLinkResponse, error) {
	u, err := url.ParseRequestURI(in.OriginalLink)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	l, err := handler.linksService.CreateLink(ctx, u)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &desc.CreateLinkResponse{ShortLink: l}, nil
}
