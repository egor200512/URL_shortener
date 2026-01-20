package links

import (
	r "github.com/egor200512/URL_shortener/services/links/internal/repository"
	s "github.com/egor200512/URL_shortener/services/links/internal/service"
	"github.com/egor200512/URL_shortener/shared/configs"
	c "github.com/egor200512/URL_shortener/shared/pkg/cache"
)

type linksService struct {
	linksRepo r.ILinksRepo
	cache     c.ICache
	jwtConf   configs.IJwtConf
}

func NewLinksService(authRepo r.ILinksRepo, cache c.ICache, jwtConf configs.IJwtConf) s.ILinksService {
	return &linksService{
		linksRepo: authRepo,
		cache:     cache,
		jwtConf:   jwtConf,
	}
}
