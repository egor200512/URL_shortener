package handler

import (
	"context"

	desc "github.com/egor200512/URL_shortener/shared/gen/links"
)

func (handler *LinksHandler) GetLinkInfo(ctx context.Context, in *desc.GetLinkInfoRequest) (*desc.LinkInfo, error) {
	return nil, nil
}
