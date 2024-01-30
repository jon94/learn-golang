package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	//"github.com/gorilla/mux"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	muxtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
)

// Response is a struct to represent the JSON response
type Response struct {
	Message string `json:"message"`
}

func main() {

	//init dd tracer
	tracer.Start()
	defer tracer.Stop()

	// Create a new router using gorilla/mux
	router := muxtrace.NewRouter()

	// Define a handler function for the API endpoint
	apiHandler := func(w http.ResponseWriter, r *http.Request) {
		// Set the content type to JSON
		w.Header().Set("Content-Type", "application/json")

		// Create a Response struct
		response := Response{Message: "Hello, this is a simple GET API!"}
		fmt.Println("Simple get API Success!!")

		// Encode the Response struct to JSON and write it to the response writer
		json.NewEncoder(w).Encode(response)
	}

	// Register the API handler function for the "/api" route using gorilla/mux
	router.HandleFunc("/api", apiHandler).Methods("GET")

	// Attach the router to the default serve mux
	http.Handle("/", router)

	// Start the server on port 8080
	fmt.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
