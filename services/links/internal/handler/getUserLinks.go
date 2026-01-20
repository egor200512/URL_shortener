package handler

import (
	"context"
	"fmt"

	desc "github.com/egor200512/URL_shortener/shared/gen/links"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (handler *LinksHandler) GetUserLinks(ctx context.Context, in *desc.GetUserLinksRequest) (*desc.GetUserLinksResponse, error) {
	if in.Limit <= 0 {
		return nil, status.Error(codes.InvalidArgument, "limit must be positive")
	}
	if in.Offset < 0 {
		return nil, status.Error(codes.InvalidArgument, "offset must be non-negative")
	}

	links, total, err := handler.linksService.GetUserLinks(ctx, in.Limit, in.Offset)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get user links: %s", err.Error()))
	}

	resp := &desc.GetUserLinksResponse{
		Links:      make([]string, 0, len(links)),
		TotalCount: total,
	}

	resp.Links = append(resp.Links, links...)

	return resp, nil
}
