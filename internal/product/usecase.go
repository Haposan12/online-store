package product

import (
	beegoContext "github.com/beego/beego/v2/server/web/context"
	"github.com/online-store/internal/domain"
	"github.com/online-store/pkg/database"
)

type UseCase interface {
	GetListProduct(beegoCtx *beegoContext.Context, req domain.GetProductListRequest) (*database.Paginator, error)
}
