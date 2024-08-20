package main

import (
	"log"
	"net/http"
	"os"

	"api/config"
	_ "api/docs"
	"api/routes"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/swaggo/http-swagger"
)

// @title Library Management API
// @version 1.0
// @description This is a sample server for a library management system.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host https://api-0bw2.onrender.com/
// @BasePath /
func createRouter() http.Handler {
	// Initialize the router
	router := mux.NewRouter()

	// Register routes and pass the database client
	client := config.ConnectDB()
	routes.RegisterBookRoutes(router, client)

	// Custom route for server status
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`<html><body><h1>Server is running!</h1></body></html>`)); err != nil {
			log.Println("Error writing response: ", err)
		}
	})

	// Configure CORS
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		log.Println("Warning: FRONTEND_URL is not set in environment variables")
		frontendURL = "*"
	}
	log.Printf("Frontend URL allowed for CORS: %s\n", frontendURL)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{frontendURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            true, // Enable debugging
	})

	// Add Swagger endpoint
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Custom 404 handler
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Custom 404 - Page Not Found", http.StatusNotFound)
	})

	// Handle preflight requests for CORS
	router.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", frontendURL)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Max-Age", "86400") // Cache preflight results for 24 hours
			w.WriteHeader(http.StatusOK)
			return
		}
		// Handle actual requests here
	}).Methods("GET", "POST", "PUT", "DELETE", "OPTIONS")

	// Wrap the router with CORS middleware
	handler := corsHandler.Handler(router)

	return handler
}

func main() {
	// Load environment variables
	err := config.LoadEnv()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	// Get the port from the environment variable or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server will run on port %s\n", port)

	// Create the router and wrap it with CORS middleware
	handler := createRouter()

	// Start the server
	log.Printf("Server is running and listening on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
