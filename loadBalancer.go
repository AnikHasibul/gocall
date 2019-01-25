package gocall

import (
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Hosts
type Hosts map[string]int

// LoadBalancer interface implements a load balancer
type LoadBalancer interface {
	TestAll()
	FindTheLightest() *url.URL
}

type balancer struct {
	healthyHosts Hosts
	hosts        []string
	mu           sync.RWMutex
}

func NewLoadBalancer(hosts []string) *balancer {
	points := make(Hosts)
	for _, v := range hosts {
		points[v] = 0
	}
	return &balancer{
		healthyHosts: points,
		hosts:        hosts,
	}
}

func (b *balancer) TestAll() {
	for _, v := range b.hosts {
		go func(host string) {
			if !healthCheck(host) {
				delete(b.healthyHosts, host)
				return
			}
			if _, ok := b.read(host); !ok {
				b.write(host, 0)
			}
		}(v)
	}
}

func (b *balancer) FindTheLightest() string {
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

func healthCheck(url string) bool {
	c := http.Client{
		Timeout: 1 * time.Second,
	}
	resp, err := c.Get(url + "/___/health")
	if err != nil {
		return false
	}
	if resp.StatusCode != 200 {
		return false
	}
	return true
}

func (b *balancer) read(key string) (int, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	c, ok := b.healthyHosts[key]
	return c, ok
}

func (b *balancer) readO(key string) int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.healthyHosts[key]
}
func (b *balancer) write(key string, value int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.healthyHosts[key] = value
}
