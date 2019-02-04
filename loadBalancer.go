package gocall

import (
	"errors"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"github.com/anikhasibul/sure"
)

// LoadBalancer holds basic load balancing mechanisms
type LoadBalancer struct {
	Fallback     func(http.ResponseWriter, *http.Request, error)
	HealthURL    string
	healthyHosts *sync.Map
	hosts        []string
	mu           sync.Mutex
}

// NewLoadBalancer returns a load balancer
/*
	var lb = gocall.NewLoadBalancer([]string{
		"http://127.0.0.1:1234",
		"http://127.0.0.1:1235",
		"http://127.0.0.1:1236",
		}, "/health", 10*time.Second)
	func main() {
        http.HandleFunc("/", proxify)
        http.ListenAndServe(":8080", nil)
	}
	func proxify(w http.ResponseWriter, r *http.Request) {
	// check basic auth here
	// now proxy the request
	lb.ProxyTheHealthiest(w, r)
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
func (b *LoadBalancer) ProxyTheHealthiest(w http.ResponseWriter, r *http.Request) {
	// find the ligtest server
	host := b.FindTheHealthiest()
	if b.Fallback == nil {
		b.Fallback = DefaultFallback
	}
	if host == "" {
		b.Fallback(w, r, errors.New("gocall: no server found for handling this request"))
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
	// parse the url
	uri, _ := url.Parse(host)
	// Update the headers to allow for SSL redirection
	r.URL.Host = uri.Host
	r.URL.Scheme = uri.Scheme
	r.Host = uri.Host
	proxy := httputil.NewSingleHostReverseProxy(uri)
	proxy.ErrorHandler = b.Fallback
	proxy.ServeHTTP(w, r)
}

func (b *LoadBalancer) healthCheck(uri string) bool {
	// we don't need slow servers
	// so timeout is just 1 second
	c := http.Client{
		Timeout: 1 * time.Second,
	}
	// check health
	resp, err := c.Get(uri + b.HealthURL)
	if err != nil {
		return false
	}
	if resp.StatusCode != 200 {
		return false
	}
	return true
}

// DefaultFallback function to response if any error occurs on reverse proxy
func DefaultFallback(w http.ResponseWriter, _ *http.Request, err error) {
	log.Println(err)
	http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
}
