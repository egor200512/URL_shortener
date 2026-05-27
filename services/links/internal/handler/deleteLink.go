package handler

import (
	"context"
	"fmt"
	"log/slog"

	desc "github.com/egor200512/URL_shortener/shared/gen/links"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (handler *LinksHandler) DeleteLink(ctx context.Context, in *desc.DeleteLinkRequest) (*emptypb.Empty, error) {
	slog.Info("delete link endpoint called", "short_link", in.ShortLink)

	if len(in.ShortLink) == 0 {
		slog.Warn("delete link validation failed", "error", "short link is required")
		return nil, status.Error(codes.InvalidArgument, "short link is required")
	}

	if err := handler.linksService.DeleteLink(ctx, in.ShortLink); err != nil {
		slog.Error("delete link failed", "short_link", in.ShortLink, "error", err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to delete link: %s", err.Error()))
	}

	slog.Info("delete link completed", "short_link", in.ShortLink)
	return &emptypb.Empty{}, nil
}
