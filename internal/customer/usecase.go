package customer

import (
	beegoContext "github.com/beego/beego/v2/server/web/context"
	"github.com/online-store/internal/domain"
)

type UseCase interface {
	InsertCustomer(beegoCtx *beegoContext.Context, req domain.InsertCustomerRequest) error
	LoginCustomer(beegoCtx *beegoContext.Context, req domain.LoginRequest) (*domain.Customer, error)
}
