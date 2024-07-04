package http

import (
	"context"
	"errors"
	"github.com/online-store/internal/domain"
	"github.com/online-store/internal/product"
	"github.com/online-store/pkg"
	paging "github.com/online-store/pkg/paging"
	"github.com/online-store/pkg/response"
	"github.com/online-store/pkg/validator"
	"net/http"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/i18n"
)

type ProductHandler struct {
	beego.Controller
	product.UseCase
	i18n.Locale
	response.APIResponseInterface
	time.Duration
}

func NewProductHandler(useCase product.UseCase, executionTimeout time.Duration, apiResponse response.APIResponseInterface) {
	handler := &ProductHandler{
		UseCase:              useCase,
		APIResponseInterface: apiResponse,
		Duration:             executionTimeout,
	}

	beego.Router("/api/v1/products", handler, "get:GetListProduct")
}

func (h *ProductHandler) Prepare() {
	// check user access when needed
	h.Lang = pkg.GetLangVersion(h.Ctx)
	requestTime := time.Now().UnixNano() / int64(time.Millisecond)
	h.Ctx.Input.SetData("request_time", requestTime)
}

func (h *ProductHandler) GetListProduct() {
	ctx, cancel := context.WithTimeout(h.Ctx.Request.Context(), h.Duration)
	defer cancel()

	h.Ctx.Request = h.Ctx.Request.WithContext(ctx)

	var request domain.GetProductListRequest
	request.ProductCategory = h.Ctx.Input.Query("product_category")
	request.Search = h.Ctx.Input.Query("search")

	//validate request
	if err := validator.Validate.ValidateStruct(&request); err != nil {
		h.ResponseError(h.Ctx, http.StatusBadRequest, domain.ApiValidationErrorCode, domain.ErrorCodeText(domain.ApiValidationErrorCode, h.Locale.Lang), err)
		return
	}

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

	//call use case
	res, err := h.UseCase.GetListProduct(h.Ctx, request)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			h.ResponseError(h.Ctx, http.StatusRequestTimeout, domain.RequestTimeoutErrorCode, domain.ErrorCodeText(domain.RequestTimeoutErrorCode, h.Locale.Lang), err)
			return
		}

		h.ResponseError(h.Ctx, http.StatusInternalServerError, domain.ServerErrorCode, domain.ErrorCodeText(domain.ServerErrorCode, h.Locale.Lang), err)
		return
	}

	h.Ok(h.Ctx, h.Tr("message.success"), res, nil)
}
