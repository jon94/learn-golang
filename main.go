package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	muxtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// CustomError is a custom error type
type CustomError struct {
	Message string
}

// Error implements the error interface for CustomError
func (e *CustomError) Error() string {
	return e.Message
}

// Response is a struct to represent the JSON response
type Response struct {
	Message string `json:"message"`
}

func main() {
	// obtain DD_AGENT_HOST
	// ddagenthost := "DD_AGENT_HOST"
	// ddagenthostvalue := os.Getenv(ddagenthost)
	// os.Setenv("DD_AGENT_HOST", ddagenthostvalue)
	// source code integration
	ddgiturl := "DD_GIT_REPOSITORY_URL"
	ddgiturlvalue := os.Getenv(ddgiturl)
	os.Setenv("DD_GIT_REPOSITORY_URL", ddgiturlvalue)
	ddgitsha := "DD_GIT_COMMIT_SHA"
	ddgitshavalue := os.Getenv(ddgitsha)
	os.Setenv("DD_GIT_COMMIT_SHA", ddgitshavalue)

	// Initialize Datadog tracer
	tracer.Start(tracer.WithRuntimeMetrics())
	defer tracer.Stop()

	// Create a new router using gorilla/mux
	router := muxtrace.NewRouter()

	// Define a handler function for the API endpoint
	apiHandler := func(w http.ResponseWriter, r *http.Request) {
		// Start a trace span for this handler
		ctx := r.Context()
		span, ctx := tracer.StartSpanFromContext(ctx, "apiHandler", tracer.ResourceName("/simplegetapi"))
		defer span.Finish()

		// Set the content type to JSON
		w.Header().Set("Content-Type", "application/json")

		// Create a Response struct
		response := Response{Message: "Hello, this is a simple GET API!"}
		log.Println("Simple get API Success!!")

		// Encode the Response struct to JSON and write it to the response writer
		json.NewEncoder(w).Encode(response)
	}

	// Define a handler function for the /getErrorRequest endpoint
	getErrorRequestHandler := func(w http.ResponseWriter, r *http.Request) {
		// Start a trace span for this handler
		ctx := r.Context()
		span, ctx := tracer.StartSpanFromContext(ctx, "getErrorRequestHandler", tracer.ResourceName("/getErrorRequest"))
		defer span.Finish()

		// Set the content type to JSON
		w.Header().Set("Content-Type", "application/json")

		// Return a custom error
		customError := &CustomError{Message: "Custom Error: Something went wrong"}
		w.WriteHeader(http.StatusInternalServerError)

		// Log the error message
		log.Printf("Error Request API triggered with message: %s\n", customError.Message)
		debug.PrintStack() // Capture and print the stack trace

		// Encode the custom error message to JSON and write it to the response writer
		json.NewEncoder(w).Encode(Response{Message: customError.Message})
	}

	// Register the API handler function for the "/api" route using gorilla/mux
	router.HandleFunc("/api", apiHandler).Methods("GET")

	// Register the handler function for the "/getErrorRequest" route
	router.HandleFunc("/getErrorRequest", getErrorRequestHandler).Methods("GET")

	// Attach the router to the default serve mux
	http.Handle("/", router)

	// Start the server on port 8080
	fmt.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("Error:", err)
	}
}
