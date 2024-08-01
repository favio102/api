// api/index.go
package handler

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/swaggo/http-swagger"
	"library-management/config"
	_ "library-management/docs"
	"library-management/routes"
)

func createRouter() http.Handler {
	router := mux.NewRouter()

	// Initialize the database client
	client := config.ConnectDB()
	routes.RegisterBookRoutes(router, client)

	// Custom route for server status
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`<html><body><h1>Server is running!</h1></body></html>`)); err != nil {
			log.Println("Error writing response:", err)
		}
	})

	// Configure CORS
	frontendURL := os.Getenv("FRONTEND_URL")
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{frontendURL},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	})

	// Wrap the router with CORS middleware
	handler := corsHandler.Handler(router)

	// Add Swagger endpoint
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Custom 404 handler
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Custom 404 - Page Not Found", http.StatusNotFound)
	})

	return handler
}

func Handler(w http.ResponseWriter, r *http.Request) {
	handler := createRouter()
	handler.ServeHTTP(w, r)
}
