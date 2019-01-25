package gocall

import (
	"testing"
)

func TestFindTheLightest(t *testing.T) {
	b := NewLoadBalancer([]string{
		"http://google.com",
		"http://m.google.com",
		"http://www.google.com",
	})
	b.hosts["http://m.google.com"] = 100
	b.hosts["http://www.google.com"] = 150
	b.hosts["http://google.com"] = 3000
	if b.FindTheLightest() != "http://m.google.com" {
		t.Fail()
	}
}