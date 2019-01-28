// Package gocall gives you the ability to create your own out of the box load balancer and API gateway!
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
package gocall
