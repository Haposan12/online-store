package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/online-store/pkg/validator"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/beego/i18n"
	"github.com/iancoleman/strcase"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	validatorGo "github.com/go-playground/validator/v10"
)

type (
	APIResponseInterface interface {
		Ok(ctx *context.Context, message string, data interface{}, pageInfo interface{}) error
		ResponseError(ctx *context.Context, httpStatus int, errorCode string, message string, err error) error
	}
)

func NewAPIResponse() APIResponseInterface {
	return &APIResponse{}
}

type APIResponse struct {
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Page      interface{} `json:"page,omitempty"`
	Errors    []Errors    `json:"errors"`
	RequestId string      `json:"request_id"`
	TimeStamp string      `json:"timestamp"`
}

type Errors struct {
	Field       string `json:"field"`
	Description string `json:"message"`
}

func (r APIResponse) Ok(ctx *context.Context, message string, data interface{}, pageInfo interface{}) error {
	ctx.Output.SetStatus(http.StatusOK)
	result := APIResponse{
		Code:      strconv.Itoa(http.StatusOK),
		RequestId: ctx.ResponseWriter.ResponseWriter.Header().Get("x-request-id"),
		Message:   message,
		Data:      data,
		Page:      pageInfo,
		TimeStamp: time.Now().Format("2006-01-02 15:04:05"),
	}
	return ctx.Output.JSON(result, beego.BConfig.RunMode != "prod", false)
}

func (r APIResponse) ResponseError(ctx *context.Context, httpStatus int, errorCode string, message string, err error) error {
	var apiResponse APIResponse
	var errorValidations []Errors = nil

	ctx.Output.SetStatus(httpStatus)

	if err != nil {
		if ctx.Input.RequestBody != nil {
			validateJsonError := checkJsonRequest(err)
			if len(validateJsonError) > 0 {
				errorValidations = validateJsonError
			} else {
				if fields, ok := err.(validatorGo.ValidationErrors); ok {
					lang := "id"
					acceptLang := ctx.Input.Header("Accept-Language")
					if i18n.IsExist(acceptLang) {
						lang = acceptLang
					}
					if trans, found := validator.Validate.GetTranslator(lang); found {
						for _, v := range fields {

							fieldName := v.Field()

							if v.Tag() == "check_fk" {
								param := strings.Split(v.Param(), `:`)
								paramFieldValue := param[0]

								fieldName = paramFieldValue
							}

							if v.Tag() == "number_format" || v.Tag() == "lob_validate" {
								if strings.Contains(v.Field(), "[") {
									fieldName = strings.Split(v.Field(), "[")[0]
								}
							}

							errorValidations = append(errorValidations, Errors{
								Field:       strcase.ToSnake(fieldName),
								Description: strings.ReplaceAll(v.Translate(trans), fieldName, strcase.ToSnake(fieldName)),
							})
						}
					}
				} else {
					errorValidations = append(errorValidations, Errors{
						Field:       "json",
						Description: err.Error(),
					})
				}
			}
		}
	}

	apiResponse.RequestId = ctx.ResponseWriter.ResponseWriter.Header().Get("x-request-id")
	apiResponse.Code = strconv.Itoa(httpStatus)
	apiResponse.Message = message
	apiResponse.TimeStamp = time.Now().Format("2006-01-02 15:04:05")
	apiResponse.Errors = errorValidations
	return ctx.Output.JSON(apiResponse, beego.BConfig.RunMode != "prod", false)
}

// checkJsonRequest Response API
func checkJsonRequest(err error) (response []Errors) {
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError
	var invalidUnmarshalError *json.InvalidUnmarshalError
	switch {
	case errors.As(err, &syntaxError):
		msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
		errorValidation := Errors{
			Field:       "json",
			Description: msg,
		}
		response = append(response, errorValidation)
		return
	case errors.Is(err, io.ErrUnexpectedEOF):
		msg := "Request body contains badly-formed JSON"
		errorValidation := Errors{
			Field:       "json",
			Description: msg,
		}
		response = append(response, errorValidation)
		return
	case errors.As(err, &unmarshalTypeError):
		if ute, ok := err.(*json.UnmarshalTypeError); ok {
			errorValidation := Errors{
				Field:       ute.Field,
				Description: fmt.Sprintf("Parameter %s is invalid (type: %s)", ute.Field, ute.Type),
			}
			response = append(response, errorValidation)
			return
		}
	case errors.As(err, &invalidUnmarshalError):
		if ute, ok := err.(*json.InvalidUnmarshalError); ok {
			errorValidation := Errors{
				Field:       ute.Type.Name(),
				Description: ute.Error(),
			}
			response = append(response, errorValidation)
			return
		}
	}
	return response
}
