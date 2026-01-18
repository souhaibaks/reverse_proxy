package main

import (
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Backend struct {
	URL          *url.URL `json:"url"`
	Alive        bool     `json:"alive"`
	CurrentConns int64    `json:"current_connections"`
	mux          sync.RWMutex
}

func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	b.Alive = alive
	b.mux.Unlock()

}
func (b *Backend) GetRealStatus() bool { //creating a simple request to check the real status of our backend
	client := &http.Client{
		Timeout: 2 * time.Second, // wait for 2 seconds
	}

	resp, err := client.Get(b.URL.String())
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return true

}
