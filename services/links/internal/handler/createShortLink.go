package handler

import (
	"context"
	"log/slog"
	"net/url"

	desc "github.com/egor200512/URL_shortener/shared/gen/links"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (handler *LinksHandler) CreateLink(ctx context.Context, in *desc.CreateLinkRequest) (*desc.CreateLinkResponse, error) {
	slog.Info("create link endpoint called", "original_link", in.OriginalLink)

	u, err := url.ParseRequestURI(in.OriginalLink)
	if err != nil {
		slog.Warn("create link validation failed", "original_link", in.OriginalLink, "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	l, err := handler.linksService.CreateLink(ctx, u)
	if err != nil {
		slog.Error("create link failed", "original_link", in.OriginalLink, "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	slog.Info("create link completed", "original_link", in.OriginalLink, "short_link", l)
	return &desc.CreateLinkResponse{ShortLink: l}, nil
}
