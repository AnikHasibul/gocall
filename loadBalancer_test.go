package gocall

import (
	"testing"
	"time"
)

func BenchmarkFindTheLightest(b *testing.B) {
	var lb = NewLoadBalancer([]string{
		"127.0.0.1:1234",
		"127.0.0.1:1235",
		"127.0.0.1:1236",
	}, "/health", 10*time.Second)
	for i := 0; i < b.N; i++ {
		lb.FindTheHealthiest()
	}
}
func BenchmarkFindTheLightest2(b *testing.B) {
	var lb = NewLoadBalancer([]string{
		"127.0.0.1:1234",
		"127.0.0.1:1235",
		"127.0.0.1:1236",
	}, "/health", 10*time.Second)
	for i := 0; i < b.N; i++ {
		lb.FindTheHealthiest2()
	}
}
