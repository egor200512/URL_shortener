package handler

import (
	"context"

	desc "github.com/egor200512/URL_shortener/services/auth/internal/gen_auth"
	pkg "github.com/egor200512/URL_shortener/shared/pkg/validation"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (router *AuthHandler) Register(ctx context.Context, in *desc.RegisterRequest) (*emptypb.Empty, error) {
	if err := pkg.ValidateEmailPassword(in.Email, in.Password); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %s", err.Error())
	}

	if err := router.authService.Register(ctx, in.Email, in.Password); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register: %s", err.Error())
	}

	return nil, nil
}
