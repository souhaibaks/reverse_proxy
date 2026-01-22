package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

func StartAdminServer(sb *ServerPool) {
	//GET
	handleGetBackends := func(w http.ResponseWriter, _ *http.Request) {
		jsonResp, err := json.Marshal(sb.Backends)
		if err != nil {
			log.Fatal("Error marshalling JSON:", err)
		}
		w.Write(jsonResp)
	}
	//POST
	handlePostBackends := func(w http.ResponseWriter, r *http.Request) {
		var req struct { //temporary struct to store the backend url in a string format 
			URL string `json:"url"`
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		if  err != nil {
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
		w.WriteHeader(http.StatusCreated)
	}
	//DELETE
	handleDeleteBackends := func(w http.ResponseWriter, r *http.Request) {
		var req struct { //temporary struct to store the backend url in a string format 
			URL string `json:"url"`
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		if  err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		u, err := url.Parse(req.URL)
		if err != nil {
			http.Error(w, "Invalid URL format", http.StatusBadRequest)
			return
		}

		sb.RemoveBackend(u)
		w.WriteHeader(http.StatusOK)
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
	log.Println("Admin server listening on :8081")
	http.ListenAndServe(":8081", nil)
}
