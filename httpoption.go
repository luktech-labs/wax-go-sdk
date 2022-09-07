package wax

import (
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

type HttpOption func(client *http.Client)

func EmptyHttpOption() HttpOption {
	return func(client *http.Client) {}
}

func proxyPool(proxies []string) func(*http.Request) (*url.URL, error) {
	return func(r *http.Request) (*url.URL, error) {
		strProxy := proxies[rand.Intn(len(proxies))]
		return &url.URL{Host: strProxy, Scheme: "https"}, nil
	}
}

func WithProxies(proxies []string) HttpOption {
	return func(c *http.Client) {
		if proxies == nil || len(proxies) == 0 {
			return
		}
		c.Transport = &http.Transport{Proxy: proxyPool(proxies)}
	}
}

// WithGlobalTimeout sets a global timeout for every future http calls
// it also overrides any subsequent client calls with a timeout that is set as part
// of the context.
func WithGlobalTimeout(timeout time.Duration) HttpOption {
	return func(c *http.Client) {
		c.Timeout = timeout
	}
}
