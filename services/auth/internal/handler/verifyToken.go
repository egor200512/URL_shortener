package handler

import (
	"context"

	desc "github.com/egor200512/URL_shortener/shared/gen/auth"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (router *AuthHandler) VerifyToken(ctx context.Context, in *desc.VerifyTokenRequest) (*desc.VerifyTokenResponse, error) {
	if len(in.AccessToken) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "no token")
	}

	user_id, err := router.authService.VerifyToken(ctx, in.AccessToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to verify token: %s", err.Error())
	}

	return &desc.VerifyTokenResponse{UserId: user_id.String()}, nil
}
