package main

import (
	"fmt"
	"net/url"
	"sync"
	"sync/atomic"
)

type ServerPool struct {
	Backends []*Backend `json:"backends"`
	Current  uint64     `json:"current"` // Used for Round-Robin
	mux      sync.RWMutex
}

/*	GetNextValidPeer() *Backend
	AddBackend(backend *Backend)
	SetBackendStatus(uri *url.URL, alive bool)*/

func (sp *ServerPool) GetNextValidPeer() *Backend {
	sp.mux.RLock()
	defer sp.mux.RUnlock() // safe read lock for backends

	if len(sp.Backends) == 0 {
		return nil
	}

	rrindex := atomic.AddUint64(&sp.Current, 1) % uint64(len(sp.Backends))

	for i := 0; i < len(sp.Backends); i++ {
		idx := (rrindex + uint64(i)) % uint64(len(sp.Backends))
		backend := sp.Backends[idx]

		backend.mux.RLock()
		alive := backend.Alive
		backend.mux.RUnlock()

		if alive {
			return backend
		}
	}

	fmt.Println("No Server Available")
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
func (sp *ServerPool) RemoveBackend(backendURL *url.URL) {
	sp.mux.Lock()
	defer sp.mux.Unlock()
	sli := sp.Backends
	for i, v := range sli {
		if v.URL.String() == backendURL.String() {
			sp.Backends = append(sli[:i], sli[i+1:]...)
			return
		}
	}
}
