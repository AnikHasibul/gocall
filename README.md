# gocall
--
    import "github.com/anikhasibul/gocall"

Package gocall gives you the ability to create your own out of the box load
balancer and API gateway!

```go
var lb = NewLoadBalancer([]string{
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
```

## Usage

#### func  ReverseProxy

```go
func ReverseProxy(target string, res http.ResponseWriter, req *http.Request)
```
ReverseProxy proxies the target with given http request

#### type LoadBalancer

```go
type LoadBalancer struct {
	Fallback  func(http.ResponseWriter, *http.Request, error)
	HealthURL string
}
```

LoadBalancer holds basic load balancing mechanisms

#### func  NewLoadBalancer

```go
func NewLoadBalancer(hosts []string, healthCheckURL string, healthCheckDelay time.Duration) *LoadBalancer
```
NewLoadBalancer returns a load balancer

```go
var lb = NewLoadBalancer([]string{
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
```

#### func (*LoadBalancer) FindTheHealthiest

```go
func (b *LoadBalancer) FindTheHealthiest() string
```
FindTheHealthiest returns the server that has the minimum amount of request at
calling time

#### func (*LoadBalancer) ProxyTheHealthiest

```go
func (b *LoadBalancer) ProxyTheHealthiest(w http.ResponseWriter, r *http.Request)
```
ProxyTheHealthiest sends a reverse proxy request to the heathiest server
returned by FindTheHealthiest()

#### func  DefaultFallback

```go
func DefaultFallback(w http.ResponseWriter, _ *http.Request, err error)
```
DefaultFallback function to response if any error occurs on reverse proxy

