package handler

import (
	"github.com/egor200512/URL_shortener/services/auth/internal/service"
	desc "github.com/egor200512/URL_shortener/shared/gen/auth"
)

type AuthHandler struct {
	desc.UnimplementedAuthServiceServer
	authService service.IAuthService
}

func NewAuthRouter(authService service.IAuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}
