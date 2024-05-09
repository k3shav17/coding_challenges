package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type BackendServers struct {
	URL string
}

type LoadBalancer struct {
	servers      []BackendServers
	lbLocks      sync.Mutex
	currentIndex int
}

func NewLoadBalancer(servers []string) *LoadBalancer {
	lb := &LoadBalancer{
		servers: make([]BackendServers, len(servers)),
	}
	for i, server := range servers {
		lb.servers[i] = BackendServers{URL: server}
	}
	return lb
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	lb.lbLocks.Lock()
	defer lb.lbLocks.Unlock()

	index := lb.currentIndex
	lb.currentIndex = (index + 1) % len(lb.servers)
	fmt.Print(index)

	proxy := httputil.NewSingleHostReverseProxy(&url.URL{Host: lb.servers[index].URL})
	proxy.ServeHTTP(w, r)
}

func main() {

	servers := []string{"localhost:3001", "localhost:3002"}

	lb := NewLoadBalancer(servers)

	fmt.Printf("load balancer is listening on port 3000")
	http.Handle("/", lb)
	http.ListenAndServe(":3000", nil)

}
