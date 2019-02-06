package gocall

import (
	"log"

	"github.com/valyala/fasthttp"
)

// ReverseProxy proxies the target with given http request
func ReverseProxy(host string, ctx *fasthttp.RequestCtx) {
	c := fasthttp.Client{}
	ctx.Request.SetHost(host)
	err := c.Do(&ctx.Request, &ctx.Response)
	if err != nil {
		log.Println(err)
		return
	}
}
