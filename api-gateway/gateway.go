package api_gateway

import (
	"net/http/httputil"
	"net/url"
)

func NewServiceProxy(target string) *httputil.ReverseProxy {
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)
	return proxy
}
