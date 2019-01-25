package gocall

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

// Hosts
type Hosts map[string]int

type LoadBalancer struct {
	healthyHosts Hosts
	hosts        []string
	mu           sync.RWMutex
}

func NewLoadBalancer(hosts []string, healthCheckDelay time.Duration) *LoadBalancer {
	points := make(Hosts)
	for _, v := range hosts {
		points[v] = 0
	}
	b := &LoadBalancer{
		healthyHosts: points,
		hosts:        hosts,
	}
	go func() {
		for {
			b.TestAll()
			time.Sleep(healthCheckDelay)
		}
	}()
	return b
}

func (b *LoadBalancer) TestAll() {
	var wg sync.WaitGroup
	b.mu.RLock()
	for _, v := range b.hosts {
		b.mu.RUnlock()
		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			if !healthCheck(host) {
				delete(b.healthyHosts, host)
				return
			}
			if _, ok := b.read(host); !ok {
				b.write(host, 0)
			}
		}(v)
	}
	wg.Wait()
}

func (b *LoadBalancer) FindTheLightest() string {
	// seems inefficeient
	var lite int
	var lightest string
	// just for setting up the initial value
	for k, v := range b.healthyHosts {
		lightest, lite = k, v
		break
	}
	// iteration...
	for k, v := range b.healthyHosts {
		if v < lite {
			lite = v
			lightest = k
		}
	}
	return lightest
}

func (b *LoadBalancer) ProxyTheLightest(w http.ResponseWriter, r *http.Request) {
	host := b.FindTheLightest()
	if host == "" {
		http.Error(w, "Oops!", 500)
		return
	}
	fmt.Println(host)
	b.add(host, 1)
	defer b.add(host, -1)
	// parse the url
	uri, _ := url.Parse(host)
	// Update the headers to allow for SSL redirection
	r.URL.Host = uri.Host
	r.URL.Scheme = uri.Scheme
	r.Header.Set("X-Called", r.Header.Get("Host"))
	r.Host = uri.Host
	proxy := httputil.NewSingleHostReverseProxy(uri)

	proxy.ErrorHandler = b.fallback
	proxy.ServeHTTP(w, r)
}

func healthCheck(url string) bool {
	c := http.Client{
		Timeout: 1 * time.Second,
	}
	resp, err := c.Get(url + "/")
	if err != nil {
		return false
	}
	if resp.StatusCode != 200 {
		return false
	}
	return true
}

func (b *LoadBalancer) read(key string) (int, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	c, ok := b.healthyHosts[key]
	return c, ok
}

func (b *LoadBalancer) readO(key string) int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.healthyHosts[key]
}

func (b *LoadBalancer) add(key string, n int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.healthyHosts[key] = b.healthyHosts[key] + n
}

func (b *LoadBalancer) write(key string, value int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.healthyHosts[key] = value
}

func (b *LoadBalancer) fallback(w http.ResponseWriter, r *http.Request, err error) {
	if !b.Log(err) {
		b.ProxyTheLightest(w, r)
		return
	}
}

func (b *LoadBalancer) Log(err error) bool {
	log.Println(err)
	return true
}
