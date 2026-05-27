package handler

import (
	"context"
	"log/slog"

	desc "github.com/egor200512/URL_shortener/shared/gen/auth"
	pkg "github.com/egor200512/URL_shortener/shared/pkg/validation"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (router *AuthHandler) Register(ctx context.Context, in *desc.RegisterRequest) (*emptypb.Empty, error) {
	slog.Info("register endpoint called", "email", in.Email)

	if err := pkg.ValidateEmailPassword(in.Email, in.Password); err != nil {
		slog.Warn("register validation failed", "email", in.Email, "error", err)
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %s", err.Error())
	}

	if err := router.authService.Register(ctx, in.Email, in.Password); err != nil {
		slog.Error("register failed", "email", in.Email, "error", err)
		return nil, status.Errorf(codes.Internal, "failed to register: %s", err.Error())
	}

	slog.Info("register completed", "email", in.Email)
	return &emptypb.Empty{}, nil
}
