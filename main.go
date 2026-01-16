package main

import (
	"log"
	"net/http"
	"net/url"
)

func parseURL(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		log.Fatal(err)
	}
	return u
}

func main() {
	//just a simple test
	backend1 := &Backend{URL: parseURL("http://localhost:9001"), Alive: true}
	backend2 := &Backend{URL: parseURL("http://localhost:9002"), Alive: true}
	serverPool := &ServerPool{
		Backends: []*Backend{backend1, backend2},
		Current:  0,
	}

	handler := &ProxyHandler{
		lb: serverPool,
	}

	http.ListenAndServe(":8080", handler)
}
