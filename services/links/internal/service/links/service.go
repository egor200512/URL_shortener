package links

import (
	r "github.com/egor200512/URL_shortener/services/links/internal/repository"
	s "github.com/egor200512/URL_shortener/services/links/internal/service"
	"github.com/egor200512/URL_shortener/shared/configs"
)

type linksService struct {
	linksRepo r.ILinksRepo
	jwtConf   configs.IJwtConf
}

func NewLinksService(authRepo r.ILinksRepo, jwtConf configs.IJwtConf) s.ILinksService {
	return &linksService{
		linksRepo: authRepo,
		jwtConf:   jwtConf,
	}
}
