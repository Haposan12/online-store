package http

import (
	"context"
	"errors"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
	"github.com/online-store/internal/cart"
	"github.com/online-store/internal/domain"
	"github.com/online-store/pkg"
	paging "github.com/online-store/pkg/paging"
	"github.com/online-store/pkg/response"
	"github.com/online-store/pkg/validator"
	"net/http"
	"time"
)

type CartHandler struct {
	beego.Controller
	cart.UseCase
	i18n.Locale
	response.APIResponseInterface
	time.Duration
}

func NewProductHandler(useCase cart.UseCase, executionTimeout time.Duration, apiResponse response.APIResponseInterface) {
	handler := &CartHandler{
		UseCase:              useCase,
		APIResponseInterface: apiResponse,
		Duration:             executionTimeout,
	}

	beego.Router("/customer/v1/cart", handler, "post:CreateCart")
	beego.Router("/customer/v1/cart", handler, "get:GetListCart")
	beego.Router("/customer/v1/cart/:id", handler, "delete:DeleteCart")
}

func (h *CartHandler) Prepare() {
	// check user access when needed
	h.Lang = pkg.GetLangVersion(h.Ctx)
	requestTime := time.Now().UnixNano() / int64(time.Millisecond)
	h.Ctx.Input.SetData("request_time", requestTime)
}

func (h *CartHandler) CreateCart() {
	ctx, cancel := context.WithTimeout(h.Ctx.Request.Context(), h.Duration)
	defer cancel()

	h.Ctx.Request = h.Ctx.Request.WithContext(ctx)

	var request domain.CreateCartRequest
	if err := h.BindJSON(&request); err != nil {
		h.ResponseError(h.Ctx, http.StatusBadRequest, domain.ApiValidationErrorCode, domain.ErrorCodeText(domain.ApiValidationErrorCode, h.Locale.Lang), err)
		return
	}

	request.CustomerID = h.Ctx.Input.GetData("userID").(int)
	if err := validator.Validate.ValidateStruct(&request); err != nil {
		h.ResponseError(h.Ctx, http.StatusBadRequest, domain.ApiValidationErrorCode, domain.ErrorCodeText(domain.ApiValidationErrorCode, h.Locale.Lang), err)
		return
	}

	if err := h.UseCase.InsertCartItem(h.Ctx, request); err != nil {
		if errors.Is(err, domain.ErrForeignKeyConstraint) {
			customMsg := fmt.Sprintf("Data product")
			h.ResponseError(h.Ctx, http.StatusBadRequest, domain.ForeignKeyConstraintErrorCode, domain.ErrorCodeText(domain.ForeignKeyConstraintErrorCode, h.Locale.Lang, customMsg), nil)
			return
		}
		h.ResponseError(h.Ctx, http.StatusInternalServerError, domain.ServerErrorCode, domain.ErrorCodeText(domain.ServerErrorCode, h.Locale.Lang), nil)
		return
	}

	h.Ok(h.Ctx, h.Tr("message.success"), nil, nil)
}

func (h *CartHandler) GetListCart() {
	ctx, cancel := context.WithTimeout(h.Ctx.Request.Context(), h.Duration)
	defer cancel()

	h.Ctx.Request = h.Ctx.Request.WithContext(ctx)

	var request domain.GetListCartRequest
	limit, page, err := paging.PageAndPageSizeValidation(h.Ctx.Input.Query("limit"), h.Ctx.Input.Query("page"))
	if err != nil {
		h.ResponseError(
			h.Ctx,
			http.StatusBadRequest,
			domain.InvalidUrlQueryParamErrorCode,
			domain.ErrorCodeText(domain.InvalidUrlQueryParamErrorCode, h.Locale.Lang),
			domain.ErrInvalidUrlQueryParam,
		)
		return
	}

	request.Limit = limit
	request.Page = page
	request.CustomerID = h.Ctx.Input.GetData("userID").(int)

	res, err := h.UseCase.GetListCartItem(h.Ctx, request)
	if err != nil {
		h.ResponseError(h.Ctx, http.StatusInternalServerError, domain.ServerErrorCode, domain.ErrorCodeText(domain.ServerErrorCode, h.Locale.Lang), nil)
		return
	}

	h.Ok(h.Ctx, h.Tr("message.success"), res, nil)
}

func (h *CartHandler) DeleteCart() {
	ctx, cancel := context.WithTimeout(h.Ctx.Request.Context(), h.Duration)
	defer cancel()

	h.Ctx.Request = h.Ctx.Request.WithContext(ctx)

	cartID := h.Ctx.Input.Param(":id")
	customerID := h.Ctx.Input.GetData("userID").(int)

	if err := h.UseCase.DeleteCartItem(h.Ctx, cartID, customerID); err != nil {
		h.ResponseError(h.Ctx, http.StatusInternalServerError, domain.ServerErrorCode, domain.ErrorCodeText(domain.ServerErrorCode, h.Locale.Lang), nil)
		return
	}

	h.Ok(h.Ctx, h.Tr("message.success"), nil, nil)
}
