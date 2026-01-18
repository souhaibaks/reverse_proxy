package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Define a simple handler for the "/" route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi from server 2 ")
	})

	// Define another route for demonstration
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "pong")
	})

	// Start the server on port 8080
	fmt.Println("Server is running on http://localhost:9002")
	err := http.ListenAndServe(":9002", nil)
	if err != nil {
		panic(err)
	}
}
