package main

import (
	"sync"
	"sync/atomic"
	"time"
	"net/url"
	"net/http/httputil"
)

type Backend struct {
	URL          *url.URL `json:"url"`
	Alive        bool     `json:"alive"`
	CurrentConns int64    `json:"current_connections"`
	mux          sync.RWMutex
}
type ServerPool struct {
	Backends []*Backend `json:"backends"`
	Current  uint64     `json:"current"` // Used for Round-Robin
	mux          sync.RWMutex
}
type ProxyConfig struct {
	Port            int           `json:"port"`
	Strategy        string        `json:"strategy"` // e.g., "round-robin" or "least-conn"
	HealthCheckFreq time.Duration `json:"health_check_frequency"`
}
type LoadBalancer interface {
	GetNextValidPeer() *Backend
	AddBackend(backend *Backend)
	SetBackendStatus(uri *url.URL, alive bool)
}

func (sp *ServerPool) GetNextValidPeer() *Backend {
	sp.mux.Lock()
	if len(sp.Backends)==0{
		return nil
	}
	rrindex:=atomic.AddUint64(&sp.Current,1)%uint64(len(sp.Backends))

	for i:=rrindex;i<uint64(len(sp.Backends));i++{
			backend:=sp.Backends[i]
			backend.mux.Lock()
			if backend.Alive{
				return backend
			}
			backend.mux.Unlock()

	}

	return nil
}
func (sp *ServerPool) AddBackend(backend *Backend) {
	sp.mux.Lock()
	sp.Backends = append(sp.Backends, backend)
	sp.mux.Unlock()

}
func (sp *ServerPool) SetBackendStatus(uri *url.URL, alive bool) {
	sp.mux.Lock()
	for _, backend := range sp.Backends {
		if backend.URL == uri {
			backend.mux.Lock()
			backend.Alive = alive
			backend.mux.Unlock()
			return
		}

	}
	sp.mux.Unlock()

}
func main() {

}
