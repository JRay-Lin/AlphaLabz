package main

import (
	"alphalabz/pkg/pocketbase"
	"alphalabz/pkg/routes/login"
	"alphalabz/pkg/routes/user"
	"alphalabz/pkg/settings"
	"encoding/json"
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
		log.Println("Server is healthy")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Corrected JSON encoding
		response := map[string]string{"message": "server is healthy"}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	})

	// Login to system
	r.Route("/login", func(r chi.Router) {
		r.Post("/account", func(w http.ResponseWriter, r *http.Request) {
			login.HandleAccountLogin(w, r, pbClient)
		})

		r.Post("/oauth", func(w http.ResponseWriter, r *http.Request) {
			// routes.HandleOAuth(w, r, pbClient)
		})

		r.Post("/sso", func(w http.ResponseWriter, r *http.Request) {
			// routes.HandleSSO(w, r, pbClient)
		})
	})

	// Users route
	r.Route("/users", func(r chi.Router) {
		r.Get("/list", func(w http.ResponseWriter, r *http.Request) {
			user.HandleUserList(w, r, pbClient)
		})

		r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
			user.HandleRegister(w, r, pbClient)
		})

		r.Post("/invite", func(w http.ResponseWriter, r *http.Request) {
			user.HandleInviteNewUser(w, r)
		})

		r.Delete("/remove", func(w http.ResponseWriter, r *http.Request) {
			// routes.HandleUserRemove(w, r, pbClient)
		})

		r.Post("/update", func(w http.ResponseWriter, r *http.Request) {
			// routes.HandleUserUpdate(w, r, pbClient)
		})
	})

	// Lab_book route
	r.Route("/lab_book", func(r chi.Router) {
		r.Get("/list", func(w http.ResponseWriter, r *http.Request) {
			// routes.HandleLabBookList(w, r, pbClient)
		})

		r.Post("/create", func(w http.ResponseWriter, r *http.Request) {
			// routes.HandleLabBookCreate(w, r, pbClient)
		})

		r.Delete("/remove", func(w http.ResponseWriter, r *http.Request) {
			// routes.HandleLabBookRemove(w, r, pbClient)
		})

		r.Post("/update", func(w http.ResponseWriter, r *http.Request) {
			// routes.HandleLabBookUpdate(w, r, pbClient)
		})
	})

	// Schedule routes
	r.Route("/schedule", func(r chi.Router) {
		r.Get("/list", func(w http.ResponseWriter, r *http.Request) {
			// routes.HandleScheduleList(w, r, pbClient)
		})
		r.Post("/create", func(w http.ResponseWriter, r *http.Request) {
			// routes.HandleScheduleCreate(w, r, pbClient)
		})
		r.Delete("/remove", func(w http.ResponseWriter, r *http.Request) {
			// routes.HandleScheduleRemove(w, r, pbClient)
		})
		r.Post("/update", func(w http.ResponseWriter, r *http.Request) {
			// routes.HandleScheduleUpdate(w, r, pbClient)
		})
	})

	// Link routes
	r.Route("/link", func(r chi.Router) {
		r.Get("/list", func(w http.ResponseWriter, r *http.Request) {
			// routes.HandleScheduleList(w, r, pbClient)
		})
		r.Post("/create", func(w http.ResponseWriter, r *http.Request) {
			// routes.HandleScheduleCreate(w, r, pbClient)
		})
		r.Delete("/remove", func(w http.ResponseWriter, r *http.Request) {
			// routes.HandleScheduleRemove(w, r, pbClient)
		})
		r.Post("/update", func(w http.ResponseWriter, r *http.Request) {
			// routes.HandleScheduleUpdate(w, r, pbClient)
		})
	})

	r.Route("/resources", func(r chi.Router) {
		r.Get("/list", func(w http.ResponseWriter, r *http.Request) {
			// routes.HandleScheduleList(w, r, pbClient)
		})
		r.Post("/create", func(w http.ResponseWriter, r *http.Request) {
			// routes.HandleScheduleCreate(w, r, pbClient)
		})
		r.Delete("/remove", func(w http.ResponseWriter, r *http.Request) {
			// routes.HandleScheduleRemove(w, r, pbClient)
		})
		r.Post("/update", func(w http.ResponseWriter, r *http.Request) {
			// routes.HandleScheduleUpdate(w, r, pbClient)
		})

		r.Route("/tags", func(r chi.Router) {
			r.Get("/list", func(w http.ResponseWriter, r *http.Request) {
				// routes.HandleScheduleList(w, r, pbClient)
			})
			r.Post("/create", func(w http.ResponseWriter, r *http.Request) {
				// routes.HandleScheduleCreate(w, r, pbClient)
			})
			r.Delete("/remove", func(w http.ResponseWriter, r *http.Request) {
				// routes.HandleScheduleRemove(w, r, pbClient)
			})
			r.Post("/update", func(w http.ResponseWriter, r *http.Request) {
				// routes.HandleScheduleUpdate(w, r, pbClient)
			})
		})
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

	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	if adminEmail == "" || adminPassword == "" {
		log.Fatal("Missing required environment variables: ADMIN_EMAIL and ADMIN_PASSWORD")
	}

	pbClient, err = pocketbase.NewPocketBaseClient(pbHost, adminEmail, adminPassword)
	if err != nil {
		log.Fatalf("Failed to initialize PocketBase client: %v", err)
	}

	// Verify PocketBase connection with retries
	err = initPocketbase(pbClient, 10, 5*time.Second)
	if err != nil {
		log.Fatalf("Failed to initialize PocketBase client after retries: %v", err)
	}

	// Setup and start server
	r := setupRouter()
	port := settings.Server.Port
	log.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
