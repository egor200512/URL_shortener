package handler

import (
	"context"
	"fmt"

	desc "github.com/egor200512/URL_shortener/shared/gen/links"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (handler *LinksHandler) GetOriginalLink(ctx context.Context, in *desc.GetOriginalLinkRequest) (*desc.GetOriginalLinkResponse, error) {
	if len(in.ShortLink) == 0 {
		return nil, status.Error(codes.InvalidArgument, "short link is required")
	}

	l, err := handler.linksService.GetLinkInfo(ctx, in.ShortLink)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get link info: %s", err.Error()))
	}

	if l == nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("link for %s doesn't exist", in.ShortLink))
	}

	return &desc.GetOriginalLinkResponse{OriginalUrl: l.OriginalLink}, nil
}
