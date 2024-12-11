package main

import (
	"elimt/pkg/pocketbase"
	"elimt/pkg/routes"
	"elimt/pkg/settings"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

var pbClient *pocketbase.PocketBaseClient

func initPocketbase(pbClient *pocketbase.PocketBaseClient, maxRetries int, retryInterval time.Duration) error {
	for i := 0; i < maxRetries; i++ {
		err := pbClient.CheckConnection()
		if err == nil {
			log.Println("Successfully connected to PocketBase")
			return nil
		}
		log.Printf("Failed to connect to PocketBase, attempt %d/%d. Retrying in %s...", i+1, maxRetries, retryInterval)
		time.Sleep(retryInterval)
	}
	return fmt.Errorf("failed to connect to PocketBase after %d attempts", maxRetries)
}

func setupRouter() *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is healthy"))
		log.Println("Server is healthy")
	})

	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		routes.HandleLogin(w, r, pbClient)
	})

	r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
		routes.HandleRegister(w, r, pbClient)
	})

	r.Get("/users", func(w http.ResponseWriter, r *http.Request) {
		routes.HandleUserList(w, r, pbClient)
	})

	return r
}

func main() {
	settings, err := settings.LoadSettings("settings.yml")
	if err != nil {
		log.Fatal(err)
	}

	pbHost := os.Getenv("POCKETBASE_URL")
	if pbHost == "" {
		pbHost = "http://localhost:8090"
	}

	pbClient = pocketbase.NewPocketBaseClient(pbHost)

	err = initPocketbase(pbClient, 10, 5*time.Second)
	if err != nil {
		log.Fatalf("Failed to initialize PocketBase client: %v", err)
	}

	// Setup and start server
	r := setupRouter()
	port := settings.Server.Port
	log.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
