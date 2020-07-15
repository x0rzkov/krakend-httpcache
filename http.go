// Package httpcache introduces an in-memory-cached http client into the KrakenD stack
package httpcache

import (
	"context"
	"net/http"

	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/proxy"
	"github.com/devopsfaith/krakend/transport/http/client"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
)

// Namespace is the key to use to store and access the custom config data
const Namespace = "github.com/x0rzkov/krakend-httpcache"

var (
	cacheTransport httpcache.Cache
	cacheClient    = http.Client{}
)

// NewHTTPClient creates a HTTPClientFactory using an in-memory-cached http client
func NewHTTPClient(cfg *config.Backend, store string) client.HTTPClientFactory {
	_, ok := cfg.ExtraConfig[Namespace]
	if !ok {
		return client.NewHTTPClient
	}

	switch store {
	case "disk":
		diskCachePath := "./shared/data/cache/krakend/httpcache"
		err = os.MkdirAll(diskCachePath, os.ModePerm)
		if err != nil {
			log.Fatal("ERROR:", err.Error())
		}
		backend := diskcache.New(diskCachePath)
		cacheTransport = httpcache.NewTransport(backend)
  	case "memory":
	  fallthrough
        default:
		cacheTransport = httpcache.NewMemoryCacheTransport()
	}

	cacheClient.Transport = cacheTransport
	return func(_ context.Context) *http.Client {
		return &cacheClient
	}
}

// BackendFactory returns a proxy.BackendFactory that creates backend proxies using
// an in-memory-cached http client
func BackendFactory(cfg *config.Backend) proxy.BackendFactory {
	return proxy.CustomHTTPProxyFactory(NewHTTPClient(cfg))
}
