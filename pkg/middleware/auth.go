package middleware

import (
	"errors"
	beego "github.com/beego/beego/v2/server/web"
	beegoContext "github.com/beego/beego/v2/server/web/context"
	"github.com/online-store/internal/domain"
	"github.com/online-store/pkg"
	"github.com/online-store/pkg/response"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	response.APIResponseInterface
}

func NewAuthMiddleware(apiResponse response.APIResponseInterface) *AuthMiddleware {
	return &AuthMiddleware{APIResponseInterface: apiResponse}
}

func (m *AuthMiddleware) ValidateAuth() beego.FilterChain {
	return func(next beego.FilterFunc) beego.FilterFunc {
		return func(ctx *beegoContext.Context) {
			req := ctx.Request

			authHeader := req.Header.Get("Authorization")
			if authHeader == "" {
				m.APIResponseInterface.ResponseError(
					ctx,
					http.StatusUnauthorized,
					domain.MissingTokenErrorCode,
					domain.ErrorCodeText(domain.MissingTokenErrorCode, pkg.GetLangVersion(ctx)),
					errors.New("token is missing"))
				return
			}

			tokenString := strings.Split(authHeader, "Bearer ")[1]
			userID, err := pkg.ValidateJWT(tokenString)
			if err != nil {
				m.APIResponseInterface.ResponseError(
					ctx,
					http.StatusUnauthorized,
					domain.InvalidTokenErrorCode,
					domain.ErrorCodeText(domain.InvalidTokenErrorCode, pkg.GetLangVersion(ctx)),
					err)
				return
			}

			ctx.Input.SetData("userID", userID)
			next(ctx)
		}
	}
}
