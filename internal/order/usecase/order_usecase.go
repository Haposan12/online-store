package usecase

import (
	beegoContext "github.com/beego/beego/v2/server/web/context"
	"github.com/online-store/internal/domain"
	"github.com/online-store/internal/order"
	"github.com/online-store/pkg/zaplogger"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type OrderUseCase struct {
	orderRepo order.Repository
	zapLogger zaplogger.Logger
}

func NewOrderUseCase(orderRepo order.Repository, zapLogger zaplogger.Logger) order.UseCase {
	return &OrderUseCase{
		orderRepo: orderRepo,
		zapLogger: zapLogger,
	}
}

func (u *OrderUseCase) CheckoutOrder(beegoCtx *beegoContext.Context, request domain.CreateOrderCheckoutRequest) (*domain.Order, error) {
	var (
		orderItem  []domain.OrderItem
		orderReq   domain.Order
		totalPrice float64
		orderData  *domain.Order

		err error
	)
	for _, v := range request.Order {
		totalPrice += v.Price
	}

	orderReq = domain.Order{
		TotalPrice: totalPrice,
		CustomerID: request.CustomerID,
		CreatedAt:  time.Now(),
		CreatedBy:  "System",
	}

	//start transaction
	errs := u.orderRepo.DB().Transaction(func(tx *gorm.DB) error {
		//insert orderReq
		orderData, err = u.orderRepo.InsertOrder(beegoCtx.Request.Context(), tx, orderReq)
		if err != nil {
			beegoCtx.Input.SetData("stackTrace", u.zapLogger.SetMessageLog(err))
			return err
		}

		for _, v := range request.Order {
			//get product
			orderItem = append(orderItem, domain.OrderItem{
				ProductID: v.ProductID,
				Price:     v.Price,
				Quantity:  v.Quantity,
				OrderID:   orderData.ID,
				CreatedAt: time.Now(),
				CreatedBy: "System",
			})
		}

		//insert orderReq item
		err = u.orderRepo.InsertOrderItem(beegoCtx.Request.Context(), tx, orderItem)
		if err != nil {
			beegoCtx.Input.SetData("stackTrace", u.zapLogger.SetMessageLog(err))
			return err
		}

		return nil
	})

	if errs != nil {
		beegoCtx.Input.SetData("stackTrace", u.zapLogger.SetMessageLog(errs))
		return nil, errs
	}

	return orderData, nil
}

func (u *OrderUseCase) MakePayment(beegoCtx *beegoContext.Context, request domain.PaymentRequest) (*domain.Payment, error) {
	var (
		data *domain.Payment
		err  error
	)

	orderID, err := strconv.Atoi(request.OrderID)
	if err != nil {
		beegoCtx.Input.SetData("stackTrace", u.zapLogger.SetMessageLog(err))
		return nil, err
	}
	errs := u.orderRepo.DB().Transaction(func(tx *gorm.DB) error {
		//insert payment
		data, err = u.orderRepo.InsertPayment(beegoCtx.Request.Context(), tx, domain.Payment{
			Method:    request.Method,
			Amount:    request.Amount,
			Status:    "Success",
			CreatedAt: time.Now(),
			CreatedBy: "System",
		})
		if err != nil {
			beegoCtx.Input.SetData("stackTrace", u.zapLogger.SetMessageLog(err))
			return err
		}

		//update order
		err = u.orderRepo.UpdateOrder(beegoCtx.Request.Context(), tx, data.ID, orderID)
		if err != nil {
			beegoCtx.Input.SetData("stackTrace", u.zapLogger.SetMessageLog(err))
			return err
		}

		return nil
	})

	if errs != nil {
		beegoCtx.Input.SetData("stackTrace", u.zapLogger.SetMessageLog(errs))
		return nil, errs
	}

	return data, nil
}
