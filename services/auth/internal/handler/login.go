package handler

import (
	"context"
	"log/slog"

	desc "github.com/egor200512/URL_shortener/shared/gen/auth"
	pkg "github.com/egor200512/URL_shortener/shared/pkg/validation"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (router *AuthHandler) Login(ctx context.Context, in *desc.LoginRequest) (*desc.LoginResponse, error) {
	slog.Info("login endpoint called", "email", in.Email)

	if err := pkg.ValidateEmailPassword(in.Email, in.Password); err != nil {
		slog.Warn("login validation failed", "email", in.Email, "error", err)
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %s", err.Error())
	}

	t, err := router.authService.Login(ctx, in.Email, in.Password)
	if err != nil {
		slog.Error("login failed", "email", in.Email, "error", err)
		return nil, status.Errorf(codes.Internal, "failed to login: %s", err.Error())
	}

	slog.Info("login completed", "email", in.Email)
	return &desc.LoginResponse{AccessToken: t.Token}, nil
}
