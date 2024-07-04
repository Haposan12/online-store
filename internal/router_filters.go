package internal

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/online-store/pkg/httpclient"
	"github.com/online-store/pkg/middleware"
	"github.com/online-store/pkg/response"
	"github.com/online-store/pkg/zaplogger"
)

func InitRouterFilters(restyHttpClient *httpclient.RestyHttpClient, log zaplogger.Logger, apiResponse response.APIResponseInterface) {
	beego.InsertFilterChain("/customer/*", middleware.NewAuthMiddleware(apiResponse).ValidateAuth())
}
