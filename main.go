package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Define a handler function for the API endpoint
	apiHandler := func(w http.ResponseWriter, r *http.Request) {
		// Set the content type to JSON
		w.Header().Set("Content-Type", "application/json")

		// Send a JSON response
		fmt.Fprintf(w, `{"message": "Hello, this is a simple GET API!"}`)
	}

	// Register the API handler function for a specific route
	http.HandleFunc("/api", apiHandler)

	// Start the server on port 8080
	fmt.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
