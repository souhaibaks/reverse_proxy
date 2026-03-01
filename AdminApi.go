package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

func StartAdminServer(sb *ServerPool) {
	// GET /status — returns summary + backend list
	handleGetStatus := func(w http.ResponseWriter, _ *http.Request) {
		sb.mux.RLock()
		backends := sb.Backends
		sb.mux.RUnlock()

		var active int
		for _, b := range backends {
			b.mux.RLock()
			if b.Alive {
				active++
			}
			b.mux.RUnlock()
		}

		resp := map[string]interface{}{
			"total_backends":  len(backends),
			"active_backends": active,
			"backends":        backends,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}

	// GET /backends — returns raw backend list
	handleGetBackends := func(w http.ResponseWriter, _ *http.Request) {
		sb.mux.RLock()
		backends := sb.Backends
		sb.mux.RUnlock()

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(backends); err != nil {
			log.Println("Error encoding backends:", err)
		}
	}
	//POST
	handlePostBackends := func(w http.ResponseWriter, r *http.Request) {
		var req struct { //temporary struct to store the backend url in a string format
			URL string `json:"url"`
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		u, err := url.Parse(req.URL)
		if err != nil {
			http.Error(w, "Invalid URL format", http.StatusBadRequest)
			return
		}

		backend := &Backend{ //creating the real backend struct after parsing the url
			URL:   u,
			Alive: true,
		}
		sb.AddBackend(backend)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "backend added"})
	}
	//DELETE
	handleDeleteBackends := func(w http.ResponseWriter, r *http.Request) {
		var req struct { //temporary struct to store the backend url in a string format
			URL string `json:"url"`
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		u, err := url.Parse(req.URL)
		if err != nil {
			http.Error(w, "Invalid URL format", http.StatusBadRequest)
			return
		}

		sb.RemoveBackend(u)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "backend removed"})
	}

	// handler dispatches based on method
	handler := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetBackends(w, r)
		case http.MethodPost:
			handlePostBackends(w, r)
		case http.MethodDelete:
			handleDeleteBackends(w, r)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	}

	http.HandleFunc("/backends", handler)
	http.HandleFunc("/status", handleGetStatus)
	log.Println("Admin server listening on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
