package handler

import (
	desc "github.com/egor200512/URL_shortener/services/auth/internal/gen_auth"
	"github.com/egor200512/URL_shortener/services/auth/internal/service"
)

type AuthHandler struct {
	desc.UnimplementedAuthServiceServer
	authService service.AuthService
}

func NewAuthRouter(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}
