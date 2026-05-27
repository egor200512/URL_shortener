package handler

import (
	"context"
	"fmt"
	"log/slog"

	desc "github.com/egor200512/URL_shortener/shared/gen/links"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (handler *LinksHandler) GetOriginalLink(ctx context.Context, in *desc.GetOriginalLinkRequest) (*desc.GetOriginalLinkResponse, error) {
	slog.Info("get original link endpoint called", "short_link", in.ShortLink)

	if len(in.ShortLink) == 0 {
		slog.Warn("get original link validation failed", "error", "short link is required")
		return nil, status.Error(codes.InvalidArgument, "short link is required")
	}

	l, err := handler.linksService.GetLinkInfo(ctx, in.ShortLink)
	if err != nil {
		slog.Error("get original link failed", "short_link", in.ShortLink, "error", err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get link info: %s", err.Error()))
	}

	if l == nil {
		slog.Warn("get original link not found", "short_link", in.ShortLink)
		return nil, status.Error(codes.NotFound, fmt.Sprintf("link for %s doesn't exist", in.ShortLink))
	}

	slog.Info("get original link completed", "short_link", in.ShortLink)
	return &desc.GetOriginalLinkResponse{OriginalUrl: l.OriginalLink}, nil
}
