package auth

import (
	r "github.com/egor200512/URL_shortener/services/auth/internal/repository"
	s "github.com/egor200512/URL_shortener/services/auth/internal/service"
	"github.com/egor200512/URL_shortener/shared/configs"
)

type authService struct {
	authRepo r.IAuthRepo
	jwtConf  configs.IJwtConf
}

func NewAuthService(authRepo r.IAuthRepo, jwtConf configs.IJwtConf) s.IAuthService {
	return &authService{
		authRepo: authRepo,
		jwtConf:  jwtConf,
	}
}
