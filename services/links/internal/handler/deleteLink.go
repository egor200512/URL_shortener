package handler

import (
	"context"
	"fmt"

	desc "github.com/egor200512/URL_shortener/shared/gen/links"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (handler *LinksHandler) DeleteLink(ctx context.Context, in *desc.DeleteLinkRequest) (*emptypb.Empty, error) {
	if len(in.ShortLink) == 0 {
		return nil, status.Error(codes.InvalidArgument, "short link is required")
	}

	if err := handler.linksService.DeleteLink(ctx, in.ShortLink); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to delete link: %s", err.Error()))
	}

	return &emptypb.Empty{}, nil
}
