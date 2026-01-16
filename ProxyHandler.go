package main
import (
	"net/http"
	"net/http/httputil"
	"sync/atomic"
)
type ProxyHandler struct {
    lb LoadBalancer
}
func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    backend := p.lb.GetNextValidPeer() //getting the proper backend using RR methof
    if backend == nil {
        http.Error(w, "Service unavailable", http.StatusServiceUnavailable) //error if no backend was available 
        return
    }

    // Increase active connections
    atomic.AddInt64(&backend.CurrentConns, 1)
    defer atomic.AddInt64(&backend.CurrentConns, -1)

    proxy := httputil.NewSingleHostReverseProxy(backend.URL) 

    // Preserve the request context
    originalDirector := proxy.Director
    proxy.Director = func(req *http.Request) {
        originalDirector(req)
        req = req.WithContext(r.Context())
    }

    // Handle backend errors
    proxy.ErrorHandler = func(w http.ResponseWriter, req *http.Request, err error) {
        backend.mux.Lock()
        backend.Alive = false
        backend.mux.Unlock()

        http.Error(w, "Backend unavailable", http.StatusBadGateway)
    }

    proxy.ServeHTTP(w, r)
}
