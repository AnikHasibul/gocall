package gocall

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

// ReverseProxy proxies the target with given http request
func ReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	// parse the url
	uri, _ := url.Parse(target)
	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(uri)
	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
}
