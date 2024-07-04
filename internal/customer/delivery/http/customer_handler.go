package http

import (
	"context"
	"errors"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
	"github.com/online-store/internal/customer"
	"github.com/online-store/internal/domain"
	"github.com/online-store/pkg"
	"github.com/online-store/pkg/response"
	"github.com/online-store/pkg/validator"
	"net/http"
	"time"
)

type CustomerHandler struct {
	beego.Controller
	customer.UseCase
	i18n.Locale
	response.APIResponseInterface
	time.Duration
}

func NewCustomerHandler(useCase customer.UseCase, executionTimeout time.Duration, apiResponse response.APIResponseInterface) {
	handler := &CustomerHandler{
		UseCase:              useCase,
		APIResponseInterface: apiResponse,
		Duration:             executionTimeout,
	}

	beego.Router("/auth/v1/customer/login", handler, "post:LoginCustomer")
	beego.Router("/auth/v1/customer/register", handler, "post:CreateCustomer")
}

func (h *CustomerHandler) Prepare() {
	// check user access when needed
	h.Lang = pkg.GetLangVersion(h.Ctx)
	requestTime := time.Now().UnixNano() / int64(time.Millisecond)
	h.Ctx.Input.SetData("request_time", requestTime)
}

func (h *CustomerHandler) LoginCustomer() {
	ctx, cancel := context.WithTimeout(h.Ctx.Request.Context(), h.Duration)
	defer cancel()

	h.Ctx.Request = h.Ctx.Request.WithContext(ctx)

	var request domain.LoginRequest
	if err := h.BindJSON(&request); err != nil {
		h.ResponseError(h.Ctx, http.StatusBadRequest, domain.ApiValidationErrorCode, domain.ErrorCodeText(domain.ApiValidationErrorCode, h.Locale.Lang), err)
		return
	}

	if err := validator.Validate.ValidateStruct(&request); err != nil {
		h.ResponseError(h.Ctx, http.StatusBadRequest, domain.ApiValidationErrorCode, domain.ErrorCodeText(domain.ApiValidationErrorCode, h.Locale.Lang), err)
		return
	}

	res, err := h.UseCase.LoginCustomer(h.Ctx, request)
	if err != nil {
		h.ResponseError(h.Ctx, http.StatusBadRequest, domain.ApiValidationErrorCode, domain.ErrorCodeText(domain.ApiValidationErrorCode, h.Locale.Lang), err)
		return
	}

	token, err := pkg.GenerateJWT(res.CustomerID)
	if err != nil {
		h.ResponseError(h.Ctx, http.StatusInternalServerError, domain.ServerErrorCode, domain.ErrorCodeText(domain.ServerErrorCode, h.Locale.Lang), err)
		return
	}

	h.Ok(h.Ctx, h.Tr("message.success"), map[string]interface{}{
		"access_token": token,
	}, nil)
}

func (h *CustomerHandler) CreateCustomer() {
	ctx, cancel := context.WithTimeout(h.Ctx.Request.Context(), h.Duration)
	defer cancel()

	h.Ctx.Request = h.Ctx.Request.WithContext(ctx)

	var request domain.InsertCustomerRequest
	if err := h.BindJSON(&request); err != nil {
		h.ResponseError(h.Ctx, http.StatusBadRequest, domain.ApiValidationErrorCode, domain.ErrorCodeText(domain.ApiValidationErrorCode, h.Locale.Lang), err)
		return
	}

	if err := validator.Validate.ValidateStruct(&request); err != nil {
		h.ResponseError(h.Ctx, http.StatusBadRequest, domain.ApiValidationErrorCode, domain.ErrorCodeText(domain.ApiValidationErrorCode, h.Locale.Lang), err)
		return
	}

	if err := h.UseCase.InsertCustomer(h.Ctx, request); err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			h.ResponseError(h.Ctx, http.StatusRequestTimeout, domain.RequestTimeoutErrorCode, domain.ErrorCodeText(domain.RequestTimeoutErrorCode, h.Locale.Lang), err)
			return
		}

		if errors.Is(err, domain.ErrUniqueConstraint) {
			h.ResponseError(h.Ctx, http.StatusBadRequest, domain.DataAlreadyExist, domain.ErrorCodeText(domain.DataAlreadyExist, h.Locale.Lang), nil)
			return
		}

		h.ResponseError(h.Ctx, http.StatusInternalServerError, domain.ServerErrorCode, domain.ErrorCodeText(domain.ServerErrorCode, h.Locale.Lang), err)
		return
	}

	h.Ok(h.Ctx, h.Tr("message.success"), nil, nil)
}
