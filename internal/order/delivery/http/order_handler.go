package http

import (
	"context"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
	"github.com/online-store/internal/domain"
	"github.com/online-store/internal/order"
	"github.com/online-store/pkg"
	"github.com/online-store/pkg/response"
	"github.com/online-store/pkg/validator"
	"net/http"
	"time"
)

type OrderHandler struct {
	beego.Controller
	order.UseCase
	i18n.Locale
	response.APIResponseInterface
	time.Duration
}

func NewOrderHandler(useCase order.UseCase, executionTimeout time.Duration, apiResponse response.APIResponseInterface) {
	handler := &OrderHandler{
		UseCase:              useCase,
		APIResponseInterface: apiResponse,
		Duration:             executionTimeout,
	}

	beego.Router("/customer/v1/order/check-out", handler, "post:OrderCheckout")
	beego.Router("/customer/v1/order/payment/:order_id", handler, "post:OrderPayment")
}

func (h *OrderHandler) Prepare() {
	// check user access when needed
	h.Lang = pkg.GetLangVersion(h.Ctx)
	requestTime := time.Now().UnixNano() / int64(time.Millisecond)
	h.Ctx.Input.SetData("request_time", requestTime)
}

func (h *OrderHandler) OrderCheckout() {
	ctx, cancel := context.WithTimeout(h.Ctx.Request.Context(), h.Duration)
	defer cancel()

	h.Ctx.Request = h.Ctx.Request.WithContext(ctx)

	var request domain.CreateOrderCheckoutRequest
	if err := h.BindJSON(&request); err != nil {
		h.ResponseError(h.Ctx, http.StatusBadRequest, domain.ApiValidationErrorCode, domain.ErrorCodeText(domain.ApiValidationErrorCode, h.Locale.Lang), err)
		return
	}

	if err := validator.Validate.ValidateStruct(&request); err != nil {
		h.ResponseError(h.Ctx, http.StatusBadRequest, domain.ApiValidationErrorCode, domain.ErrorCodeText(domain.ApiValidationErrorCode, h.Locale.Lang), err)
		return
	}

	request.CustomerID = h.Ctx.Input.GetData("userID").(int)

	res, err := h.UseCase.CheckoutOrder(h.Ctx, request)
	if err != nil {
		h.ResponseError(h.Ctx, http.StatusInternalServerError, domain.ServerErrorCode, domain.ErrorCodeText(domain.ServerErrorCode, h.Locale.Lang), nil)
		return
	}

	h.Ok(h.Ctx, h.Tr("message.success"), res, nil)
}

func (h *OrderHandler) OrderPayment() {
	ctx, cancel := context.WithTimeout(h.Ctx.Request.Context(), h.Duration)
	defer cancel()

	h.Ctx.Request = h.Ctx.Request.WithContext(ctx)

	var request domain.PaymentRequest
	if err := h.BindJSON(&request); err != nil {
		h.ResponseError(h.Ctx, http.StatusBadRequest, domain.ApiValidationErrorCode, domain.ErrorCodeText(domain.ApiValidationErrorCode, h.Locale.Lang), err)
		return
	}

	if err := validator.Validate.ValidateStruct(&request); err != nil {
		h.ResponseError(h.Ctx, http.StatusBadRequest, domain.ApiValidationErrorCode, domain.ErrorCodeText(domain.ApiValidationErrorCode, h.Locale.Lang), err)
		return
	}

	request.OrderID = h.Ctx.Input.Param(":order_id")
	res, err := h.UseCase.MakePayment(h.Ctx, request)
	if err != nil {
		h.ResponseError(h.Ctx, http.StatusInternalServerError, domain.ServerErrorCode, domain.ErrorCodeText(domain.ServerErrorCode, h.Locale.Lang), nil)
		return
	}

	h.Ok(h.Ctx, h.Tr("message.success"), res, nil)
}
