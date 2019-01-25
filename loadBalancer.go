package gocall

import (
	"net/url"
)

// Hosts
type Hosts map[string]int

// LoadBalancer interface implements a load balancer
type LoadBalancer interface {
	TestAll() (Hosts, error)
	FindTheLightest() *url.URL
}

type balancer struct {
	hosts Hosts
}

func NewLoadBalancer(hosts []string) *balancer {
	points := make(Hosts)
	for _, v := range hosts {
		points[v] = 0
	}
	return &balancer{
		hosts: points,
	}
}

func (b *balancer) TestAll() (Hosts, error) {
	for _, v := range b.hosts {
		_ = v
	}
	return nil, nil
}

func (b *balancer) FindTheLightest() string {
	// seems inefficeient
	var lite int
	var lightest string
	// just for setting up the initial value
	for k, v := range b.hosts {
		lightest, lite = k, v
		break
	}
	// iteration...
	for k, v := range b.hosts {
		if v < lite {
			lite = v
			lightest = k
		}
	}
	return lightest
}
