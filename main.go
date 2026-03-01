package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

// configFile holds the JSON structure of the config file
type configFile struct {
	Port            int      `json:"port"`
	Strategy        string   `json:"strategy"`
	HealthCheckFreq int64    `json:"health_check_frequency"` // nanoseconds
	BackendURLs     []string `json:"backends"`
}

func parseURL(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		log.Fatal(err)
	}
	return u
}

func main() {
	configPath := flag.String("config", "config.json", "Path to the JSON config file")
	flag.Parse()

	// Load config file
	data, err := os.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("Failed to read config file %q: %v", *configPath, err)
	}

	var cfg configFile
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("Failed to parse config file: %v", err)
	}

	// Apply config to global proxyConfig
	proxyConfig = ProxyConfig{
		Port:            cfg.Port,
		Strategy:        cfg.Strategy,
		HealthCheckFreq: time.Duration(cfg.HealthCheckFreq),
	}

	// Build server pool from config backends
	serverPool := &ServerPool{}
	for _, rawURL := range cfg.BackendURLs {
		serverPool.AddBackend(&Backend{
			URL:   parseURL(rawURL),
			Alive: true,
		})
	}

	// If no backends in config, fall back to hardcoded defaults (for dev/testing)
	if len(serverPool.Backends) == 0 {
		serverPool.AddBackend(&Backend{URL: parseURL("http://localhost:9001"), Alive: true})
		serverPool.AddBackend(&Backend{URL: parseURL("http://localhost:9002"), Alive: true})
		serverPool.AddBackend(&Backend{URL: parseURL("http://localhost:9003"), Alive: true})
	}

	handler := &ProxyHandler{lb: serverPool}

	// Start admin API on :8081
	go StartAdminServer(serverPool)

	// Start periodic health checker
	go func() {
		t := time.NewTicker(proxyConfig.HealthCheckFreq)
		for range t.C {
			log.Println("Starting health check ...")
			serverPool.HealthCheck()
			log.Println("Health check completed")
		}
	}()

	addr := fmt.Sprintf(":%d", proxyConfig.Port)
	fmt.Printf("Reverse proxy listening on %s (strategy: %s)\n", addr, proxyConfig.Strategy)
	log.Fatal(http.ListenAndServe(addr, handler))
}
