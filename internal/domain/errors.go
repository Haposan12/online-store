package domain

import (
	"errors"
	"fmt"
	"github.com/beego/i18n"
	"strings"
)

const (
	ServerErrorCode = "STR-API-500"

	RequestForbiddenErrorCode     = "STR-API-001"
	ResourceNotFoundErrorCode     = "STR-API-002"
	RequestTimeoutErrorCode       = "STR-API-003"
	ApiValidationErrorCode        = "STR-API-004"
	DataNotFoundErrorCode         = "STR-API-005"
	ServiceCommunicationErrorCode = "STR-API-006"
	InvalidCredentialErrorCode    = "STR-API-007"
	InvalidUrlParamErrorCode      = "STR-API-008"
	InvalidUrlQueryParamErrorCode = "STR-API-009"
	DataAlreadyExist              = "STR-API-010"
	DataAlreadyExistByCondition   = "STR-API-011"
	ForeignKeyConstraintErrorCode = "STR-API-012"

	PgCodeUniqueConstraint     = "23505"
	PgCodeForeignKeyConstraint = "23503"
	PgCodeNoPartitionRelation  = "23514"

	MissingTokenErrorCode = "STR-AUTH-001"
	InvalidTokenErrorCode = "STR-AUTH-002"
)

var (
	ErrInvalidUrlQueryParam = errors.New("query param is invalid")
	ErrForeignKeyConstraint = errors.New("foreign key constraint")
	ErrUniqueConstraint     = errors.New("unique_constraint")
)

func ErrorCodeText(code, locale string, args ...interface{}) string {
	switch code {
	case ApiValidationErrorCode:
		return i18n.Tr(locale, "message.errorValidation", args)
	case ServerErrorCode:
		return i18n.Tr(locale, "message.errorServerError", args)
	case RequestTimeoutErrorCode:
		return i18n.Tr(locale, "message.errorRequestTimeout", args)
	case DataAlreadyExist:
		return i18n.Tr(locale, "message.errorDataAlreadyExist", args)
	case InvalidTokenErrorCode:
		return i18n.Tr(locale, "message.errorInvalidToken", args)
	case MissingTokenErrorCode:
		return i18n.Tr(locale, "message.errorMissingToken", args)
	case ForeignKeyConstraintErrorCode:
		msg := i18n.Tr(locale, "message.errorForeignKeyConstraint", nil)
		if len(args) > 0 {
			msg = strings.ReplaceAll(msg, "data", fmt.Sprintf("%v", args[0]))
		}
		return msg
	default:
		return ""
	}
}
