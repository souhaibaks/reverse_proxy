package main

import (
	"fmt"
	"sync"
	"url"
	"time"
)
type Backend struct{
	URL *url.URL `json:"url"`
	Alive bool `json:"alive"`
	CurrentConns int64 `json:"current_connections"`
	mux sync.RWMutex
}
type ServerPool struct {
	Backends []*Backend `json:"backends"`
	Current uint64 `json:"current"` // Used for Round-Robin
}
type ProxyConfig struct {
	Port int `json:"port"`
	Strategy string `json:"strategy"` // e.g., "round-robin" or "least-conn"
	HealthCheckFreq time.Duration `json:"health_check_frequency"`
}
type LoadBalancer interface {
	GetNextValidPeer() *Backend
	AddBackend(backend *Backend)
	SetBackendStatus(uri *url.URL, alive bool)
}
func main(){

}