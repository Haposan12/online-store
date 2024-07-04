package httpclient

import (
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

type ConfigOption func(config *RestyHttpClient)

func ConfigDebug(enable bool) ConfigOption {
	return func(cfg *RestyHttpClient) {
		cfg.client.SetDebug(enable)
	}
}

func ConfigLogger(logger resty.Logger) ConfigOption {
	return func(cfg *RestyHttpClient) {
		cfg.client.SetLogger(logger)
	}
}

func ConfigTransport(transport *http.Transport) ConfigOption {
	return func(cfg *RestyHttpClient) {
		cfg.client.SetTransport(transport)
	}
}
func ConfigTimeout(timeout time.Duration) ConfigOption {
	return func(cfg *RestyHttpClient) {
		cfg.client.SetTimeout(timeout)
	}
}

func ConfigRetryCount(count int) ConfigOption {
	return func(cfg *RestyHttpClient) {
		cfg.client.SetRetryCount(count)
	}
}

func ConfigRetryWaitTime(waitTime time.Duration) ConfigOption {
	return func(cfg *RestyHttpClient) {
		cfg.client.SetRetryWaitTime(waitTime)
	}
}

func ConfigRetryMaxWaitTime(maxWaitTime time.Duration) ConfigOption {
	return func(cfg *RestyHttpClient) {
		cfg.client.SetRetryMaxWaitTime(maxWaitTime)
	}
}

type RestyHttpClient struct {
	client *resty.Client
}

func restyHttpClientDefaultConfig() RestyHttpClient {
	return RestyHttpClient{
		client: resty.New(),
	}
}

type RestyHttpClientInterface interface {
	Client() *resty.Client
	SetOnBeforeRequest(middleware resty.RequestMiddleware)
	SetOnAfterResponse(middleware resty.ResponseMiddleware)
}

func NewRestyHttpClient(opts ...ConfigOption) *RestyHttpClient {

	cfg := restyHttpClientDefaultConfig()
	for _, fn := range opts {
		if nil != fn {
			fn(&cfg)
		}
	}

	return &RestyHttpClient{
		client: cfg.client,
	}
}

func (r *RestyHttpClient) Client() *resty.Client {
	return r.client
}

func (r *RestyHttpClient) SetOnBeforeRequest(middleware resty.RequestMiddleware) {
	r.client.OnBeforeRequest(middleware)
}

func (r *RestyHttpClient) SetOnAfterResponse(middleware resty.ResponseMiddleware) {
	r.client.OnAfterResponse(middleware)
}
