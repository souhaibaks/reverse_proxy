package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Define a simple handler for the "/" route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World! This is a simple Go backend.")
	})

	// Define another route for demonstration
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	})

	// Start the server on port 8080
	fmt.Println("Server is running on http://localhost:9001")
	err := http.ListenAndServe(":9001", nil)
	if err != nil {
		panic(err)
	}
}
