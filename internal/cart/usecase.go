package cart

import (
	beegoContext "github.com/beego/beego/v2/server/web/context"
	"github.com/online-store/internal/domain"
	"github.com/online-store/pkg/database"
)

type UseCase interface {
	InsertCartItem(beegoCtx *beegoContext.Context, request domain.CreateCartRequest) error
	GetListCartItem(beegoCtx *beegoContext.Context, request domain.GetListCartRequest) (*database.Paginator, error)
	DeleteCartItem(beegoCtx *beegoContext.Context, cartIDReq string, customerIDReq int) error
}
