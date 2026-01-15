package handler

import (
	"context"

	desc "github.com/egor200512/URL_shortener/services/auth/internal/gen_auth"
	pkg "github.com/egor200512/URL_shortener/shared/pkg/validation"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (router *AuthHandler) Login(ctx context.Context, in *desc.LoginRequest) (*desc.LoginResponse, error) {
	if err := pkg.ValidateEmailPassword(in.Email, in.Password); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %s", err.Error())
	}

	token, err := router.authService.Login(ctx, in.Email, in.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to login: %s", err.Error())
	}

	return &desc.LoginResponse{AccessToken: token}, nil
}
