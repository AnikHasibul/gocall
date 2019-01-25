package gocall

import (
	"testing"
)

func TestFindTheLightest(t *testing.T) {
	b := new(balancer)
	b.healthyHosts["http://m.google.com"] = 100
	b.healthyHosts["http://www.google.com"] = 150
	b.healthyHosts["http://google.com"] = 3000
	if b.FindTheLightest() != "http://m.google.com" {
		t.Fail()
	}
}
