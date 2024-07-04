package order

import (
	beegoContext "github.com/beego/beego/v2/server/web/context"
	"github.com/online-store/internal/domain"
)

type UseCase interface {
	CheckoutOrder(beegoCtx *beegoContext.Context, request domain.CreateOrderCheckoutRequest) (*domain.Order, error)
	MakePayment(beegoCtx *beegoContext.Context, request domain.PaymentRequest) (*domain.Payment, error)
}
