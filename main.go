package main

import (
	"context"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	beegoContext "github.com/beego/beego/v2/server/web/context"
	"github.com/beego/i18n"
	"github.com/online-store/internal"
	"github.com/online-store/internal/domain"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/online-store/pkg"
	"github.com/online-store/pkg/cache/redis"
	"github.com/online-store/pkg/database"
	"github.com/online-store/pkg/httpclient"
	"github.com/online-store/pkg/response"
	"github.com/online-store/pkg/zaplogger"
	"net"
	"net/http"
	"runtime"
	"time"

	redisConfig "github.com/online-store/pkg/cache"

	productHandler "github.com/online-store/internal/product/delivery/http"
	productRepository "github.com/online-store/internal/product/repository"
	productUC "github.com/online-store/internal/product/usecase"

	customerHandler "github.com/online-store/internal/customer/delivery/http"
	customerRepository "github.com/online-store/internal/customer/repository"
	customerUseCase "github.com/online-store/internal/customer/usecase"

	cartHandler "github.com/online-store/internal/cart/delivery/http"
	cartRepository "github.com/online-store/internal/cart/repository"
	cartUseCase "github.com/online-store/internal/cart/usecase"

	orderHandler "github.com/online-store/internal/order/delivery/http"
	orderRepository "github.com/online-store/internal/order/repository"
	orderUseCase "github.com/online-store/internal/order/usecase"
)

func main() {
	var dbSectionConfig map[string]string

	// zap logger
	zapLog := zaplogger.NewZapLogger(beego.AppConfig.DefaultString("logPath", "./logs/api.log"), "")

	// language
	languages := strings.Split(beego.AppConfig.DefaultString("lang", "en|id"), "|")
	for i := range languages {
		if err := i18n.SetMessage(languages[i], "conf/"+languages[i]+".ini"); err != nil {
			panic("Failed to set message file for l10n")
		}
	}

	//database
	if config, err := beego.AppConfig.GetSection("database"); err != nil {
		panic(err)
	} else {
		dbSectionConfig = config
	}

	gormDb, err := database.New(database.ConfigFromEnvironment(dbSectionConfig))
	if err != nil {
		zapLog.Fatal(err)
	}

	// beego config
	beego.BConfig.Log.AccessLogs = false
	beego.BConfig.Log.EnableStaticLogs = false
	beego.BConfig.Listen.ServerTimeOut = beego.AppConfig.DefaultInt64("serverTimeout", 60)

	if beego.BConfig.RunMode != "prod" {
		// static files swagger
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	beego.BConfig.RecoverFunc = func(context *beegoContext.Context, config *beego.Config) {
		if err := recover(); err != nil {
			fmt.Println("masuk selalu", err)
			var stack string
			hasIndent := beego.BConfig.RunMode != beego.PROD
			out := response.APIResponse{
				Code:      domain.ServerErrorCode,
				Message:   domain.ErrorCodeText(domain.ServerErrorCode, pkg.GetLangVersion(context)),
				Data:      nil,
				Errors:    nil,
				RequestId: context.ResponseWriter.ResponseWriter.Header().Get("X-REQUEST-ID"),
				TimeStamp: time.Now().Format("2006-01-02 15:04:05"),
			}
			for i := 1; ; i++ {
				_, file, line, ok := runtime.Caller(i)
				if !ok {
					break
				}
				stackInfo := fmt.Sprintf("%s:%d", file, line)
				stack = stack + fmt.Sprintln(stackInfo)
			}
			context.Input.SetData("stackTrace", &zaplogger.ListErrors{
				Error: stack,
				Extra: err,
			})
			if context.Output.Status != 0 {
				context.ResponseWriter.WriteHeader(context.Output.Status)
			} else {
				context.ResponseWriter.WriteHeader(500)
			}
			context.Output.JSON(out, hasIndent, false)
		}
	}

	//redis client
	redisConConfig := beego.AppConfig.DefaultString("redisConConfig", `{"conn":"127.0.0.1:6379"}`)
	redisConf := redisConfig.SetConf(redisConConfig)
	redisClient := redisConfig.NewRedisClient(redisConf)

	// resty http client
	restyClient := httpclient.NewRestyHttpClient(
		func(config *httpclient.RestyHttpClient) {
			if beego.BConfig.RunMode != "prod" {
				// active if local & dev environment
				config.Client().SetDebug(true)
			}
		},
		httpclient.ConfigTimeout(time.Duration(beego.AppConfig.DefaultInt64("restyTimeout", 120))*time.Second),
		httpclient.ConfigRetryCount(beego.AppConfig.DefaultInt("restyRetryCount", 3)),
		httpclient.ConfigRetryWaitTime(time.Duration(beego.AppConfig.DefaultInt64("restyRetryWaitTime", 100))*time.Millisecond),
		httpclient.ConfigRetryMaxWaitTime(time.Duration(beego.AppConfig.DefaultInt64("restyRetryMaxWaitTime", 2))*time.Second),
		httpclient.ConfigLogger(zapLog),
		httpclient.ConfigTransport(&http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   time.Duration(beego.AppConfig.DefaultInt64("httpclientconfig::dialcontexttimeout", 30)) * time.Second,
				KeepAlive: time.Duration(beego.AppConfig.DefaultInt64("httpclientconfig::dialcontextkeepalive", 30)) * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     beego.AppConfig.DefaultBool("httpclientconfig::forceattempthttp2", true),
			MaxIdleConns:          beego.AppConfig.DefaultInt("httpclientconfig::maxidleconns", 100),
			MaxIdleConnsPerHost:   beego.AppConfig.DefaultInt("httpclientconfig::maxidleconnsperhost", 100),
			IdleConnTimeout:       time.Duration(beego.AppConfig.DefaultInt64("httpclientconfig::idleconntimeout", 90)) * time.Second,
			TLSHandshakeTimeout:   time.Duration(beego.AppConfig.DefaultInt64("httpclientconfig::tlshandshaketimeout", 10)) * time.Second,
			ExpectContinueTimeout: time.Duration(beego.AppConfig.DefaultInt64("httpclientconfig::expectcontinuetimeout", 1)) * time.Second,
		}),
	)

	// init redis
	redisRepository := redis.NewRedisRepository(redisClient)

	apiResponseInterface := response.NewAPIResponse()

	// init routers filters
	internal.InitRouterFilters(restyClient, zapLog, apiResponseInterface)

	//init repository
	productRepository := productRepository.NewProductRepository(gormDb.Conn())
	customerRepo := customerRepository.NewCustomerRepository(gormDb.Conn())
	cartRepo := cartRepository.NewCartRepository(gormDb.Conn())
	orderRepo := orderRepository.NewOrderRepository(gormDb.Conn())

	//init use case
	productUseCase := productUC.NewProductUseCase(productRepository, redisRepository, zapLog)
	customerUC := customerUseCase.NewCustomerUseCase(customerRepo, zapLog)
	cartUC := cartUseCase.NewCustomerUseCase(cartRepo, zapLog, redisRepository)
	orderUC := orderUseCase.NewOrderUseCase(orderRepo, zapLog)

	// default error handler
	beego.ErrorController(&internal.BaseController{})

	//init handler
	productHandler.NewProductHandler(productUseCase, time.Duration(beego.AppConfig.DefaultInt("executionTimeout", 5))*time.Second, apiResponseInterface)
	customerHandler.NewCustomerHandler(customerUC, time.Duration(beego.AppConfig.DefaultInt("executionTimeout", 5))*time.Second, apiResponseInterface)
	cartHandler.NewProductHandler(cartUC, time.Duration(beego.AppConfig.DefaultInt("executionTimeout", 5))*time.Second, apiResponseInterface)
	orderHandler.NewOrderHandler(orderUC, time.Duration(beego.AppConfig.DefaultInt("executionTimeout", 5))*time.Second, apiResponseInterface)

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		beego.Run()
	}()

	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit

	pid := syscall.Getpid()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer func() {
		cancel()
	}()

	switch sig {
	case syscall.SIGINT:
		log.Println(pid, "Received SIGINT.")
		log.Println(pid, "Waiting for connections to finish...")
		if err := beego.BeeApp.Server.Shutdown(ctx); err != nil {
			log.Fatal("failed shutdown server:", err)
		}
	case syscall.SIGTERM:
		log.Println(pid, "Received SIGTERM.")
		log.Println(pid, "Waiting for connections to finish...")
		if err := beego.BeeApp.Server.Shutdown(ctx); err != nil {
			log.Fatal("failed shutdown server:", err)
		}
	default:
		log.Printf("Received %v: nothing i care about...\n", sig)
	}

	// catching ctx.Done(). timeout of 5 seconds.
	if <-ctx.Done(); true {
		log.Println("shutdown server success.")
	}

	log.Println("server exiting")
}
