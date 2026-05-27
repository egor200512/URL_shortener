package handler

import (
	"context"
	"fmt"
	"log/slog"

	desc "github.com/egor200512/URL_shortener/shared/gen/links"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (handler *LinksHandler) GetUserLinks(ctx context.Context, in *desc.GetUserLinksRequest) (*desc.GetUserLinksResponse, error) {
	slog.Info("get user links endpoint called", "limit", in.Limit, "offset", in.Offset)

	if in.Limit <= 0 {
		slog.Warn("get user links validation failed", "limit", in.Limit, "error", "limit must be positive")
		return nil, status.Error(codes.InvalidArgument, "limit must be positive")
	}
	if in.Offset < 0 {
		slog.Warn("get user links validation failed", "offset", in.Offset, "error", "offset must be non-negative")
		return nil, status.Error(codes.InvalidArgument, "offset must be non-negative")
	}

	links, total, err := handler.linksService.GetUserLinks(ctx, in.Limit, in.Offset)
	if err != nil {
		slog.Error("get user links failed", "limit", in.Limit, "offset", in.Offset, "error", err)
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get user links: %s", err.Error()))
	}

	resp := &desc.GetUserLinksResponse{
		Links:      make([]string, 0, len(links)),
		TotalCount: total,
	}

	resp.Links = append(resp.Links, links...)

	slog.Info("get user links completed", "limit", in.Limit, "offset", in.Offset, "count", len(links), "total_count", total)
	return resp, nil
}
