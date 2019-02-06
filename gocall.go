// Package gocall gives you the ability to create your own out of the box load balancer and API gateway!
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
package gocall
