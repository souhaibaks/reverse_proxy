package main
import (
	"sync"
	"net/url"
)
type Backend struct {
	URL          *url.URL `json:"url"`
	Alive        bool     `json:"alive"`
	CurrentConns int64    `json:"current_connections"`
	mux          sync.RWMutex
}