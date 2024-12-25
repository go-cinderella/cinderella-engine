package utils

import (
	"github.com/go-cinderella/cinderella-engine/engine/constant"
	"github.com/go-resty/resty/v2"
	"github.com/klauspost/compress/gzhttp"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"net"
	"net/http"
	"runtime"
	"time"
)

// NewDefaultHttpClient creates new resty Client instance
func NewDefaultHttpClient() *resty.Client {
	client := resty.New()
	client.SetTimeout(1 * time.Minute)
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}
	client.SetTransport(gzhttp.Transport(&nethttp.Transport{
		RoundTripper: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           dialer.DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
			MaxConnsPerHost:       10000,
		},
	}))
	client.SetRetryCount(constant.HttpRetryCount)
	client.SetDebug(true)
	return client
}
