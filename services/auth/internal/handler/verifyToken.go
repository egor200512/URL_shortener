package handler

import (
	"context"
	"log/slog"

	desc "github.com/egor200512/URL_shortener/shared/gen/auth"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (router *AuthHandler) VerifyToken(ctx context.Context, in *desc.VerifyTokenRequest) (*desc.VerifyTokenResponse, error) {
	slog.Info("verify token endpoint called")

	if len(in.AccessToken) == 0 {
		slog.Warn("verify token validation failed", "error", "no token")
		return nil, status.Errorf(codes.InvalidArgument, "no token")
	}

	user_id, err := router.authService.VerifyToken(ctx, in.AccessToken)
	if err != nil {
		slog.Error("verify token failed", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to verify token: %s", err.Error())
	}

	slog.Info("verify token completed", "user_id", user_id.String())
	return &desc.VerifyTokenResponse{UserId: user_id.String()}, nil
}
