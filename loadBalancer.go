package gocall

import (
	"errors"
	"sync"
	"time"

	"github.com/anikhasibul/sure"
	"github.com/valyala/fasthttp"
)

// LoadBalancer holds basic load balancing mechanisms
type LoadBalancer struct {
	Fallback     func(*fasthttp.RequestCtx, error)
	HealthURL    string
	healthyHosts *sync.Map
	hosts        []string
	mu           sync.Mutex
}

// NewLoadBalancer returns a load balancer
/*
var lb = gocall.NewLoadBalancer([]string{
	"127.0.0.1:1234",
	"127.0.0.1:1235",
	"127.0.0.1:1236",
	}, "/health", 10*time.Second)
func main() {
		fasthttp.ListenAndServe(":8081", proxify)
}
func proxify(ctx *fasthttp.RequestCtx) {
	// Check auth here
	// ....
	// Now pass the request to the target server
	lb.ProxyTheHealthiest(ctx)
	}
*/
func NewLoadBalancer(hosts []string, healthCheckURL string, healthCheckDelay time.Duration) *LoadBalancer {
	// setting up the servers
	var servers = new(sync.Map)
	for _, v := range hosts {
		servers.Store(v, 0)
	}
	// preparing the LoadBalancer
	b := &LoadBalancer{
		HealthURL:    healthCheckURL,
		healthyHosts: servers,
		hosts:        hosts,
	}
	b.Fallback = DefaultFallback
	// health Checker
	go func() {
		for {
			b.healthChecker()
			time.Sleep(healthCheckDelay)
		}
	}()
	// finding up the healthiest server
	go func() {
	}()
	return b
}

// healthChecker checks the health of a server
func (b *LoadBalancer) healthChecker() {
	var wg sync.WaitGroup
	// check available servers
	for _, v := range b.hosts {
		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			// cut out unavailable servers
			if !b.healthCheck(host) {
				b.healthyHosts.Delete(host)
				return
			}
			// add up avalilable servers
			if _, ok := b.healthyHosts.Load(host); !ok {
				b.healthyHosts.Store(host, 0)
			}
		}(v)
	}
	wg.Wait()
}

// FindTheHealthiest returns the server that has the minimum amount of request at calling time
func (b *LoadBalancer) FindTheHealthiest() string {
	// seems inefficeient
	// but it's not that much bad ;)
	var lite int
	var lightest string
	// just for setting up the initial value
	b.healthyHosts.Range(func(k, v interface{}) bool {
		lightest = sure.String(k)
		lite = sure.Int(v)
		return false
	})
	// iteration...
	b.healthyHosts.Range(func(k, v interface{}) bool {
		val := sure.Int(v)
		key := sure.String(k)
		if val < lite {
			lite = val
			lightest = key
		}
		return true
	})
	return lightest
}

// ProxyTheHealthiest sends a reverse proxy request to the heathiest server returned by FindTheHealthiest()
func (b *LoadBalancer) ProxyTheHealthiest(ctx *fasthttp.RequestCtx) {
	// find the ligtest server
	host := b.FindTheHealthiest()
	if host == "" {
		b.Fallback(ctx, errors.New("gocall: no server found for handling this request"))
		return
	}
	b.mu.Lock()
	c, _ := b.healthyHosts.Load(host)
	b.healthyHosts.Store(host, sure.Int(c)+1)
	b.mu.Unlock()
	defer func() {
		b.mu.Lock()
		c, _ := b.healthyHosts.Load(host)
		b.healthyHosts.Store(host, sure.Int(c)+-1)
		b.mu.Unlock()
	}()
	ctx.Response.Header.Set("X-Gate", "BOOM!")
	ctx.URI().SetHost(host)
	ctx.URI().SetScheme("https")
	client := &fasthttp.Client{}
	err := client.Do(&ctx.Request, &ctx.Response)
	if err != nil {
		b.Fallback(ctx, err)
		return
	}
}

func (b *LoadBalancer) healthCheck(uri string) bool {
	// we don't need slow servers
	// so timeout is just 1 second
	// check health
	req := fasthttp.AcquireRequest()
	req.SetRequestURI("http://" + uri + b.HealthURL)
	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{}
	err := client.Do(req, resp)
	if err != nil {
		return false
	}
	//	bodyBytes := resp.Body()
	return true
}

// DefaultFallback function to response if any error occurs on reverse proxy
func DefaultFallback(ctx *fasthttp.RequestCtx, err error) {
	ctx.SetStatusCode(fasthttp.StatusServiceUnavailable)
	ctx.SetBody([]byte(err.Error()))
	//panic(err)
}
