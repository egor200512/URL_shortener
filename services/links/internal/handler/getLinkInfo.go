package handler

import (
	"context"
	"fmt"
	"log/slog"

	desc "github.com/egor200512/URL_shortener/shared/gen/links"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (handler *LinksHandler) GetLinkInfo(ctx context.Context, in *desc.GetLinkInfoRequest) (*desc.LinkInfo, error) {
	slog.Info("get link info endpoint called", "short_link", in.ShortLink)

	if len(in.ShortLink) == 0 {
		slog.Warn("get link info validation failed", "error", "short link is required")
		return nil, status.Error(codes.InvalidArgument, "short link is required")
	}

	l, err := handler.linksService.GetLinkInfo(ctx, in.ShortLink)
	if err != nil {
		slog.Error("get link info failed", "short_link", in.ShortLink, "error", err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get link info: %s", err.Error()))
	}

	if l == nil {
		slog.Warn("get link info not found", "short_link", in.ShortLink)
		return nil, status.Error(codes.NotFound, fmt.Sprintf("link for %s doesn't exist", in.ShortLink))
	}

	slog.Info("get link info completed", "short_link", in.ShortLink, "link_id", l.ID.String())
	return &desc.LinkInfo{
		Id:               l.ID.String(),
		UserId:           l.UserID.String(),
		ShortLink:        in.ShortLink,
		OriginalLinkHost: l.OriginalLinkHost,
		OriginalLink:     l.OriginalLink,
		CreatedAt:        timestamppb.New(l.CreatedAt.Time),
	}, nil
}
