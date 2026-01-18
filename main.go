package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
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
	backend3 := &Backend{URL: parseURL("http://localhost:9003"), Alive: true}
	serverPool := &ServerPool{
		Backends: []*Backend{backend1, backend2, backend3},
		Current:  0,
	}

	handler := &ProxyHandler{
		lb: serverPool,
	}
	config:=&ProxyConfig{
		HealthCheckFreq: 30*time.Second,
	}

	fmt.Println("Running on :9000")
	go func() {
		t := time.NewTicker(config.HealthCheckFreq) //interval
		for range t.C { //using a channel to perform health checks
			log.Println("Starting health check ...")
			serverPool.HealthCheck()
			log.Println("Health Check completed")
		}
	}()
	log.Fatal(http.ListenAndServe(":9000", handler))

}
