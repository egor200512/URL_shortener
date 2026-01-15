package auth

import (
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
