package handler

import (
	"context"
	"fmt"

	desc "github.com/egor200512/URL_shortener/shared/gen/links"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (handler *LinksHandler) GetLinkInfo(ctx context.Context, in *desc.GetLinkInfoRequest) (*desc.LinkInfo, error) {
	if len(in.ShortLink) == 0 {
		return nil, status.Error(codes.InvalidArgument, "short link is required")
	}

	l, err := handler.linksService.GetLinkInfo(ctx, in.ShortLink)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Errorf("failed to get link info: %s", err.Error()).Error())
	}

	if l == nil {
		return nil, status.Error(codes.NotFound, fmt.Errorf("link for %s doesn't exist", in.ShortLink).Error())
	}

	return &desc.LinkInfo{
		Id:              l.ID.String(),
		UserId:          l.UserID.String(),
		ShortLink:       in.ShortLink,
		OriginalUrlHost: l.OriginalLinkHost,
		OriginalUrl:     l.OriginalLink,
		CreatedAt:       timestamppb.New(l.CreatedAt.Time),
	}, nil
}
