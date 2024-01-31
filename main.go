package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	muxtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Response is a struct to represent the JSON response
type Response struct {
	Message string `json:"message"`
}

func main() {
	// obtain DD_AGENT_HOST
	// ddagenthost := "DD_AGENT_HOST"
	// ddagenthostvalue := os.Getenv(ddagenthost)
	// ddagenthostvalueWithPort := ddagenthostvalue + ":8126"
	// os.Setenv("DD_AGENT_HOST", ddagenthostvalue)
	// log.Println(ddagenthostvalueWithPort)

	// Initialize Datadog tracer
	tracer.Start()
	defer tracer.Stop()

	// Create a new router using gorilla/mux
	router := muxtrace.NewRouter()

	// Define a handler function for the API endpoint
	apiHandler := func(w http.ResponseWriter, r *http.Request) {
		// Start a trace span for this handler
		span := tracer.StartSpan("apiHandler")
		defer span.Finish()

		// Set span tags if needed
		span.SetTag(ext.SpanType, ext.SpanTypeWeb)
		span.SetTag(ext.HTTPMethod, r.Method)
		span.SetTag(ext.HTTPURL, r.URL.Path)

		// Set the content type to JSON
		w.Header().Set("Content-Type", "application/json")

		// Create a Response struct
		response := Response{Message: "Hello, this is a simple GET API!"}
		log.Println("Simple get API Success!!")

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
		log.Println("Error:", err)
	}
}
